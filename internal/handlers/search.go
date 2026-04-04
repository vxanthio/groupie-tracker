package handlers

import (
	"encoding/json"
	"groupie-tracker/internal/store"
	"net/http"
)

// SearchHandler handles live search requests and holds the store dependency
// needed to query artists. It is registered as a GET-only route and responds
// with JSON, making it suitable for consumption by the frontend search.js script
// without a full page reload.
type SearchHandler struct {
	Store store.Store
}

// Search handles GET /api/search?q= requests by delegating to the store's
// SearchArtists method and encoding the result as a JSON array.
// The q parameter is required — a missing or empty value returns 400 immediately
// so the store is never called with a blank query.
// The response Content-Type is set to application/json before encoding.
// If JSON encoding fails after the header is already written, a 500 is returned,
// though in practice this only occurs if the response writer itself is broken.
func (h *SearchHandler) Search(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	if r.URL.Query().Get("q") == "" {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}
	query := r.URL.Query().Get("q")
	result := h.Store.SearchArtists(query)
	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(result)
	if err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}
}
