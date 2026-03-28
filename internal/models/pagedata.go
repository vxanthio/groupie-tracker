package models

// ArtistPageData is the view model passed to the artist detail template.
// It combines the core Artist data with the resolved concert locations and dates,
// which are fetched separately from the API and joined together in the store layer.
// This struct is intentionally flat so the template can access all fields directly
// without nested lookups.
type ArtistPageData struct {
	Artist    Artist
	Locations []string
	Dates     []string
}
