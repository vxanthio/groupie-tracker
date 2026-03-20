package store

import "groupie-tracker/internal/models"

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
	return r.AllArtists()
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
