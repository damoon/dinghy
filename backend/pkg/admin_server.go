package dinghy

import (
	"context"
	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// Healthy allows to check the health of a service.
type Healthy interface {
	healthy(context.Context) error
}

// AdminServer answers to administration requests.
type AdminServer struct {
	router  *http.ServeMux
	Storage Healthy
}

// NewAdminServer creates a new administration server.
func NewAdminServer() *AdminServer {
	srv := &AdminServer{}
	srv.routes()

	return srv
}

func (s *AdminServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}

func (s *AdminServer) routes() {
	s.router = http.NewServeMux()
	s.router.HandleFunc("/healthz", s.handleHealthz())
	s.router.Handle("/metrics", promhttp.Handler())
}

func (s *AdminServer) handleHealthz() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := s.Storage.healthy(r.Context())
		if err != nil {
			w.WriteHeader(http.StatusServiceUnavailable)
		}
	}
}
