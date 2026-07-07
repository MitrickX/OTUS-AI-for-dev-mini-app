package main

import (
	"log"
	"net/http"

	"github.com/dmitrypavlov/mini-questionnaire/api"
	"github.com/dmitrypavlov/mini-questionnaire/internal/handler"
)

func main() {
	h := handler.New()

	mux := http.NewServeMux()
	api.HandlerFromMux(h, mux)

	s := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	log.Printf("Starting server on %s", s.Addr)
	log.Fatal(s.ListenAndServe())
}
