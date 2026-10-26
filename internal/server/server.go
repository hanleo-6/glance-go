package server

// HTTP server

import (
	_ "embed"
	"html/template"
	"net/http"
	"os"
	"path/filepath"

	"github.com/hanleo-6/glance-go/internal/translate"
)

//go:embed template.html
var templateHTML string


var pageTmpl = template.Must(template.New("page").Parse(templateHTML))

type pageData struct {
	Title string
	Body  template.HTML
}

type Server struct {
	filepath string
	mux      *http.ServeMux
}

func New(filepath string) *Server {
	s := &Server{
		filepath: filepath,
		mux:      http.NewServeMux(),
	}
	s.mux.HandleFunc("/", s.handleIndex)
	return s
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.mux.ServeHTTP(w, r)
}

func (s *Server) handleIndex(w http.ResponseWriter, r *http.Request) {
	rawBytes, err := os.ReadFile(s.filepath)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	parsedHTML, err := translate.Markdown(rawBytes)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	data := pageData{
		Title:  filepath.Base(s.filepath),
		Body:  template.HTML(parsedHTML),
	}

	err = pageTmpl.Execute(w, data)
}

