package model

type Conference string

const (
	// doğu
	ConferenceEastern Conference = "Eastern"
	// batı
	ConferenceWestern Conference = "Western"
)

type Group struct {
	BaseModel
	Name       string     `json:"name"`
	Conference Conference `json:"conference"`
	Teams      []Team     `json:"teams"`
}
