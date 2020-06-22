package notify

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/pprof"
	"sync"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// Server handles http requests.
type Server struct {
	BearerToken string
	router      *http.ServeMux
	C           *sync.Cond
}

// NewServer creates a new http server.
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
	s.router.Handle("/metrics", promhttp.Handler())
	s.router.Handle("/webhook", s.webhook())
	s.router.HandleFunc("/debug/pprof/", pprof.Index)
	s.router.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
	s.router.HandleFunc("/debug/pprof/profile", pprof.Profile)
	s.router.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
	s.router.HandleFunc("/debug/pprof/trace", pprof.Trace)
}

func (s *Server) webhook() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !s.isValidMinioRequest(r) {
			return
		}

		t, err := eventType(r.Body)
		if err != nil {
			log.Printf("get event type: %v", err)
			w.WriteHeader(http.StatusBadRequest)

			return
		}

		if !meansChange(t) {
			return
		}

		s.C.L.Lock()
		s.C.Broadcast()
		s.C.L.Unlock()
	}
}

func (s *Server) isValidMinioRequest(r *http.Request) bool {
	return r.Header.Get("Authorization") == "Bearer "+s.BearerToken
}

type minioNotification struct {
	EventName string
}

func eventType(r io.Reader) (string, error) {
	notification := &minioNotification{}

	err := json.NewDecoder(r).Decode(notification)
	if err != nil {
		return "", fmt.Errorf("get event type: %v", err)
	}

	return notification.EventName, nil
}

func meansChange(event string) bool {
	ss := []string{"s3:ObjectRemoved:Delete", "s3:ObjectCreated:Put", "s3:ObjectCreated:Copy"}
	for _, s := range ss {
		if event == s {
			return true
		}
	}

	return false
}
