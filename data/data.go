package data

// This function seeding database.
// 1- Create groups.
// 2- Create teams.
// 3- Create players.
// 4- Create games(matches).
// 5- Create players in the matches.
// 6- Create matches events(scores, catch the ball, pass, dribling to ball).
func InitData() {
	groups := CreateGroups()
	teams := CreateTeams(groups)
	_ = CreatePlayers(teams)
	games := CreateGames(teams)
	CreateGamesPlayers(games)
	PlayGame(games)
}
