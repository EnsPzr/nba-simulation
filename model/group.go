package model

type Conference string

const (
	ConferenceEastern Conference = "Eastern"
	ConferenceWestern Conference = "Western"
)

// Group This structure is store groups.
// Count of groups is 6.
// Each group has 5 teams.
// 3 groups in eastern conference and 3 groups in western conference.
type Group struct {
	BaseModel
	Name       string     `json:"name"`
	Conference Conference `json:"conference"`
	Teams      []Team     `json:"teams"`
}
