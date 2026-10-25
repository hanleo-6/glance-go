package server

// HTTP server

import (
	"github.com/hanleo-6/glance-go/internal/translate"

	_ "embed"

	"path/filepath"
	"html/template"
	"net/http"
	"os"
)

//go:embed template.html
var templateHTML string

type pageData struct {
	Title string
	Body  template.HTML
}

func NewPageData(title string, body template.HTML) *pageData {
	return &pageData{
		Title: title,
		Body: body,
	}
}

type Server struct {
	filepath string
	mux      *http.ServeMux
}

func NewServer(filepath string) *Server {
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

	data := NewPageData(filepath.Base(s.filepath), template.HTML(parsedHTML))

	pageTmpl, err := template.New("page").Parse(templateHTML)
	err = pageTmpl.Execute(w, data)
}

