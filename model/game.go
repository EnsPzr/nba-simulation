package model

// Game This structure is store game data.
type Game struct {
	BaseModel
	HomeTeamID  int  `json:"home_team_id"`
	HomeTeam    Team `json:"home_team"`
	AwayTeamID  int  `json:"away_team_id"`
	AwayTeam    Team `json:"away_team"`
	AttackCount int  `json:"attack_count"`
	TotalScore  int  `json:"total_score"`
}
