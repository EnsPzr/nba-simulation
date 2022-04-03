package data

import (
	"github.com/enspzr/nba-simulation/model"
	"github.com/enspzr/nba-simulation/storage/postgre"
	"github.com/enspzr/nba-simulation/utils"
)

// CreateGames This function crate games.
// Just create 15 games.
// Each game has 2 teams.
// Teams are randomly selected from the list of teams.
func CreateGames(teams []model.Team) []model.Game {
	games := make([]model.Game, 0)

	// Check if the table is empty.
	err := postgre.DB().Find(&games).Error
	if err != nil {
		panic("Games find error: " + err.Error())
	}

	// If the table is not empty, return.
	if len(games) > 0 {
		return games
	}

	// Get teams Ids.
	avaibleTeamsIDs := model.GetTeamsIds(teams)
	for i := 0; i < 15; i++ {
		// Select home team.
		homeTeamID := avaibleTeamsIDs[random.Intn(len(avaibleTeamsIDs))]

		// Remove home team id in the avaible teams list.
		avaibleTeamsIDs = utils.Remove(avaibleTeamsIDs, homeTeamID)

		// Select away team.
		awayTeamID := avaibleTeamsIDs[random.Intn(len(avaibleTeamsIDs))]

		// Remove away team id in the avaible teams list.
		avaibleTeamsIDs = utils.Remove(avaibleTeamsIDs, awayTeamID)

		// Create game and add it to the list.
		games = append(games, model.Game{
			HomeTeamID: homeTeamID,
			AwayTeamID: awayTeamID,
		})
	}
	//Create records.
	err = postgre.DB().Create(&games).Error
	if err != nil {
		panic("Games create error: " + err.Error())
	}
	return games
}
