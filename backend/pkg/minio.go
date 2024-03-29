package dinghy

import (
	"context"
	"fmt"
	"io"
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
)

type MinioAdapter struct {
	Client *s3.S3
	Bucket string
}

func (m MinioAdapter) exists(ctx context.Context, path string) (bool, string, string, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "s3: stat object")
	defer span.Finish()

	span.LogFields(log.String("path", path))

	head, err := m.Client.HeadObjectWithContext(ctx, &s3.HeadObjectInput{
		Bucket: aws.String(m.Bucket),
		Key:    aws.String(path),
	})

	if err == nil {
		contentType := *head.ContentType
		etag := strings.Trim(*head.ETag, "\"")
		return true, etag, contentType, nil
	}

	if strings.Contains(err.Error(), "NotFound: Not Found") {
		return false, "", "", nil
	}

	span.LogFields(log.Error(err))

	return false, "", "", fmt.Errorf("stat object %s: %v", path, err)
}

type Directory struct {
	Path        string
	Directories []string
	Files       []File
}

type File struct {
	Name        string
	Path        string
	DownloadURL string
	Size        int64
	Icon        string
	Thumbnail   string `json:"Thumbnail,omitempty"`
	Archive     bool
}

type byFileName []File

func (s byFileName) Len() int {
	return len(s)
}

func (s byFileName) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s byFileName) Less(i, j int) bool {
	return strings.ToLower(s[i].Name) < strings.ToLower(s[j].Name)
}

type byCaseInsensitiveString []string

func (s byCaseInsensitiveString) Len() int {
	return len(s)
}

func (s byCaseInsensitiveString) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s byCaseInsensitiveString) Less(i, j int) bool {
	return strings.ToLower(s[i]) < strings.ToLower(s[j])
}

func (m MinioAdapter) list(ctx context.Context, prefix string) (Directory, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "s3: list prefix")
	defer span.Finish()

	span.LogFields(log.String("prefix", prefix))

	l := Directory{
		Path:        strings.TrimPrefix(prefix, "/"),
		Directories: []string{},
		Files:       []File{},
	}

	perPage := func(page *s3.ListObjectsV2Output, lastPage bool) bool {
		for _, object := range page.CommonPrefixes {
			name := *object.Prefix
			name = strings.TrimPrefix(name, filesDirectory+prefix)
			name = strings.TrimSuffix(name, "/")
			l.Directories = append(l.Directories, name)
		}

		for _, object := range page.Contents {
			name := strings.TrimPrefix(*object.Key, filesDirectory+prefix)

			redirect := ""
			if shouldUsePresignRedirect(name) {
				redirect = "?redirect"
			}

			url := strings.TrimPrefix(*object.Key+redirect, filesDirectory+"/")

			file := File{
				Name:        name,
				Path:        prefix + name,
				Size:        *object.Size,
				DownloadURL: url,
				Icon:        icon(name),
				Archive:     canBeExtracted(name, l.Directories),
			}

			if thumbnailSupported(name) {
				file.Thumbnail = url + "&thumbnail"
			}

			l.Files = append(l.Files, file)
		}

		return lastPage
	}

	err := m.Client.ListObjectsV2PagesWithContext(ctx, &s3.ListObjectsV2Input{
		Bucket:    aws.String(m.Bucket),
		Prefix:    aws.String(strings.TrimPrefix(filesDirectory, "/") + prefix),
		Delimiter: aws.String("/"),
	}, perPage)
	if err != nil {
		span.LogFields(log.Error(err))
		return Directory{}, fmt.Errorf("list %s: %v", prefix, err)
	}

	sort.Sort(byFileName(l.Files))
	sort.Sort(byCaseInsensitiveString(l.Directories))

	return l, nil
}

func shouldUsePresignRedirect(name string) bool {
	extensions := []string{
		".html",
		".htm",
		".css",
		".js",
	}

	for _, ext := range extensions {
		if strings.HasSuffix(name, ext) {
			return false
		}
	}

	return true
}

