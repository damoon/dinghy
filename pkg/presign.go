package server

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"
)

type presignStorage interface {
	exists(ctx context.Context, objectName string) (bool, error)
	get(objectName string) (*url.URL, error)
	head(objectName string) (*url.URL, error)
	put(objectName string) (*url.URL, error)
}

func NewPresignHandler(storage presignStorage, redirectURL string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		objectName := strings.TrimPrefix(r.RequestURI, "/")
		log.Printf("objectName: %v", objectName)

		if r.Method == http.MethodGet {
			found, err := storage.exists(r.Context(), objectName)
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

		url, err := presignRequest(storage, r.Method, objectName)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Printf("presigning request: %v", err)
			return
		}

		log.Printf("redirect to; %v", url.String())
		http.Redirect(w, r, url.String(), http.StatusTemporaryRedirect)
	}
}

func presignRequest(storage presignStorage, method, objectName string) (*url.URL, error) {
	switch method {
	case http.MethodGet:
		return storage.get(objectName)
	case http.MethodHead:
		return storage.head(objectName)
	case http.MethodPut:
		return storage.put(objectName)
	}
	return nil, fmt.Errorf("method %v not known", method)
}
