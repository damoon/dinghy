package dinghy

import (
	"context"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
)

type ObjectStore interface {
	list(ctx context.Context, prefix string) (Directory, error)
	upload(ctx context.Context, path string, file io.ReadSeeker) error
	delete(ctx context.Context, path string) error
	download(ctx context.Context, path string, w io.WriterAt) error
	exists(ctx context.Context, path string) (bool, error)
	presign(method, path string) (string, error)
}

func (s *ServiceServer) get(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path

	found, err := s.Storage.exists(r.Context(), path)
	if err != nil {
		log.Printf("GET %s: %v", path, err)
		w.WriteHeader(http.StatusInternalServerError)

		return
	}

	if found && path != "/" {
		s.download(w, r)
		return
	}

	if strings.HasSuffix(path, "/") {
		s.list(w, r)
		return
	}

	w.WriteHeader(http.StatusNotFound)
}

func (s *ServiceServer) download(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path

	redirect, err := shouldRedirect(r.URL.RawQuery)
	if err != nil {
		log.Printf("GET %s: check redirect: %v", path, err)
	}

	if redirect {
		url, err := s.Storage.presign(http.MethodGet, path)
		if err != nil {
			log.Printf("GET %s: redirect: %v", path, err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, url, http.StatusTemporaryRedirect)
		return
	}

	err = s.delieverFile(r.Context(), path, w)
	if err != nil {
		log.Printf("GET %s: send object: %v", path, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (s *ServiceServer) delieverFile(ctx context.Context, path string, w io.Writer) error {
	tmpfile, err := ioutil.TempFile("", "s3_download")
	if err != nil {
		return fmt.Errorf("create temp file: %v", err)
	}
	defer os.Remove(tmpfile.Name())

	err = s.Storage.download(ctx, path, tmpfile)
	if err != nil {
		return fmt.Errorf("download: %v", err)
	}

	_, err = tmpfile.Seek(0, 0)
	if err != nil {
		return fmt.Errorf("seek temp file: %v", err)
	}

	_, err = io.Copy(w, tmpfile)
	if err != nil {
		return fmt.Errorf("write reponse: %v", err)
	}

	return nil
}

func shouldRedirect(rawQuery string) (bool, error) {
	m, err := url.ParseQuery(rawQuery)
	if err != nil {
		return false, fmt.Errorf("parse url parameters: %v", err)
	}
	_, ok := m["redirect"]
	if ok {
		return true, nil
	}

	return false, nil
}

func (s *ServiceServer) delete(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path

	redirect, err := shouldRedirect(r.URL.RawQuery)
	if err != nil {
		log.Printf("DELETE %s: check redirect: %v", path, err)
	}

	if redirect {
		url, err := s.Storage.presign(http.MethodDelete, path)
		if err != nil {
			log.Printf("DELETE %s: redirect: %v", path, err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, url, http.StatusTemporaryRedirect)
		return
	}

	err = s.Storage.delete(r.Context(), path)
	if err != nil {
		log.Printf("DELETE %s: %v", path, err)
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func (s *ServiceServer) put(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path

	redirect, err := shouldRedirect(r.URL.RawQuery)
	if err != nil {
		log.Printf("PUT %s: check redirect: %v", path, err)
	}

	if redirect {
		url, err := s.Storage.presign(http.MethodPut, path)
		if err != nil {
			log.Printf("PUT %s: redirect: %v", path, err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, url, http.StatusTemporaryRedirect)
		return
	}

	err = s.receiveFile(r.Context(), path, r.Body)
	if err != nil {
		log.Printf("PUT %s: receive file: %v", path, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (s *ServiceServer) receiveFile(ctx context.Context, path string, r io.Reader) error {
	tmpfile, err := ioutil.TempFile("", "s3_upload")
	if err != nil {
		return fmt.Errorf("create temp file: %v", err)
	}
	defer os.Remove(tmpfile.Name())

	_, err = io.Copy(tmpfile, r)
	if err != nil {
		return fmt.Errorf("write local temp file: %v", err)
	}

	_, err = tmpfile.Seek(0, 0)
	if err != nil {
		return fmt.Errorf("seek temp file: %v", err)
	}

	err = s.Storage.upload(ctx, path, tmpfile)
	if err != nil {
		return fmt.Errorf("upload: %v", err)
	}

	return nil
}

func (s *ServiceServer) list(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path

	l, err := s.Storage.list(r.Context(), path)
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
