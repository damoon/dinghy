package dinghy

import (
	"context"
	"fmt"
	"io"
	"mime"
	"net/http"
	"path/filepath"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
)

type Storage struct {
	Client *s3.S3
	Bucket string
}

func (m Storage) exists(ctx context.Context, path string) (bool, string, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "s3: stat object")
	defer span.Finish()

	span.LogFields(log.String("path", path))

	head, err := m.Client.HeadObjectWithContext(ctx, &s3.HeadObjectInput{
		Bucket: aws.String(m.Bucket),
		Key:    aws.String(path),
	})

	if err == nil {
		etag := strings.Trim(*head.ETag, "\"")
		return true, etag, nil
	}

	if strings.Contains(err.Error(), "NotFound: Not Found") {
		return false, "", nil
	}

	span.LogFields(log.Error(err))

	return false, "", fmt.Errorf("stat object %s: %v", path, err)
}

func (m Storage) healthy(ctx context.Context) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "s3: stat bucket")
	defer span.Finish()

	span.LogFields(log.String("bucket", m.Bucket))

	_, err := m.Client.HeadBucketWithContext(ctx, &s3.HeadBucketInput{
		Bucket: aws.String(m.Bucket),
	})
	if err != nil {
		span.LogFields(log.Error(err))
		return fmt.Errorf("check bucket: %v", err)
	}

	return nil
}

type Directory struct {
	Path        string
	Directories []string
	Files       []File
}

type File struct {
	Name        string
	DownloadURL string
	Size        int64
	Icon        string
	Thumbnail   string `json:"Thumbnail,omitempty"`
}

func (m Storage) list(ctx context.Context, prefix string) (Directory, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "s3: list prefix")
	defer span.Finish()

	span.LogFields(log.String("prefix", prefix))

	//	prefix = strings.TrimPrefix(prefix, "/")

	l := Directory{
		Path:        strings.TrimPrefix(prefix, "/"),
		Directories: []string{},
		Files:       []File{},
	}

	ls, err := m.Client.ListObjectsV2WithContext(ctx, &s3.ListObjectsV2Input{
		Bucket:    aws.String(m.Bucket),
		Prefix:    aws.String(strings.TrimPrefix(filesDirectory, "/") + prefix),
		Delimiter: aws.String("/"),
	})
	if err != nil {
		span.LogFields(log.Error(err))
		return Directory{}, fmt.Errorf("list %s: %v", prefix, err)
	}

	for _, object := range ls.CommonPrefixes {
		name := *object.Prefix
		name = strings.TrimPrefix(name, filesDirectory+prefix)
		name = strings.TrimSuffix(name, "/")
		l.Directories = append(l.Directories, name)
	}

	for _, object := range ls.Contents {
		name := strings.TrimPrefix(*object.Key, filesDirectory+prefix)
		url := strings.TrimPrefix(*object.Key+"?redirect", filesDirectory+"/")
		file := File{
			Name:        name,
			Size:        *object.Size,
			DownloadURL: url,
			Icon:        icon(name),
		}

		if thumbnailSupported(name) {
			file.Thumbnail = url + "&thumbnail"
		}

		l.Files = append(l.Files, file)
	}

	return l, nil
}

func (m Storage) presign(ctx context.Context, method, path string) (string, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "s3: presign")
	defer span.Finish()

	span.LogFields(
		log.String("method", method),
		log.String("path", path),
	)

	var req *request.Request
	switch method {
	case http.MethodGet:
		extention := filepath.Ext(path)
		typee := mime.TypeByExtension(extention)
		req, _ = m.Client.GetObjectRequest(&s3.GetObjectInput{
			Bucket: aws.String(m.Bucket),
			Key:    aws.String(path),
			//			ResponseContentDisposition: aws.String("attachment"), // forces browser to download; vs inline to open directly
			ResponseContentType: aws.String(typee),
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

func (m Storage) delete(ctx context.Context, path string) error {
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

func (m Storage) upload(ctx context.Context, path string, file io.ReadSeeker) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "s3: upload")
	defer span.Finish()

	span.LogFields(
		log.String("path", path),
	)

	_, err := m.Client.PutObjectWithContext(ctx, &s3.PutObjectInput{
		Bucket: aws.String(m.Bucket),
		Key:    aws.String(path),
		Body:   file,
	})
	if err != nil {
		span.LogFields(log.Error(err))
		return err
	}

	return nil
}

func (m Storage) download(ctx context.Context, path string, w io.WriterAt) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "s3: download")
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
