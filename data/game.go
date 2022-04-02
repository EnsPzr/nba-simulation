package data

import (
	"github.com/enspzr/nba-simulation/model"
	"github.com/enspzr/nba-simulation/storage/postgre"
	"github.com/enspzr/nba-simulation/utils"
)

func CreateGames(teams []model.Team) []model.Game {
	games := make([]model.Game, 0)
	err := postgre.DB().Find(&games).Error
	if err != nil {
		panic("Games find error: " + err.Error())
	}
	if len(games) > 0 {
		return games
	}

	avaibleTeamsIDs := model.GetTeamsIds(teams)
	for i := 0; i < 15; i++ {
		homeTeamID := avaibleTeamsIDs[random.Intn(len(avaibleTeamsIDs))]
		avaibleTeamsIDs = utils.Remove(avaibleTeamsIDs, homeTeamID)
		awayTeamID := avaibleTeamsIDs[random.Intn(len(avaibleTeamsIDs))]
		avaibleTeamsIDs = utils.Remove(avaibleTeamsIDs, awayTeamID)
		games = append(games, model.Game{
			HomeTeamID: homeTeamID,
			AwayTeamID: awayTeamID,
		})
	}
	err = postgre.DB().Create(&games).Error
	if err != nil {
		panic("Games create error: " + err.Error())
	}
	return games
}
