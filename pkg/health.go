package server

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/minio/minio-go"
)

func newHealth(mc *minio.Client, bucket string) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		timeout, cancel := context.WithTimeout(r.Context(), 1*time.Second)
		defer cancel()
		found, err := bucketExists(timeout, mc, bucket)
		if err != nil {
			log.Printf("healthcheck: failed to check for bucket: %s", err)
		}
		if !found {
			log.Printf("healthcheck: bucket is missing: %s", err)
			w.WriteHeader(http.StatusServiceUnavailable)
			return
		}
		w.WriteHeader(http.StatusOK)
	})
}
