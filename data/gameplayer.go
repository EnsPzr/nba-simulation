package data

import (
	"github.com/enspzr/nba-simulation/model"
	"github.com/enspzr/nba-simulation/storage/postgre"
	"github.com/enspzr/nba-simulation/utils"
)

func CreateGamesPlayers(games []model.Game) {
	var gamePlayersCount int64
	err := postgre.DB().Model(&model.GamePlayer{}).Count(&gamePlayersCount).Error
	if err != nil {
		panic("GamePlayers find error: " + err.Error())
	}
	if gamePlayersCount > 0 {
		return
	}

	gamesPlayers := make([]model.GamePlayer, 0)
	for _, game := range games {
		var homeTeamPlayers []model.Player
		err = postgre.DB().Model(&model.Player{}).Where("team_id = ?", game.HomeTeamID).Find(&homeTeamPlayers).Error
		if err != nil {
			panic("Team Players find error: " + err.Error())
		}

		homeTeamPlayersIDs := model.GetPlayerIds(homeTeamPlayers)

		for i := 0; i < 12; i++ {
			playerType := model.GamePlayerTypeMain
			if i > 4 {
				playerType = model.GamePlayerTypeSubstitute
			}
			next := random.Intn(len(homeTeamPlayersIDs))
			playerId := homeTeamPlayersIDs[next]
			gamesPlayers = append(gamesPlayers, model.GamePlayer{
				GameId:     game.ID,
				PlayerId:   playerId,
				PlayerType: playerType,
			})

			homeTeamPlayersIDs = utils.Remove(homeTeamPlayersIDs, playerId)
		}

		var awayTeamPlayers []model.Player
		err = postgre.DB().Model(&model.Player{}).Where("team_id = ?", game.AwayTeamID).Find(&awayTeamPlayers).Error
		if err != nil {
			panic("Team Players find error: " + err.Error())
		}

		awayTeamPlayerIds := model.GetPlayerIds(awayTeamPlayers)

		for i := 0; i < 12; i++ {
			playerType := model.GamePlayerTypeMain
			if i > 4 {
				playerType = model.GamePlayerTypeSubstitute
			}
			next := random.Intn(len(homeTeamPlayersIDs))
			playerId := awayTeamPlayerIds[next]
			gamesPlayers = append(gamesPlayers, model.GamePlayer{
				GameId:     game.ID,
				PlayerId:   playerId,
				PlayerType: playerType,
			})

			awayTeamPlayerIds = utils.Remove(awayTeamPlayerIds, playerId)
		}
	}

	err = postgre.DB().Create(&gamesPlayers).Error
	if err != nil {
		panic("GamePlayers create error: " + err.Error())
	}
}
