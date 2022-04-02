package ws

type WsEvent int

const (
	WsEventStart WsEvent = iota
	WsEventDataReady
	WsEventContinue
	WsEventDataSent
)

type Message struct {
	Event    WsEvent  `json:"event"`
	Step     int      `json:"step"`
	ResultVM ResultVM `json:"result"`
}

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
type PlayerEventVM struct {
	PlayerName                     string `json:"player_name"`
	AssistCount                    int    `json:"assist_count"`
	SuccessfulThreePointShootCount int    `json:"successful_three_point_shoot_count"`
	SuccessfulTwoPointShootCount   int    `json:"successful_two_point_shoot_count"`
}
