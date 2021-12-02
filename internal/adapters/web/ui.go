package web

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"text/template"
)

type UIAdapter struct{}

func NewUIAdapter() *UIAdapter {
	return &UIAdapter{}
}

func NewServer(handler http.Handler, addr string) *http.Server {
	return &http.Server{
		Handler: handler,
		Addr:    addr,
	}
}

func NewPublicHandler(public *UIAdapter, adapter *HTTPAdapter) http.Handler {
	mux := http.NewServeMux()
	mux.Handle("/api/", http.StripPrefix("/api", NewAPIHandler(adapter)))
	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./data/static"))))
	mux.HandleFunc("/", public.index)
	return mux
}

func (u *UIAdapter) index(w http.ResponseWriter, r *http.Request) {
	lp := filepath.Join("data", "templates", "layout.html")
	fp := filepath.Join("data", "templates", "index.html")
	templates, err := template.ParseFiles(lp, fp)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%+v", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	if err := templates.ExecuteTemplate(w, "layout", nil); err != nil {
		fmt.Fprintf(os.Stderr, "%+v", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
}
