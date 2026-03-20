// Package handlers provides tests for all HTTP request handlers.
package handlers

import (
	"html/template"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"groupie-tracker/internal/models"
	"groupie-tracker/internal/store"
)

// --- Mock store for handler tests ---

// testStore is a simple store.Store implementation for use in handler tests.
type testStore struct {
	artists []models.Artist
}

func (s *testStore) AllArtists() []models.Artist {
	return s.artists
}

func (s *testStore) ArtistByID(id int) (models.Artist, bool) {
	for _, a := range s.artists {
		if a.ID == id {
			return a, true
		}
	}
	return models.Artist{}, false
}

func (s *testStore) SearchArtists(query string) []models.Artist {
	return s.artists
}

// Compile-time check: testStore satisfies store.Store.
var _ store.Store = (*testStore)(nil)

// --- Template helpers for tests ---

// mustParseTemplate parses a template string and panics on error.
// Only used in test setup — panics are acceptable there.
func mustParseTemplate(src string) *template.Template {
	return template.Must(template.New("base").Parse(src))
}

// brokenTemplate returns a template that always fails on execution.
func brokenTemplate() *template.Template {
	// Calling a nil value forces an execution error without panicking the handler.
	tmpl, _ := template.New("base").Parse(`{{call .}}`)
	return tmpl
}

// --- GET / (HomeHandler) tests ---

func TestHomeHandler(t *testing.T) {
	twoArtists := []models.Artist{
		{ID: 1, Name: "Foo Fighters", Image: "http://img/1.jpg", CreationDate: 1994},
		{ID: 2, Name: "Queen", Image: "http://img/2.jpg", CreationDate: 1970},
	}

	// Minimal template that renders artist names — mirrors what home.html does.
	homeTmpl := mustParseTemplate(`{{range .}}{{.Name}}{{end}}`)

	tests := []struct {
		name             string
		artists          []models.Artist
		tmpl             *template.Template
		wantStatusCode   int
		wantBodyContains []string
	}{
		{
			name:             "happy_path_two_artists",
			artists:          twoArtists,
			tmpl:             homeTmpl,
			wantStatusCode:   http.StatusOK,
			wantBodyContains: []string{"Foo Fighters", "Queen"},
		},
		{
			name:             "empty_store_returns_200",
			artists:          []models.Artist{},
			tmpl:             homeTmpl,
			wantStatusCode:   http.StatusOK,
			wantBodyContains: []string{},
		},
		{
			name:             "template_error_returns_500",
			artists:          twoArtists,
			tmpl:             brokenTemplate(),
			wantStatusCode:   http.StatusInternalServerError,
			wantBodyContains: []string{},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			s := &testStore{artists: tc.artists}
			h := NewHomeHandler(s, tc.tmpl)

			req := httptest.NewRequest(http.MethodGet, "/", nil)
			rec := httptest.NewRecorder()

			h.ServeHTTP(rec, req)

			if rec.Code != tc.wantStatusCode {
				t.Errorf("status = %d, want %d", rec.Code, tc.wantStatusCode)
			}

			body := rec.Body.String()
			for _, want := range tc.wantBodyContains {
				if !strings.Contains(body, want) {
					t.Errorf("body does not contain %q\nbody: %s", want, body)
				}
			}
		})
	}
}

// --- GET /artist/{id} (ArtistHandler) tests ---

func TestArtistHandler(t *testing.T) {
	artists := []models.Artist{
		{ID: 1, Name: "Foo Fighters", Image: "http://img/1.jpg", CreationDate: 1994, FirstAlbum: "04-07-1995"},
	}

	// Minimal template that renders the artist name from ArtistPageData.
	artistTmpl := mustParseTemplate(`{{.Artist.Name}}`)

	tests := []struct {
		name             string
		url              string
		pathID           string
		artists          []models.Artist
		tmpl             *template.Template
		wantStatusCode   int
		wantBodyContains []string
	}{
		{
			name:             "valid_id_returns_200",
			url:              "/artist/1",
			pathID:           "1",
			artists:          artists,
			tmpl:             artistTmpl,
			wantStatusCode:   http.StatusOK,
			wantBodyContains: []string{"Foo Fighters"},
		},
		{
			name:             "unknown_id_returns_404",
			url:              "/artist/99",
			pathID:           "99",
			artists:          artists,
			tmpl:             artistTmpl,
			wantStatusCode:   http.StatusNotFound,
			wantBodyContains: []string{},
		},
		{
			name:             "non_numeric_id_returns_404",
			url:              "/artist/abc",
			pathID:           "abc",
			artists:          artists,
			tmpl:             artistTmpl,
			wantStatusCode:   http.StatusNotFound,
			wantBodyContains: []string{},
		},
		{
			name:             "template_error_returns_500",
			url:              "/artist/1",
			pathID:           "1",
			artists:          artists,
			tmpl:             brokenTemplate(),
			wantStatusCode:   http.StatusInternalServerError,
			wantBodyContains: []string{},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			s := &testStore{artists: tc.artists}
			h := NewArtistHandler(s, tc.tmpl)

			req := httptest.NewRequest(http.MethodGet, tc.url, nil)
			req.SetPathValue("id", tc.pathID)
			rec := httptest.NewRecorder()

			h.ServeHTTP(rec, req)

			if rec.Code != tc.wantStatusCode {
				t.Errorf("status = %d, want %d", rec.Code, tc.wantStatusCode)
			}

			body := rec.Body.String()
			for _, want := range tc.wantBodyContains {
				if !strings.Contains(body, want) {
					t.Errorf("body does not contain %q\nbody: %s", want, body)
				}
			}
		})
	}
}
