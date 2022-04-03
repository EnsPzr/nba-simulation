package data

import (
	"github.com/enspzr/nba-simulation/model"
	"github.com/enspzr/nba-simulation/storage/postgre"
	"strconv"
)

// CreatePlayers This function generate players.
// Function generate 15 players for each team.
func CreatePlayers(teams []model.Team) []model.Player {
	players := make([]model.Player, 0)

	// Check if the table is empty.
	err := postgre.DB().Find(&players).Error
	if err != nil {
		panic("Players find error: " + err.Error())
	}
	// If the table is empty, generate players.
	if len(players) > 0 {
		return players
	}
	for i := 0; i < 15*30; i++ {
		// Generate random player.
		// Assignment next team.
		players = append(players, model.Player{
			Name:   "Player " + strconv.Itoa(i+1),
			TeamID: teams[i%30].ID,
		})
	}
	// Save players to database.
	err = postgre.DB().Create(&players).Error
	if err != nil {
		panic("Players create error: " + err.Error())
	}
	return players
}
