package store

import (
	"groupie-tracker/internal/models"
	"testing"
)

func TestMockStore_AllArtists(t *testing.T) {
	m := &MockStore{}
	artists := m.AllArtists()
	if len(artists) != 2 {
		t.Errorf("expected 2 artists, got %d", len(artists))
	}
}

func TestMockStore_ArtistByID_Found(t *testing.T) {
	m := &MockStore{}
	a, ok := m.ArtistByID(1)
	if !ok {
		t.Fatal("expected artist to be found")
	}
	if a.Name != "Billie Eilish" {
		t.Errorf("expected Billie Eilish, got %s", a.Name)
	}
}

func TestMockStore_ArtistByID_NotFound(t *testing.T) {
	m := &MockStore{}
	_, ok := m.ArtistByID(99)
	if ok {
		t.Error("expected false for missing ID")
	}
}

func TestMockStore_SearchArtists(t *testing.T) {
	m := &MockStore{}

	results := m.SearchArtists("Billie")
	if len(results) != 1 {
		t.Errorf("expected 1 result, got %d", len(results))
	}

	results = m.SearchArtists("Queen")
	if len(results) != 0 {
		t.Errorf("expected 0 results, got %d", len(results))
	}

	results = m.SearchArtists("system")
	if len(results) != 1 {
		t.Errorf("expected 1 result, got %d", len(results))
	}
}
func TestRealStore_SearchArtists(t *testing.T) {
	s := &RealStore{
		Artists: []models.Artist{
			{ID: 1, Name: "Queen", Members: []string{"Freddie Mercury", "Brian May"}, CreationDate: 1970, FirstAlbum: "14-12-1973"},
			{ID: 2, Name: "Billie Eilish", Members: []string{"Billie Eilish"}, CreationDate: 2015, FirstAlbum: "26-03-2017"},
		},
	}

	tests := []struct {
		name      string
		query     string
		wantNames []string
	}{
		{"match_by_name", "queen", []string{"Queen"}},
		{"match_by_name_case_insensitive", "BILLIE", []string{"Billie Eilish"}},
		{"match_by_member", "freddie", []string{"Queen"}},
		{"match_by_creation_date", "1970", []string{"Queen"}},
		{"match_by_first_album", "26-03-2017", []string{"Billie Eilish"}},
		{"no_match_returns_empty", "zzznomatch", []string{}},
		{"match_multiple", "2", []string{"Queen", "Billie Eilish"}},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			results := s.SearchArtists(tc.query)
			if len(results) != len(tc.wantNames) {
				t.Fatalf("got %d results, want %d", len(results), len(tc.wantNames))
			}
			for i, name := range tc.wantNames {
				if results[i].Name != name {
					t.Errorf("result[%d] = %q, want %q", i, results[i].Name, name)
				}
			}
		})
	}
}
