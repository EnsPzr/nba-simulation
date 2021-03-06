package ws

type WsEvent int

// This events are used to client/server communication.
const (
	WsEventStart WsEvent = iota
	WsEventDataReading
	WsEventDataReady
	WsEventContinue
	WsEventDataSent
	WsEventEnd
	WsEventRestart
)

// Message This structure is used to pass data between client and server.
type Message struct {
	WsEvent     WsEvent    `json:"ws_event"`
	Step        int        `json:"step"`
	RealTime    int        `json:"real_time"`
	VirtualTime int        `json:"virtual_time"`
	ResultVM    []ResultVM `json:"result"`
}

// ResultVM This structure is stored period data.
// Every 5 seconds, data is updated.
type ResultVM struct {
	GameID               int                    `json:"game_id"`
	HomeTeamName         string                 `json:"home_team_name"`
	AwayTeamName         string                 `json:"away_team_name"`
	HomeTeamScore        int                    `json:"home_team_score"`
	AwayTeamScore        int                    `json:"away_team_score"`
	HomeTeamAttackCount  int                    `json:"home_team_attack_count"`
	AwayTeamAttackCount  int                    `json:"away_team_attack_count"`
	HomeTeamPlayerEvents map[int]*PlayerEventVM `json:"home_team_player_events"`
	AwayTeamPlayerEvents map[int]*PlayerEventVM `json:"away_team_player_events"`
}

// PlayerEventVM This structure is store player events.
type PlayerEventVM struct {
	PlayerName                     string `json:"player_name"`
	AssistCount                    int    `json:"assist_count"`
	SuccessfulThreePointShootCount int    `json:"successful_three_point_shoot_count"`
	SuccessfulTwoPointShootCount   int    `json:"successful_two_point_shoot_count"`
}
