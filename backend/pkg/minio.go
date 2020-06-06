package dinghy

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/minio/minio-go"
)

type MinioStorage struct {
	client   *minio.Client
	bucket   string
	location string
}

func NewMinioStorage(c *minio.Client, l, b string) MinioStorage {
	s := MinioStorage{
		client:   c,
		bucket:   b,
		location: l,
	}

	return s
}

func (m MinioStorage) EnsureBucket() {
	for {
		err := m.createBucketIfMissing()
		if err != nil {
			log.Printf("failed to ensure bucket exists: %v\n", err)
			time.Sleep(5 * time.Second)

			continue
		}

		return
	}
}

func (m MinioStorage) createBucketIfMissing() error {
	exists, err := m.client.BucketExists(m.bucket)
	if err != nil {
		return fmt.Errorf("failed to access bucket %s: %s", m.bucket, err)
	}

	if exists {
		return nil
	}

	err = m.client.MakeBucket(m.bucket, m.location)
	if err != nil {
		return fmt.Errorf("failed to create bucket %s: %s", m.bucket, err)
	}

	log.Printf("bucket %s created\n", m.bucket)

	return nil
}

func (m MinioStorage) healthy(ctx context.Context) error {
	found, err := m.bucketExists(ctx)
	if err != nil {
		return fmt.Errorf("failed to check for bucket: %s", err)
	}

	if !found {
		return fmt.Errorf("bucket is missing: %s", err)
	}

	return nil
}

type exists struct {
	found bool
	err   error
}

func (m MinioStorage) bucketExists(ctx context.Context) (bool, error) {
	res := make(chan exists)

	go func() {
		found, err := m.client.BucketExists(m.bucket)
		res <- exists{found: found, err: err}
	}()

	select {
	case res := <-res:
		return res.found, res.err
	case <-ctx.Done():
		return false, ctx.Err()
	}
}

func (m MinioStorage) exists(ctx context.Context, objectName string) (bool, error) {
	if objectName == "" {
		return false, nil
	}

	res := make(chan exists)

	go func() {
		_, err := m.client.StatObject(m.bucket, objectName, minio.StatObjectOptions{})
		if err != nil {
			errResponse := minio.ToErrorResponse(err)
			if errResponse.Code == "NoSuchKey" {
				res <- exists{found: false}
			}
			res <- exists{err: err}
		}
		res <- exists{found: true}
	}()

	select {
	case res := <-res:
		return res.found, res.err
	case <-ctx.Done():
		return false, ctx.Err()
	}
}

func (m MinioStorage) presign(method, objectName string) (*url.URL, error) {
	switch method {
	case http.MethodGet:
		return m.client.PresignedGetObject(m.bucket, objectName, 10*time.Minute, url.Values{})
	case http.MethodHead:
		return m.client.PresignedHeadObject(m.bucket, objectName, 10*time.Minute, url.Values{})
	case http.MethodPut:
		return m.client.PresignedPutObject(m.bucket, objectName, 10*time.Minute)
	}

	return nil, fmt.Errorf("method %v not known", method)
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

func (m MinioStorage) list(prefix string) (Directory, error) {
	prefix = strings.TrimSuffix(prefix, "/")
	prefix = prefix + "/"

	if prefix == "/" {
		prefix = ""
	}

	l := Directory{
		Path:        strings.TrimSuffix(prefix, "/"),
		Directories: []string{},
		Files:       []File{},
	}

	doneCh := make(chan struct{})
	defer close(doneCh)

	isRecursive := false
	objectCh := m.client.ListObjectsV2(m.bucket, prefix, isRecursive, doneCh)

	for object := range objectCh {
		if object.Err != nil {
			return Directory{}, object.Err
		}

		name := strings.TrimPrefix(object.Key, prefix)

		if strings.HasSuffix(name, "/") {
			l.Directories = append(l.Directories, strings.TrimSuffix(name, "/"))
			continue
		}

		//		url := "http://backend:8080/" + object.Key
		url, err := m.client.PresignedGetObject(m.bucket, object.Key, 10*time.Minute, url.Values{})
		if err != nil {
			return Directory{}, fmt.Errorf("presign download: %v", err)
		}

		l.Files = append(l.Files, File{
			Name:        name,
			Size:        object.Size,
			DownloadURL: url.String(),
		})
	}

	return l, nil
}

func (m MinioStorage) delete(ctx context.Context, path string) error {
	err := m.client.RemoveObject(m.bucket, path)
	if err != nil {
		return err
	}

	return nil
}

func (m MinioStorage) upload(ctx context.Context, path string, file io.Reader, size int64) error {
	_, err := m.client.PutObjectWithContext(ctx, m.bucket, path, file, size, minio.PutObjectOptions{})
	if err != nil {
		return err
	}

	return nil
}

func (m MinioStorage) download(ctx context.Context, path string, w io.Writer) error {
	object, err := m.client.GetObjectWithContext(ctx, m.bucket, path, minio.GetObjectOptions{})
	if err != nil {
		return err
	}

	if _, err = io.Copy(w, object); err != nil {
		return err
	}

	return nil
}
