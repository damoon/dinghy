package dinghy

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/minio/minio-go"
)

type MinioStorage struct {
	client   *minio.Client
	bucket   string
	location string
}

func NewMinioStorage(c *minio.Client, b, l string) MinioStorage {
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
