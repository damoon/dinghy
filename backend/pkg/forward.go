package server

import (
	"net/http"
)

type forwardStorage interface {
	//	exists(ctx context.Context, objectName string) (bool, error)
	//	store(ctx context.Context, objectName string, content io.Reader)
	//	restore(ctx context.Context, objectName string, content io.Writer)
	//	delete(ctx context.Context, objectName string) error
	//	list(ctx context.Context, path string) ([]string, error)
}

func NewForwardHandler(storage forwardStorage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotImplemented)
	}
}
