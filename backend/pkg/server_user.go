package dinghy

import (
	"log"
	"net/http"
)

// ServiceServer executes the users requests.
type ServiceServer struct {
	Storage     ObjectStore
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
