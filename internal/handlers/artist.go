package handlers

import (
	"bytes"
	"html/template"
	"net/http"
	"strconv"

	"groupie-tracker/internal/models"
	"groupie-tracker/internal/store"
)

// ArtistHandler handles GET /artist/{id} requests by rendering the artist detail page.
type ArtistHandler struct {
	store store.Store
	tmpl  *template.Template
}

// NewArtistHandler constructs an ArtistHandler with the given store and template.
// The template is parsed once at construction time and reused across requests.
func NewArtistHandler(s store.Store, tmpl *template.Template) http.Handler {
	return &ArtistHandler{store: s, tmpl: tmpl}
}

// ServeHTTP parses the artist ID from the URL, looks up the artist, and renders
// the detail template. Returns 404 for unknown or non-numeric IDs, 500 on
// template execution failure.
func (h *ArtistHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	artist, ok := h.store.ArtistByID(id)
	if !ok {
		http.NotFound(w, r)
		return
	}

	data := models.ArtistPageData{
		Artist: artist,
		// Locations and Dates are populated once Vasiliki extends the Store interface.
	}

	var buf bytes.Buffer
	if err := h.tmpl.Execute(&buf, data); err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	buf.WriteTo(w) //nolint:errcheck // response write errors are unrecoverable
}
