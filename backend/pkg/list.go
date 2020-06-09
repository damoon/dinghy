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
	exists(ctx context.Context, path string) (bool, string, error)
	presign(ctx context.Context, method, path string) (string, error)
}

func (s *ServiceServer) get(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	ctx := r.Context()

	found, etag, err := s.Storage.exists(ctx, filesDirectory+path)
	if err != nil {
		log.Printf("GET %s: %v", path, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if found && path != "/" {
		err = s.download(ctx, etag, w, r)
		if err != nil {
			log.Printf("GET %s: %v", path, err)
			w.WriteHeader(http.StatusInternalServerError)
		}
		return
	}

	if strings.HasSuffix(path, "/") {
		err = s.list(w, r)
		if err != nil {
			log.Printf("GET %s: %v", path, err)
			w.WriteHeader(http.StatusInternalServerError)
		}
		return
	}

	w.WriteHeader(http.StatusNotFound)
}

func (s *ServiceServer) download(ctx context.Context, etag string, w http.ResponseWriter, r *http.Request) error {
	path := filesDirectory + r.URL.Path

	redirect, thumbnail, err := parseRequest(r.URL.RawQuery)
	if err != nil {
		return fmt.Errorf("GET %s: parse parameters: %v", path, err)
	}

	if thumbnail {
		path, err = s.prepareThumbnail(ctx, etag, r.URL.Path)
		if err != nil {
			return fmt.Errorf("GET %s: prepare thumbnail: %v", path, err)
		}
	}

	if redirect {
		url, err := s.Storage.presign(r.Context(), http.MethodGet, path)
		if err != nil {
			return fmt.Errorf("GET %s: presign: %v", path, err)
		}

		http.Redirect(w, r, url, http.StatusTemporaryRedirect)
		return nil
	}

	err = s.delieverFile(r.Context(), path, w)
	if err != nil {
		return fmt.Errorf("GET %s: deliever file: %v", path, err)
	}

	return nil
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

func parseRequest(rawQuery string) (bool, bool, error) {
	m, err := url.ParseQuery(rawQuery)
	if err != nil {
		return false, false, fmt.Errorf("parse url parameters: %v", err)
	}

	_, redirect := m["redirect"]
	_, thumbnail := m["thumbnail"]

	return redirect, thumbnail, nil
}

func (s *ServiceServer) delete(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path

	redirect, _, err := parseRequest(r.URL.RawQuery)
	if err != nil {
		log.Printf("DELETE %s: check redirect: %v", path, err)
	}

	if redirect {
		url, err := s.Storage.presign(r.Context(), http.MethodDelete, filesDirectory+path)
		if err != nil {
			log.Printf("DELETE %s: redirect: %v", path, err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, url, http.StatusTemporaryRedirect)
		return
	}

	err = s.Storage.delete(r.Context(), filesDirectory+path)
	if err != nil {
		log.Printf("DELETE %s: %v", path, err)
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func (s *ServiceServer) put(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path

	redirect, _, err := parseRequest(r.URL.RawQuery)
	if err != nil {
		log.Printf("PUT %s: check redirect: %v", path, err)
	}

	if redirect {
		url, err := s.Storage.presign(r.Context(), http.MethodPut, filesDirectory+path)
		if err != nil {
			log.Printf("PUT %s: redirect: %v", path, err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, url, http.StatusTemporaryRedirect)
		return
	}

	err = s.receiveFile(r.Context(), filesDirectory+path, r.Body)
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

func (s *ServiceServer) list(w http.ResponseWriter, r *http.Request) error {
	path := r.URL.Path

	l, err := s.Storage.list(r.Context(), path)
	if err != nil {
		return fmt.Errorf("list %s: %v", path, err)
	}

	err = respond(w, r, l, s.FrontendURL)
	if err != nil {
		return fmt.Errorf("respond: %v", err)
	}

	return nil
}

func respond(w http.ResponseWriter, r *http.Request, l Directory, frontendURL string) error {
	if requestsJSON(r.Header.Get("Accept")) {
		err := l.toJSON(w)
		if err != nil {
			return fmt.Errorf("render json: %v", err)
		}

		return nil
	}

	if isCLIClient(r.UserAgent()) {
		err := l.toTXT(w)
		if err != nil {
			return fmt.Errorf("render text: %v", err)
		}

		return nil
	}

	http.Redirect(w, r, frontendURL+r.URL.Path, http.StatusTemporaryRedirect)

	return nil
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
