package dinghy

import (
	"context"
	"net/http"
	"net/http/pprof"

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
	s.router.HandleFunc("/healthz", handleHealthz)
	s.router.Handle("/metrics", promhttp.Handler())
	s.router.HandleFunc("/debug/pprof/", pprof.Index)
	s.router.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
	s.router.HandleFunc("/debug/pprof/profile", pprof.Profile)
	s.router.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
	s.router.HandleFunc("/debug/pprof/trace", pprof.Trace)
}

func handleHealthz(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}
