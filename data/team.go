package data

import (
	"github.com/enspzr/nba-simulation/model"
	"github.com/enspzr/nba-simulation/storage/postgre"
	"strconv"
)

func CreateTeams(groups []model.Group) []model.Team {
	teams := make([]model.Team, 0)
	err := postgre.DB().Find(&teams).Error
	if err != nil {
		panic("Teams find error: " + err.Error())
	}
	if len(teams) > 0 {
		return teams
	}
	for i := 0; i < 30; i++ {
		teams = append(teams, model.Team{
			Name:    "Team " + strconv.Itoa(i+1),
			GroupID: groups[i/5].ID,
		})
	}

	err = postgre.DB().Create(&teams).Error
	if err != nil {
		panic("Teams create error: " + err.Error())
	}
	return teams
}
