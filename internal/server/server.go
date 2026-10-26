package server

// HTTP server

import (
	_ "embed"
	"fmt"
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
	reloadCh chan struct{}
	mux      *http.ServeMux
}

func New(filepath string, reloadCh chan struct{}) *Server {
	s := &Server{
		filepath: filepath,
		reloadCh: reloadCh,
		mux:      http.NewServeMux(),
	}
	s.mux.HandleFunc("/", s.handleIndex)
	s.mux.HandleFunc("/events", s.handleEvents)
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

func (s *Server) handleEvents(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/event-stream")
    w.Header().Set("Cache-Control", "no-cache")
    w.Header().Set("Connection", "keep-alive")
	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Streaming unsupported", http.StatusInternalServerError)
		return
	}
	for {
		select {
		case <-r.Context().Done():
			// Client closed the tab
			return
		case <-s.reloadCh:
			// File watcher notification
			fmt.Println("notified")
			fmt.Fprintf(w, "data: reload\n\n")
			flusher.Flush()
		}
	}
}

