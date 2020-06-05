package dinghy

import (
	"net/http"
)

type Storage interface {
	list(prefix string) (Directory, error)
}

// ServiceServer executes the users requests.
type ServiceServer struct {
	router      *http.ServeMux
	Storage     Storage
	FrontendURL string
}

// NewServiceServer creates a new service server and initiates the routes.
func NewServiceServer() *ServiceServer {
	srv := &ServiceServer{}
	srv.routes()
	return srv
}

func (s *ServiceServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}

func (s *ServiceServer) routes() {
	s.router = http.NewServeMux()
	s.router.HandleFunc("/", s.list)
}
