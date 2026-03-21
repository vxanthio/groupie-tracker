package handlers

import (
	"encoding/json"
	"groupie-tracker/internal/store"
	"net/http"
)

// SearchHandler holds the store dependency for search operations.
type SearchHandler struct {
	Store store.Store
}

// Search handles GET /api/search?q= requests.
// It filters artists by the query string and returns a JSON array.
// Returns 405 if method is not GET.
// Returns 400 if query parameter is missing or empty.
// Returns 500 if JSON encoding fails.
func (h *SearchHandler) Search(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	if r.URL.Query().Get("q") == "" {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}
	query := r.URL.Query().Get("q")
	result := h.Store.SearchArtists(query)
	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(result)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}
