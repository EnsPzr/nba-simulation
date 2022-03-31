package model

type Player struct {
	BaseModel
	Name   string `json:"name"`
	TeamID int    `json:"team_id"`
	Team   Team   `json:"team"`
}