func (m MinioAdapter) presign(ctx context.Context, method, path string) (string, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "s3: presign")
	defer span.Finish()

	span.LogFields(
		log.String("method", method),
		log.String("path", path),
	)

	var req *request.Request
	switch method {
	case http.MethodGet:
		req, _ = m.Client.GetObjectRequest(&s3.GetObjectInput{
			Bucket: aws.String(m.Bucket),
			Key:    aws.String(path),
		})
	case http.MethodPut:
		req, _ = m.Client.PutObjectRequest(&s3.PutObjectInput{
			Bucket: aws.String(m.Bucket),
			Key:    aws.String(path),
		})
	case http.MethodDelete:
		req, _ = m.Client.DeleteObjectRequest(&s3.DeleteObjectInput{
			Bucket: aws.String(m.Bucket),
			Key:    aws.String(path),
		})
	default:
		err := fmt.Errorf("method %s not supported", method)
		span.LogFields(log.Error(err))
		return "", err
	}

	url, err := req.Presign(10 * time.Minute)
	if err != nil {
		span.LogFields(log.Error(err))
		return "", fmt.Errorf("presign path %s %s: %v", method, path, err)
	}

	return url, nil
}

func (m MinioAdapter) delete(ctx context.Context, path string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "s3: delete")
	defer span.Finish()

	span.LogFields(
		log.String("path", path),
	)

	_, err := m.Client.DeleteObjectWithContext(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(m.Bucket),
		Key:    aws.String(path),
	})
	if err != nil {
		span.LogFields(log.Error(err))
		return err
	}

	return nil
}

func (m MinioAdapter) deleteRecursive(ctx context.Context, prefix string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "s3: delete recursive")
	defer span.Finish()

	span.LogFields(
		log.String("path", prefix),
	)

	ls, err := m.Client.ListObjectsV2WithContext(ctx, &s3.ListObjectsV2Input{
		Bucket: aws.String(m.Bucket),
		Prefix: aws.String(strings.TrimPrefix(filesDirectory, "/") + prefix),
	})
	if err != nil {
		span.LogFields(log.Error(err))
		return fmt.Errorf("list %s: %v", prefix, err)
	}

	for _, object := range ls.Contents {
		err = m.delete(ctx, *object.Key)
		if err != nil {
			span.LogFields(log.Error(err))
			return fmt.Errorf("delete %s: %v", *object.Key, err)
		}
	}

	return nil
}

func (m MinioAdapter) upload(ctx context.Context, path string, file io.ReadSeeker, contentType string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "s3: upload")
	defer span.Finish()

	span.LogFields(
		log.String("path", path),
	)

	put := &s3.PutObjectInput{
		Bucket: aws.String(m.Bucket),
		Key:    aws.String(path),
		Body:   file,
	}

	if contentType != "" {
		put.ContentType = aws.String(contentType)
	}

	_, err := m.Client.PutObjectWithContext(ctx, put)
	if err != nil {
		span.LogFields(log.Error(err))
		return err
	}

	return nil
}

func (m MinioAdapter) uploadRecursive(ctx context.Context, src, target string) error {
	return filepath.Walk(src,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			if info.IsDir() {
				return nil
			}

			file, err := os.Open(path)
			if err != nil {
				return err
			}
			defer file.Close()

			extention := filepath.Ext(path)
			contentType := mime.TypeByExtension(extention)

			err = m.upload(ctx, target+strings.TrimPrefix(path, src), file, contentType)
			if err != nil {
				return err
			}

			return nil
		})
}

func (m MinioAdapter) download(ctx context.Context, path string, w io.WriterAt) error {
	span, _ := opentracing.StartSpanFromContext(ctx, "s3: download")
	defer span.Finish()

	span.LogFields(
		log.String("path", path),
	)

	downloader := s3manager.NewDownloaderWithClient(m.Client)

	_, err := downloader.Download(w, &s3.GetObjectInput{
		Bucket: aws.String(m.Bucket),
		Key:    aws.String(path),
	})
	if err != nil {
		span.LogFields(log.Error(err))
		return fmt.Errorf("download file: %v", err)
	}

	return nil
}
