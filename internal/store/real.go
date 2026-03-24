package store

import (
	"groupie-tracker/internal/models"
	"strconv"
	"strings"
)

type RealStore struct {
	Artists   []models.Artist
	Locations models.LocationsResponse
	Data      models.DatesResponse
	Relations models.RelationResponse
}

func (r *RealStore) AllArtists() []models.Artist {
	return r.Artists
}
func (r *RealStore) ArtistByID(id int) (models.Artist, bool) {
	for _, a := range r.AllArtists() {
		if a.ID == id {
			return a, true
		}
	}
	return models.Artist{}, false
}
func (r *RealStore) SearchArtists(query string) []models.Artist {
	var result []models.Artist
	for _, a := range r.AllArtists() {
		if matchesQuery(a, query) {
			result = append(result, a)
		}
	}
	return result
}
func matchesQuery(a models.Artist, query string) bool {
	q := strings.ToLower(query)
	if strings.Contains(strings.ToLower(a.Name), q) {
		return true
	}
	for _, member := range a.Members {
		if strings.Contains(strings.ToLower(member), q) {
			return true
		}
	}
	if strings.Contains(strings.ToLower(strconv.Itoa(a.CreationDate)), q) {
		return true
	}
	if strings.Contains(strings.ToLower(a.FirstAlbum), q) {
		return true
	}

	return false
}

func (r *RealStore) ArtistPageDataByID(id int) (models.ArtistPageData, bool) {
	for _, a := range r.AllArtists() {
		if a.ID == id {
			var locations []string
			for _, l := range r.Locations.Index {
				if l.ID == id {
					locations = l.Locations
				}
			}
			var dates []string
			for _, d := range r.Data.Index {
				if d.ID == id {
					dates = d.Dates
				}
			}
			return models.ArtistPageData{
				Artist:    a,
				Locations: locations,
				Dates:     dates,
			}, true
		}
	}
	return models.ArtistPageData{}, false

}
