package data

import (
	"github.com/enspzr/nba-simulation/model"
	"github.com/enspzr/nba-simulation/storage/postgre"
	"github.com/enspzr/nba-simulation/utils"
)

// CreateGamesPlayers This function generate players playing in the game.
// Generate 5 main players and substitute players.
func CreateGamesPlayers(games []model.Game) {
	var gamePlayersCount int64

	// Check if the table is empty.
	err := postgre.DB().Model(&model.GamePlayer{}).Count(&gamePlayersCount).Error
	if err != nil {
		panic("GamePlayers find error: " + err.Error())
	}

	// If the table is not empty, return.
	if gamePlayersCount > 0 {
		return
	}

	gamesPlayers := make([]model.GamePlayer, 0)
	for _, game := range games {
		var homeTeamPlayers []model.Player

		// Get home team players.
		err = postgre.DB().Model(&model.Player{}).Where("team_id = ?", game.HomeTeamID).Find(&homeTeamPlayers).Error
		if err != nil {
			panic("Team Players find error: " + err.Error())
		}

		// Create home teams players ids array.
		homeTeamPlayersIDs := model.GetPlayerIds(homeTeamPlayers)

		// Select 12 players.
		// Select 5 main players.
		// Select 7 substitute players.
		for i := 0; i < 12; i++ {
			playerType := model.GamePlayerTypeMain
			if i > 4 {
				playerType = model.GamePlayerTypeSubstitute
			}
			// Select randomly index from home team players ids array.
			next := random.Intn(len(homeTeamPlayersIDs))
			playerId := homeTeamPlayersIDs[next]

			// Create game player.
			gamesPlayers = append(gamesPlayers, model.GamePlayer{
				GameId:     game.ID,
				PlayerId:   playerId,
				PlayerType: playerType,
			})

			// Remove player from home team players ids array.
			homeTeamPlayersIDs = utils.Remove(homeTeamPlayersIDs, playerId)
		}

		// Get away team players.
		var awayTeamPlayers []model.Player
		err = postgre.DB().Model(&model.Player{}).Where("team_id = ?", game.AwayTeamID).Find(&awayTeamPlayers).Error
		if err != nil {
			panic("Team Players find error: " + err.Error())
		}

		// Create away teams players ids array.
		awayTeamPlayerIds := model.GetPlayerIds(awayTeamPlayers)

		// Select 12 players.
		// Select 5 main players.
		// Select 7 substitute players.
		for i := 0; i < 12; i++ {
			playerType := model.GamePlayerTypeMain
			if i > 4 {
				playerType = model.GamePlayerTypeSubstitute
			}
			// Select randomly index from away team players ids array.
			next := random.Intn(len(homeTeamPlayersIDs))
			playerId := awayTeamPlayerIds[next]
			// Create game player.
			gamesPlayers = append(gamesPlayers, model.GamePlayer{
				GameId:     game.ID,
				PlayerId:   playerId,
				PlayerType: playerType,
			})

			// Remove player from away team players ids array.
			awayTeamPlayerIds = utils.Remove(awayTeamPlayerIds, playerId)
		}
	}

	// Insert game players.
	err = postgre.DB().Create(&gamesPlayers).Error
	if err != nil {
		panic("GamePlayers create error: " + err.Error())
	}
}
