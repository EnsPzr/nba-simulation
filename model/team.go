package model

// Team This structure is store teams.
// Each teams has a list of players.
// Each teams is registered to a group.
type Team struct {
	BaseModel
	Name string `json:"name"`

	GroupID int   `json:"group_id"`
	Group   Group `json:"group"`

	Players []Player `json:"players"`
}

// GetTeamsIds This function returns teams ids.
func GetTeamsIds(teams []Team) []int {
	ids := make([]int, 0)
	for _, team := range teams {
		ids = append(ids, team.ID)
	}
	return ids
}

// TeamsToMap This function returb teams map.
// Key is team id and value is team.
func TeamsToMap(teams []Team) map[int]Team {
	teamsMap := make(map[int]Team)
	for _, team := range teams {
		teamsMap[team.ID] = team
	}
	return teamsMap
}
