package server

import (
	"context"
	"fmt"
	"log"
	"net/url"
	"time"

	"github.com/minio/minio-go"
)

type bucket string
type location string

type minioStorage struct {
	client *minio.Client
	b      bucket
	l      location
}

func NewMinioStorage(c *minio.Client, b, l string) minioStorage {
	s := minioStorage{
		client: c,
		b:      bucket(b),
		l:      location(l),
	}
	return s
}

func (ms minioStorage) bucket() string {
	return string(ms.b)
}

func (ms minioStorage) location() string {
	return string(ms.l)
}

func (ms minioStorage) EnsureBucket() {
	for {
		err := ms.createBucketIfMissing()
		if err != nil {
			log.Printf("failed to ensure bucket exists: %v\n", err)
			time.Sleep(5 * time.Second)
			continue
		}
		return
	}
}

func (ms minioStorage) createBucketIfMissing() error {
	exists, err := ms.client.BucketExists(ms.bucket())
	if err != nil {
		return fmt.Errorf("failed to access bucket %s: %s", ms.bucket(), err)
	}
	if exists {
		return nil
	}

	err = ms.client.MakeBucket(ms.bucket(), ms.location())
	if err != nil {
		return fmt.Errorf("failed to create bucket %s: %s", ms.bucket(), err)
	}
	log.Printf("bucket %s created\n", ms.bucket())

	return nil
}

func (mc minioStorage) healthy(ctx context.Context) error {
	found, err := mc.bucketExists(ctx)
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

func (ms minioStorage) bucketExists(ctx context.Context) (bool, error) {

	res := make(chan exists)

	go func() {
		found, err := ms.client.BucketExists(ms.bucket())
		res <- exists{found: found, err: err}
	}()

	select {
	case res := <-res:
		return res.found, res.err
	case <-ctx.Done():
		return false, ctx.Err()
	}
}

func (ms minioStorage) exists(ctx context.Context, objectName string) (bool, error) {
	if objectName == "" {
		return false, nil
	}

	res := make(chan exists)

	go func() {
		_, err := ms.client.StatObject(ms.bucket(), objectName, minio.StatObjectOptions{})
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

func (ms minioStorage) get(objectName string) (*url.URL, error) {
	return ms.client.PresignedGetObject(ms.bucket(), objectName, 10*time.Minute, url.Values{})
}

func (ms minioStorage) head(objectName string) (*url.URL, error) {
	return ms.client.PresignedHeadObject(ms.bucket(), objectName, 10*time.Minute, url.Values{})
}

func (ms minioStorage) put(objectName string) (*url.URL, error) {
	return ms.client.PresignedPutObject(ms.bucket(), objectName, 10*time.Minute)
}
