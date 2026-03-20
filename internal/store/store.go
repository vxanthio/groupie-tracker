package store

import "groupie-tracker/internal/models"

type Store interface {
	AllArtists() []models.Artist
	ArtistByID(id int) (models.Artist, bool)
	SearchArtists(query string) []models.Artist
	ArtistPageDataByID(id int) (models.ArtistPageData, bool)
}
