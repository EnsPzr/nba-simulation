package model

type EventType int

// These events are randomly selected.
// The simulation is shaped by randomly selected events.
const (
	EventTypeSuccessfulThreePointShoot EventType = iota
	EventTypeFailedThreePointShoot
	EventTypeSuccessfulTwoPointShoot
	EventTypeFailedTwoPointShoot
	EventTypeSuccessfulOnePointShoot
	EventTypeFailedOnePointShoot
	EventTypeCatchTheBall
	EventTypeToPass
	EventTypeDribblingTheBall
)

// DidTheBallChangeSides If the throw is successful or the other team player catches the ball, the ball
// is in the other team. In this case this function return true.
func (event EventType) DidTheBallChangeSides() bool {
	return event == EventTypeSuccessfulThreePointShoot ||
		event == EventTypeSuccessfulTwoPointShoot ||
		event == EventTypeSuccessfulOnePointShoot ||
		event == EventTypeCatchTheBall
}

// FailedShoot If shooting is failed, this function return true.
func (event EventType) FailedShoot() bool {
	return event == EventTypeFailedThreePointShoot ||
		event == EventTypeFailedTwoPointShoot ||
		event == EventTypeFailedOnePointShoot
}

// SuccessfulShoot If shooting is successful, this function return true.
func (event EventType) SuccessfulShoot() bool {
	return event == EventTypeSuccessfulThreePointShoot ||
		event == EventTypeSuccessfulTwoPointShoot ||
		event == EventTypeSuccessfulOnePointShoot
}

// Event This structure stores all movement by games.
type Event struct {
	BaseModel
	GameID     int       `json:"game_id"`
	Game       Game      `json:"game"`
	PlayerID   int       `json:"player_id"`
	Player     Player    `json:"player"`
	TeamID     int       `json:"team_id"`
	Team       Team      `json:"team"`
	Time       int       `json:"time"`
	Type       EventType `json:"type"`
	IsAsist    bool      `json:"is_asist"`
	AttackTime int       `json:"attack_time"`
	IsAttack   bool      `json:"is_attack"`
}
