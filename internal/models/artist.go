package models

type Artist struct {
	ID           int      `json:"id"`
	Name         string   `json:"name"`
	Image        string   `json:"image"`
	Members      []string `json:"members"`
	CreationDate int      `json:"creationDate"`
	FirstAlbum   string   `json:"firstAlbum"`
	Locations    string   `json:"locations"`
	ConcertDates string   `json:"concertDates"`
	Relations    string   `json:"relations"`
}
type LocationsResponse struct {
	Index []Locations `json:"index"`
}
type Locations struct {
	ID        int      `json:"id"`
	Locations []string `json:"locations"`
}
type DatesResponse struct {
	Index []Dates `json:"index"`
}
type Dates struct {
	ID    int      `json:"id"`
	Dates []string `json:"dates"`
}
type RelationResponse struct {
	Index []Relation `json:"index"`
}
type Relation struct {
	ID             int                 `json:"id"`
	DatesLocations map[string][]string `json:"datesLocations"`
}
