package model

type GamePlayerType int

const (
	GamePlayerTypeMain GamePlayerType = iota
	GamePlayerTypeSubstitute
)

// GamePlayer This structure is store players by games.
// There is two player type of each game.
// Main player and substitute player.
type GamePlayer struct {
	BaseModel
	GameId     int            `json:"game_id"`
	Game       Game           `json:"game"`
	PlayerId   int            `json:"player_id"`
	Player     Player         `json:"player"`
	PlayerType GamePlayerType `json:"player_type"`
}
