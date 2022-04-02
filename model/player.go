package model

type Player struct {
	BaseModel
	Name   string `json:"name"`
	TeamID int    `json:"team_id"`
	Team   Team   `json:"team"`
}

func GetPlayerIds(players []Player) []int {
	ids := make([]int, 0)
	for _, player := range players {
		ids = append(ids, player.ID)
	}
	return ids
}

func PlayersToMap(players []Player) map[int]Player {
	playersMap := make(map[int]Player)
	for _, player := range players {
		playersMap[player.ID] = player
	}
	return playersMap
}
