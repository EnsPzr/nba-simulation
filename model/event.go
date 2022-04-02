package model

type EventType int

const (
	EventTypeSuccessfulThreePointShoot EventType = iota
	EventTypeFailedThreePointShoot
	EventTypeSuccessfulTwoPointShoot
	EventTypeFailedTwoPointShoot
	EventTypeSuccessfulOnePointShoot
	EventTypeFailedOnePointShoot
	// topu kaptır
	EventTypeCatchTheBall
	// pas ver
	EventTypeToPass
	// topu sürme
	EventTypeDribblingTheBall
)

// top takım değiştirdi mi
func (event EventType) DidTheBallChangeSides() bool {
	return event == EventTypeSuccessfulThreePointShoot ||
		event == EventTypeSuccessfulTwoPointShoot ||
		event == EventTypeSuccessfulOnePointShoot ||
		event == EventTypeCatchTheBall
}

func (event EventType) FailedShoot() bool {
	return event == EventTypeFailedThreePointShoot ||
		event == EventTypeFailedTwoPointShoot ||
		event == EventTypeFailedOnePointShoot
}

func (event EventType) SuccessfulShoot() bool {
	return event == EventTypeSuccessfulThreePointShoot ||
		event == EventTypeSuccessfulTwoPointShoot ||
		event == EventTypeSuccessfulOnePointShoot
}

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
