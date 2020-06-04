package dinghy

import (
	"context"
	"log"
	"net/http"
	"time"

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
		// 1 second is the default timeout for readiness probes https://kubernetes.io/docs/tasks/configure-pod-container/configure-liveness-readiness-startup-probes/#configure-probes.
		ctx, cancel := context.WithTimeout(r.Context(), 1*time.Second)
		defer cancel()

		err := s.Storage.healthy(ctx)
		if err != nil {
			log.Printf("healthcheck: %s", err)
			w.WriteHeader(http.StatusServiceUnavailable)
		}
	}
}
