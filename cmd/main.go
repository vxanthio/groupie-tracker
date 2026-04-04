// Package main is the entry point for the groupie-tracker web server.
// It loads all artist data from the external API on startup, builds the
// in-memory store, registers all HTTP routes, and starts listening.
package main

import (
	"html/template"
	"log"
	"net/http"

	"groupie-tracker/internal/api"
	"groupie-tracker/internal/handlers"
	"groupie-tracker/internal/store"
)

// addr is the TCP address the HTTP server listens on.
const addr = ":8080"

// main loads all data from the external API, wires up the dependency graph
// (store, templates, handlers, middleware), and starts the HTTP server.
func main() {
	// Load all data from the external API once at startup.
	if err := api.LoadData(); err != nil {
		log.Fatalf("failed to load data: %v", err)
	}

	d := api.GetData()
	s := &store.RealStore{
		Artists:   d.Artists,
		Locations: d.Locations,
		Dates:     d.Dates,
		Relations: d.Relations,
	}
	log.Printf("Data loaded: %d artists", len(d.Artists))
	homeTmpl := template.Must(template.ParseFiles(
		"web/templates/base.html",
		"web/templates/home.html",
	))

	artistTmpl := template.Must(template.ParseFiles(
		"web/templates/base.html",
		"web/templates/artist.html",
	))
	errorTmpl := template.Must(template.ParseGlob("web/templates/*.html"))
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			handlers.NotFoundHandler(errorTmpl)(w, r)
			return
		}
		handlers.NewHomeHandler(s, homeTmpl).ServeHTTP(w, r)
	})
	mux.Handle("GET /artist/{id}", handlers.NewArtistHandler(s, artistTmpl))
	mux.Handle("GET /static/", http.StripPrefix("/static/", http.FileServer(http.Dir("web/static"))))
	searchHandler := &handlers.SearchHandler{Store: s}
	mux.HandleFunc("GET /api/search", searchHandler.Search)

	log.Printf("server listening on %s", addr)
	if err := http.ListenAndServe(addr, handlers.RecoveryMiddleware(errorTmpl, mux)); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
