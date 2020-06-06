package dinghy

import (
	"context"
	"encoding/json"
	"html/template"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func (s *ServiceServer) get(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/")

	ctx, cancel := context.WithTimeout(r.Context(), 30*time.Second)
	defer cancel()

	found, err := s.Storage.exists(ctx, path)
	if err != nil {
		log.Printf("GET %s: %v", path, err)
		w.WriteHeader(http.StatusInternalServerError)

		return
	}

	if found && path != "" {
		s.download(w, r)

		return
	}

	s.list(w, r)
}

func (s *ServiceServer) download(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/")

	ctx, cancel := context.WithTimeout(r.Context(), 30*time.Second)
	defer cancel()

	err := s.Storage.download(ctx, path, w)
	if err != nil {
		log.Printf("GET %s: %v", path, err)
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func (s *ServiceServer) delete(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/")

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	err := s.Storage.delete(ctx, path)
	if err != nil {
		log.Printf("DELETE %s: %v", path, err)
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func (s *ServiceServer) put(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/")

	ctx, cancel := context.WithTimeout(r.Context(), 30*time.Second)
	defer cancel()

	file := r.Body
	defer r.Body.Close()

	contentLength := r.Header.Get("Content-Length")

	size, err := strconv.ParseInt(contentLength, 10, 64)
	if err != nil {
		log.Printf("PUT %s: parse size %s: %v", path, contentLength, err)
		w.WriteHeader(http.StatusInternalServerError)

		return
	}

	err = s.Storage.upload(ctx, path, file, size)
	if err != nil {
		log.Printf("PUT %s: %v", path, err)
		w.WriteHeader(http.StatusInternalServerError)

		return
	}
}

func (s *ServiceServer) list(w http.ResponseWriter, r *http.Request) {
	setupCORS(&w, r, s.FrontendURL)

	path := strings.TrimPrefix(r.URL.Path, "/")

	l, err := s.Storage.list(path)
	if err != nil {
		log.Printf("list %s: %v", path, err)
		w.WriteHeader(http.StatusInternalServerError)

		return
	}

	addIcons(l.Files)

	respond(w, r, l, s.FrontendURL)
}

func respond(w http.ResponseWriter, r *http.Request, l Directory, frontendURL string) {
	if requestsJSON(r.Header.Get("Accept")) {
		err := l.toJSON(w)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Println(err)
		}

		return
	}

	if isCLIClient(r.UserAgent()) {
		err := l.toTXT(w)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Println(err)
		}

		return
	}

	http.Redirect(w, r, frontendURL+r.URL.Path, http.StatusTemporaryRedirect)
}

func setupCORS(w *http.ResponseWriter, req *http.Request, domain string) {
	(*w).Header().Set("Access-Control-Allow-Origin", domain)
}

func requestsJSON(ct string) bool {
	if strings.Contains(strings.ToLower(ct), "application/json") {
		return true
	}

	log.Println("test")

	return false
}

func isCLIClient(agent string) bool {
	if strings.Contains(strings.ToLower(agent), "curl") {
		return true
	}

	if strings.Contains(strings.ToLower(agent), "wget") {
		return true
	}

	return false
}

func (l Directory) toJSON(w io.Writer) error {
	e := json.NewEncoder(w)

	err := e.Encode(l)
	if err != nil {
		return err
	}

	return nil
}

func (l Directory) toTXT(w io.Writer) error {
	const letter = `{{.Path}}:
{{range .Directories}}{{ . }}/
{{end}}{{range .Files}}{{ .Name }} ({{ .Size }} Byte)
{{end}}`

	t := template.Must(template.New("letter").Parse(letter))

	err := t.Execute(w, l)
	if err != nil {
		return err
	}

	return nil
}
