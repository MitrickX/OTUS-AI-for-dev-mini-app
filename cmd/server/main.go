package main

import (
	"embed"
	"io/fs"
	"log"
	"net/http"

	"github.com/dmitrypavlov/mini-questionnaire/api"
	"github.com/dmitrypavlov/mini-questionnaire/internal/handler"
)

//go:embed static
var staticFiles embed.FS

func main() {
	h := handler.New()

	mux := http.NewServeMux()

	api.HandlerFromMux(h, mux)

	staticSub, err := fs.Sub(staticFiles, "static")
	if err != nil {
		log.Fatal(err)
	}
	fileServer := http.FileServer(http.FS(staticSub))

	mux.Handle("GET /static/", http.StripPrefix("/static/", fileServer))

	mux.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}
		data, err := staticFiles.ReadFile("static/index.html")
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Write(data)
	})

	s := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	log.Printf("Starting server on %s", s.Addr)
	log.Fatal(s.ListenAndServe())
}
