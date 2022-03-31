package model

type Team struct {
	BaseModel
	Name string `json:"name"`

	GroupID int   `json:"group_id"`
	Group   Group `json:"group"`

	Players []Player `json:"players"`
}
