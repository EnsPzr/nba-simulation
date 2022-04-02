package data

func InitData() {
	groups := CreateGroups()
	teams := CreateTeams(groups)
	_ = CreatePlayers(teams)
	games := CreateGames(teams)
	CreateGamesPlayers(games)
	PlayGame(games)
}
