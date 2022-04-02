package data

import (
	"github.com/enspzr/nba-simulation/model"
	"github.com/enspzr/nba-simulation/storage/postgre"
	"strconv"
)

func CreatePlayers(teams []model.Team) []model.Player {
	players := make([]model.Player, 0)
	err := postgre.DB().Find(&players).Error
	if err != nil {
		panic("Players find error: " + err.Error())
	}
	if len(players) > 0 {
		return players
	}
	for i := 0; i < 15*30; i++ {
		players = append(players, model.Player{
			Name:   "Player " + strconv.Itoa(i+1),
			TeamID: teams[i%30].ID,
		})
	}
	err = postgre.DB().Create(&players).Error
	if err != nil {
		panic("Players create error: " + err.Error())
	}
	return players
}
