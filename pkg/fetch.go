package server

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/minio/minio-go"
)

func fetch(mc *minio.Client, bucket, redirectURL string) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		objectName := strings.TrimPrefix(r.RequestURI, "/")
		log.Printf("objectName: %v", objectName)

		if r.Method == http.MethodGet {
			found, err := objectExists(r.Context(), mc, bucket, objectName)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				log.Printf("looking up file: %v", err)
				return
			}
			if !found {
				http.Redirect(w, r, redirectURL, http.StatusTemporaryRedirect)
				return
			}
		}

		url, err := presignRequest(mc, r.Method, bucket, objectName)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Printf("presigning request: %v", err)
			return
		}

		log.Printf("redirect to; %v", url.String())
		http.Redirect(w, r, url.String(), http.StatusTemporaryRedirect)
	}
}

func presignRequest(mc *minio.Client, method, bucket, objectName string) (*url.URL, error) {
	switch method {
	case http.MethodGet:
		return mc.PresignedGetObject(bucket, objectName, 10*time.Minute, url.Values{})
	case http.MethodHead:
		return mc.PresignedHeadObject(bucket, objectName, 10*time.Minute, url.Values{})
	case http.MethodPut:
		return mc.PresignedPutObject(bucket, objectName, 10*time.Minute)
	}
	return nil, fmt.Errorf("method %v not known", method)
}
