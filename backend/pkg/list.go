package dinghy

import (
	"encoding/json"
	"html/template"
	"io"
	"log"
	"net/http"
	"strings"
)

type Directory struct {
	Path        string
	Directories []string
	Files       []File
}

type File struct {
	Name string
	Size int
	Icon string
}

type svc struct {
}

func (s *svc) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	setupCORS(&w, r)
	if (*r).Method == "OPTIONS" {
		return
	}

	l := Directory{
		Path: r.URL.Path,
		Directories: []string{
			"backups",
			"pictures",
		},
		Files: []File{
			{Name: "apache.log", Size: 1024},
			{Name: "background.png", Size: 10240},
			{Name: "a very long name just to see it still works.png", Size: 10240},
			{Name: "fdshiofdjsaoifdjfiadsjfoidsajfsaoidjiofdsa.tiff", Size: 10240},
			{Name: "more.exe", Size: 10240},
			{Name: "apache.log", Size: 1024},
			{Name: "background.png", Size: 10240},
			{Name: "a very long name just to see it still works.zip", Size: 10240},
			{Name: "fdshiofdjsaoifdjfiadsjfoidsajfsaoidjiofdsa.png", Size: 10240},
			{Name: "more.exe", Size: 10240},
			{Name: "apache.log", Size: 1024},
			{Name: "background.png", Size: 10240},
			{Name: "a very long name just to see it still works.bmp", Size: 10240},
			{Name: "fdshiofdjsaoifdjfiadsjfoidsajfsaoidjiofdsa.tar", Size: 10240},
			{Name: "more.exe", Size: 10240},
			{Name: "apache.log", Size: 1024},
			{Name: "background.png", Size: 10240},
			{Name: "a very long name just to see it still works.jpg", Size: 10240},
			{Name: "fdshiofdjsaoifdjfiadsjfoidsajfsaoidjiofdsa.png", Size: 10240},
			{Name: "more.exe", Size: 10240},
			{Name: "apache.log", Size: 1024},
			{Name: "background.png", Size: 10240},
			{Name: "a very long name just to see it still works.jpeg", Size: 10240},
			{Name: "fdshiofdjsaoifdjfiadsjfoidsajfsaoidjiofdsa.tar.gz", Size: 10240},
			{Name: "more.exe", Size: 10240},
		},
	}

	addIcons(l.Files)

	respond(w, r, l)
}

func respond(w http.ResponseWriter, r *http.Request, l Directory) {
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

	http.Redirect(w, r, "/ui/"+r.URL.Path, http.StatusTemporaryRedirect)
}

func setupCORS(w *http.ResponseWriter, req *http.Request) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
	(*w).Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	(*w).Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
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
