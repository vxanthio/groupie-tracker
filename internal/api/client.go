package api

import (
	"encoding/json"
	"fmt"
	"groupie-tracker/internal/models"
	"net/http"
)

type appData struct {
	Artists   []models.Artist
	Locations models.LocationsResponse
	Date      models.DatesResponse
	Relations models.RelationResponse
}

var data appData

func loadData() error {
	return loadDataFromURLs(
		"https://groupietrackers.herokuapp.com/api/artists",
		"https://groupietrackers.herokuapp.com/api/locations",
		"https://groupietrackers.herokuapp.com/api/dates",
		"https://groupietrackers.herokuapp.com/api/relation",
	)
}

func loadDataFromURLs(artistsURL, locationsURL, datesURL, relationsURL string) error {
	respArt, err := http.Get(artistsURL)
	if err != nil {
		return fmt.Errorf("artist fetch failed: %w", err)
	}
	defer respArt.Body.Close()
	if err = json.NewDecoder(respArt.Body).Decode(&data.Artists); err != nil {
		return fmt.Errorf("decoding failed: %w", err)
	}
	respLoc, err := http.Get(locationsURL)
	if err != nil {
		return fmt.Errorf("location fetch failed: %w", err)
	}
	defer respLoc.Body.Close()
	if err = json.NewDecoder(respLoc.Body).Decode(&data.Locations); err != nil {
		return fmt.Errorf("location decode failed: %w", err)
	}
	respDate, err := http.Get(datesURL)
	if err != nil {
		return fmt.Errorf("date fetch failed: %w", err)
	}
	defer respDate.Body.Close()
	if err = json.NewDecoder(respDate.Body).Decode(&data.Date); err != nil {
		return fmt.Errorf("date decode failed: %w", err)
	}
	respRel, err := http.Get(relationsURL)
	if err != nil {
		return fmt.Errorf("relations fetch failed: %w", err)
	}
	defer respRel.Body.Close()
	if err = json.NewDecoder(respRel.Body).Decode(&data.Relations); err != nil {
		return fmt.Errorf("relations decode failed: %w", err)
	}
	return nil
}
