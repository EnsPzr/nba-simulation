package model

// Player This structure is store all players.
// Each player is registered to a team.
type Player struct {
	BaseModel
	Name   string `json:"name"`
	TeamID int    `json:"team_id"`
	Team   Team   `json:"team"`
}

// GetPlayerIds This function return players ids.
func GetPlayerIds(players []Player) []int {
	ids := make([]int, 0)
	for _, player := range players {
		ids = append(ids, player.ID)
	}
	return ids
}

// PlayersToMap This function returns players map.
// Key is player id and value is player.
func PlayersToMap(players []Player) map[int]Player {
	playersMap := make(map[int]Player)
	for _, player := range players {
		playersMap[player.ID] = player
	}
	return playersMap
}
