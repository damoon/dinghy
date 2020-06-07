package dinghy

import (
	"context"
	"encoding/json"
	"html/template"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

type ObjectStore interface {
	list(ctx context.Context, prefix string) (Directory, error)
	upload(ctx context.Context, path string, file io.ReadSeeker) error
	delete(ctx context.Context, path string) error
	download(ctx context.Context, path string, w io.WriterAt) error
	exists(ctx context.Context, path string) (bool, error)
}

func (s *ServiceServer) get(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path

	ctx, cancel := context.WithTimeout(r.Context(), 30*time.Second)
	defer cancel()

	found, err := s.Storage.exists(ctx, path)
	if err != nil {
		log.Printf("GET %s: %v", path, err)
		w.WriteHeader(http.StatusInternalServerError)

		return
	}

	if found && path != "/" {
		s.download(w, r)
		return
	}

	s.list(w, r)
}

func (s *ServiceServer) download(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path

	ctx, cancel := context.WithTimeout(r.Context(), 30*time.Second)
	defer cancel()

	tmpfile, err := ioutil.TempFile("", "s3_download")
	if err != nil {
		log.Printf("GET %s: create temp file: %v", path, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer os.Remove(tmpfile.Name())

	err = s.Storage.download(ctx, path, tmpfile)
	if err != nil {
		log.Printf("GET %s: download: %v", path, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	_, err = tmpfile.Seek(0, 0)
	if err != nil {
		log.Printf("GET %s: seek temp file: %v", path, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	_, err = io.Copy(w, tmpfile)
	if err != nil {
		log.Printf("GET %s: write reponse: %v", path, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (s *ServiceServer) delete(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	err := s.Storage.delete(ctx, path)
	if err != nil {
		log.Printf("DELETE %s: %v", path, err)
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func (s *ServiceServer) put(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path

	ctx, cancel := context.WithTimeout(r.Context(), 30*time.Second)
	defer cancel()

	defer r.Body.Close()

	tmpfile, err := ioutil.TempFile("", "s3_upload")
	if err != nil {
		log.Printf("PUT %s: create temp file: %v", path, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer os.Remove(tmpfile.Name())

	_, err = io.Copy(tmpfile, r.Body)
	if err != nil {
		log.Printf("PUT %s: write local temp file: %v", path, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	_, err = tmpfile.Seek(0, 0)
	if err != nil {
		log.Printf("PUT %s: seek temp file: %v", path, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = s.Storage.upload(ctx, path, tmpfile)
	if err != nil {
		log.Printf("PUT %s: upload: %v", path, err)
		w.WriteHeader(http.StatusInternalServerError)

		return
	}
}

func (s *ServiceServer) list(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path

	setupCORS(&w, r, s.FrontendURL)

	ctx, cancel := context.WithTimeout(r.Context(), 30*time.Second)
	defer cancel()

	l, err := s.Storage.list(ctx, path)
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
