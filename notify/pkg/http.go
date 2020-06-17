package notify

import (
	"log"
	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// Server handles http and grpc requests.
type Server struct {
	router *http.ServeMux
}

// NewServer creates a new server.
func NewServer() *Server {
	srv := &Server{}
	srv.routes()
	return srv
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}

func (s *Server) routes() {
	s.router = http.NewServeMux()
	s.router.Handle("/healthz", s.handleHealthz())
	s.router.Handle("/metrics", promhttp.Handler())
	s.router.Handle("/webhook", s.webhook())
}

func (s *Server) handleHealthz() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// executing this means all is OK
		log.Println("healthz")
	}
}

func (s *Server) webhook() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println("webhook")
	}
}
