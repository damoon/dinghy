package server

import (
	"net/http"

	"github.com/minio/minio-go"
)

func emulate(mc *minio.Client, bucket, redirectURL string) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotImplemented)
	}
}
