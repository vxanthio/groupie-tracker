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

const addr = ":8080"

func main() {
	// Load all data from the external API once at startup.
	if err := api.LoadData(); err != nil {
		log.Fatalf("failed to load data: %v", err)
	}

	d := api.GetData()
	s := &store.RealStore{
		Artists:   d.Artists,
		Locations: d.Locations,
		Data:      d.Date,
		Relations: d.Relations,
	}

	homeTmpl := template.Must(template.ParseFiles(
		"web/templates/base.html",
		"web/templates/home.html",
	))

	artistTmpl := template.Must(template.ParseFiles(
		"web/templates/base.html",
		"web/templates/artist.html",
	))

	mux := http.NewServeMux()

	mux.Handle("GET /", handlers.NewHomeHandler(s, homeTmpl))
	mux.Handle("GET /artist/{id}", handlers.NewArtistHandler(s, artistTmpl))

	mux.Handle("GET /static/", http.StripPrefix("/static/", http.FileServer(http.Dir("web/static"))))

	log.Printf("server listening on %s", addr)
	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
