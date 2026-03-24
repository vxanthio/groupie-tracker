package store

import (
	"groupie-tracker/internal/models"
)

type MockStore struct {
}

func (m *MockStore) AllArtists() []models.Artist {
	return []models.Artist{
		{ID: 1, Name: "Billie Eilish"},
		{ID: 2, Name: "System of a down"},
	}
}
func (m *MockStore) ArtistByID(id int) (models.Artist, bool) {
	for _, a := range m.AllArtists() {
		if a.ID == id {
			return a, true
		}
	}
	return models.Artist{}, false
}
func (m *MockStore) SearchArtists(query string) []models.Artist {
	var results []models.Artist
	for _, a := range m.AllArtists() {
		if matchesQuery(a, query) {
			results = append(results, a)
		}
	}
	return results
}
func (m *MockStore) ArtistPageDataByID(id int) (models.ArtistPageData, bool) {
	for _, a := range m.AllArtists() {
		if a.ID == id {
			return models.ArtistPageData{
				Artist:    a,
				Locations: []string{},
				Dates:     []string{},
			}, true
		}
	}
	return models.ArtistPageData{}, false
}
