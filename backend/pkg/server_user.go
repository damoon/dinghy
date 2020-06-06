package dinghy

import (
	"context"
	"io"
	"log"
	"net/http"
)

type Storage interface {
	list(prefix string) (Directory, error)
	upload(ctx context.Context, path string, file io.Reader, size int64) error
	delete(ctx context.Context, path string) error
	download(ctx context.Context, path string, w io.Writer) error
	exists(ctx context.Context, path string) (bool, error)
}

// ServiceServer executes the users requests.
type ServiceServer struct {
	Storage     Storage
	FrontendURL string
}

// NewServiceServer creates a new service server and initiates the routes.
func NewServiceServer() *ServiceServer {
	srv := &ServiceServer{}
	return srv
}

func (s *ServiceServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Printf("%s %s", r.Method, r.URL.Path)
	switch r.Method {
	case http.MethodOptions:
		return
	case http.MethodGet:
		s.get(w, r)
	case http.MethodPut:
		s.put(w, r)
	case http.MethodDelete:
		s.delete(w, r)
	default:
		log.Printf("%s %s not allowed", r.URL.Path, r.Method)
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}
