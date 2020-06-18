package dinghy

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	"gitlab.com/davedamoon/dinghy/backend/pkg/pb"
)

// ServiceServer executes the users requests.
type ServiceServer struct {
	Storage        ObjectStore
	FrontendURL    string
	Upgrader       websocket.Upgrader
	NotifierClient pb.NotifierClient
}

// NewServiceServer creates a new service server and initiates the routes.
func NewServiceServer() *ServiceServer {
	srv := &ServiceServer{}
	return srv
}

func (s *ServiceServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
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
		log.Printf("%s %s not supported", r.Method, r.URL.Path)
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}
