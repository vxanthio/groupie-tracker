// Package api handles all communication with the Groupie Trackers external REST API.
// It fetches four datasets — artists, locations, dates, and relations — over HTTP,
// decodes the JSON responses into typed models, and holds the result in a
// package-level store that the rest of the application reads at request time.
// A shared http.Client with a 10-second timeout is used for all outbound requests
// to prevent the server from hanging if the upstream API is slow or unresponsive.
package api

import (
	"encoding/json"
	"fmt"
	"groupie-tracker/internal/models"
	"net/http"
	"time"
)

// AppData holds the complete dataset fetched from the external API at startup.
// It is populated once by LoadData and then read-only for the lifetime of the server.
type AppData struct {
	Artists   []models.Artist
	Locations models.LocationsResponse
	Dates     models.DatesResponse
	Relations models.RelationResponse
}

var data AppData

// LoadData is the public entry point for populating the in-memory data store.
// It calls loadDataFromURLs with the four production API endpoints and should
// be invoked exactly once during application startup, before the HTTP server
// begins accepting requests. If any endpoint fails to respond or returns
// malformed JSON, an error is returned and the application should not start.
func LoadData() error {
	return loadDataFromURLs(
		"https://groupietrackers.herokuapp.com/api/artists",
		"https://groupietrackers.herokuapp.com/api/locations",
		"https://groupietrackers.herokuapp.com/api/dates",
		"https://groupietrackers.herokuapp.com/api/relation",
	)
}
// GetData returns a copy of the AppData struct populated by LoadData.
// It is used by main to pass the loaded datasets into the RealStore,
// which the HTTP handlers then query for every incoming request.
func GetData() AppData {
	return data
}

var httpClient = &http.Client{Timeout: 10 * time.Second}

// loadDataFromURLs sequentially fetches and JSON-decodes each of the four API
// endpoints into the package-level data variable. Each request is made with
// the shared httpClient which enforces a 10-second timeout per call.
// The function is intentionally separate from LoadData so that tests can
// inject a local httptest.Server URL instead of hitting the real API.
// Errors from any fetch or decode step are wrapped with context and returned
// immediately — no partial data is used if any step fails.
func loadDataFromURLs(artistsURL, locationsURL, datesURL, relationsURL string) error {
	respArt, err := httpClient.Get(artistsURL)
	if err != nil {
		return fmt.Errorf("artist fetch failed: %w", err)
	}
	defer respArt.Body.Close()
	if err = json.NewDecoder(respArt.Body).Decode(&data.Artists); err != nil {
		return fmt.Errorf("decoding failed: %w", err)
	}
	respLoc, err := httpClient.Get(locationsURL)
	if err != nil {
		return fmt.Errorf("location fetch failed: %w", err)
	}
	defer respLoc.Body.Close()
	if err = json.NewDecoder(respLoc.Body).Decode(&data.Locations); err != nil {
		return fmt.Errorf("location decode failed: %w", err)
	}
	respDate, err := httpClient.Get(datesURL)
	if err != nil {
		return fmt.Errorf("date fetch failed: %w", err)
	}
	defer respDate.Body.Close()
	if err = json.NewDecoder(respDate.Body).Decode(&data.Dates); err != nil {
		return fmt.Errorf("date decode failed: %w", err)
	}
	respRel, err := httpClient.Get(relationsURL)
	if err != nil {
		return fmt.Errorf("relations fetch failed: %w", err)
	}
	defer respRel.Body.Close()
	if err = json.NewDecoder(respRel.Body).Decode(&data.Relations); err != nil {
		return fmt.Errorf("relations decode failed: %w", err)
	}
	return nil
}
