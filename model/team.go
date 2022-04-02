package model

type Team struct {
	BaseModel
	Name string `json:"name"`

	GroupID int   `json:"group_id"`
	Group   Group `json:"group"`

	Players []Player `json:"players"`
}

func GetTeamsIds(teams []Team) []int {
	ids := make([]int, 0)
	for _, team := range teams {
		ids = append(ids, team.ID)
	}
	return ids
}

func TeamsToMap(teams []Team) map[int]Team {
	teamsMap := make(map[int]Team)
	for _, team := range teams {
		teamsMap[team.ID] = team
	}
	return teamsMap
}
