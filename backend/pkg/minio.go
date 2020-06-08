package dinghy

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

type Storage struct {
	Client *s3.S3
	Bucket string
}

func (m Storage) exists(ctx context.Context, path string) (bool, error) {
	_, err := m.Client.HeadObjectWithContext(ctx, &s3.HeadObjectInput{
		Bucket: aws.String(m.Bucket),
		Key:    aws.String(path),
	})

	if err == nil {
		return true, nil
	}

	if strings.Contains(err.Error(), "NotFound: Not Found") {
		return false, nil
	}

	return false, fmt.Errorf("stat object %s: %v", path, err)
}

func (m Storage) healthy(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	err := m.Client.WaitUntilBucketExistsWithContext(ctx, &s3.HeadBucketInput{
		Bucket: aws.String(m.Bucket),
	})
	if err != nil {
		return fmt.Errorf("wait for bucket: %v", err)
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
}

func (m Storage) list(ctx context.Context, prefix string) (Directory, error) {
	prefix = strings.TrimPrefix(prefix, "/")

	l := Directory{
		Path:        prefix,
		Directories: []string{},
		Files:       []File{},
	}

	ls, err := m.Client.ListObjectsV2WithContext(ctx, &s3.ListObjectsV2Input{
		Bucket:    aws.String(m.Bucket),
		Prefix:    aws.String(prefix),
		Delimiter: aws.String("/"),
	})
	if err != nil {
		return Directory{}, fmt.Errorf("list %s: %v", prefix, err)
	}

	for _, object := range ls.CommonPrefixes {
		name := strings.TrimPrefix(strings.TrimSuffix(*object.Prefix, "/"), prefix)
		l.Directories = append(l.Directories, name)
	}

	for _, object := range ls.Contents {
		name := strings.TrimPrefix(*object.Key, prefix)
		url := *object.Key + "?redirect"

		l.Files = append(l.Files, File{
			Name:        name,
			Size:        *object.Size,
			DownloadURL: url,
		})
	}

	return l, nil
}

func (m Storage) presign(method, path string) (string, error) {
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
		return "", fmt.Errorf("method %s not supported", method)
	}

	url, err := req.Presign(10 * time.Minute)
	if err != nil {
		return "", fmt.Errorf("presign path %s %s: %v", method, path, err)
	}

	return url, nil
}

func (m Storage) delete(ctx context.Context, path string) error {
	_, err := m.Client.DeleteObjectWithContext(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(m.Bucket),
		Key:    aws.String(path),
	})
	if err != nil {
		return err
	}

	return nil
}

func (m Storage) upload(ctx context.Context, path string, file io.ReadSeeker) error {
	_, err := m.Client.PutObjectWithContext(ctx, &s3.PutObjectInput{
		Bucket: aws.String(m.Bucket),
		Key:    aws.String(path),
		Body:   file,
	})
	if err != nil {
		return err
	}

	return nil
}

func (m Storage) download(ctx context.Context, path string, w io.WriterAt) error {
	downloader := s3manager.NewDownloaderWithClient(m.Client)

	_, err := downloader.Download(w, &s3.GetObjectInput{
		Bucket: aws.String(m.Bucket),
		Key:    aws.String(path),
	})
	if err != nil {
		return fmt.Errorf("download file: %v", err)
	}

	return nil
}
