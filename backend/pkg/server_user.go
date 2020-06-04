package dinghy

import (
	"net/http"
	"time"
)

// ServiceServer executes the users requests.
type ServiceServer struct {
	router *http.ServeMux
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
	svc := &svc{}
	s.router.Handle("/", http.TimeoutHandler(svc, 30*time.Second, ""))
}
