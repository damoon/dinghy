package server

import (
	"context"

	"github.com/minio/minio-go"
)

func bucketExists(ctx context.Context, mc *minio.Client, bucket string) (bool, error) {
	type bucketExists struct {
		found bool
		err   error
	}
	res := make(chan bucketExists)

	go func() {
		exists, err := mc.BucketExists(bucket)
		res <- bucketExists{found: exists, err: err}
	}()

	select {
	case res := <-res:
		return res.found, res.err
	case <-ctx.Done():
		return false, ctx.Err()
	}
}

func objectExists(ctx context.Context, mc *minio.Client, bucket, object string) (bool, error) {
	if object == "" {
		return false, nil
	}

	type exists struct {
		found bool
		err   error
	}
	res := make(chan exists)

	go func() {
		_, err := mc.StatObject(bucket, object, minio.StatObjectOptions{})
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
