package dinghy

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

// ServiceServer executes the users requests.
type ServiceServer struct {
	Storage     ObjectStore
	FrontendURL string
	upgrader    websocket.Upgrader
}

// NewServiceServer creates a new service server and initiates the routes.
func NewServiceServer() *ServiceServer {
	srv := &ServiceServer{}
	srv.upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}
	srv.upgrader.CheckOrigin = func(r *http.Request) bool {
		if r.Header.Get("Origin") != srv.FrontendURL {
			return false
		}

		return true
	}

	return srv
}

func (s *ServiceServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodOptions:
		return
	case http.MethodGet:
		s.get(w, r)
	case http.MethodPost:
		s.post(w, r)
	case http.MethodPut:
		s.put(w, r)
	case http.MethodDelete:
		s.delete(w, r)
	default:
		log.Printf("%s %s not supported", r.Method, r.URL.Path)
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}
