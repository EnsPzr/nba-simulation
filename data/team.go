package data

import (
	"github.com/enspzr/nba-simulation/model"
	"github.com/enspzr/nba-simulation/storage/postgre"
	"strconv"
)

// CreateTeams This function create Teams.
// Each teams have a group.
func CreateTeams(groups []model.Group) []model.Team {
	teams := make([]model.Team, 0)

	// Check if the table is empty.
	err := postgre.DB().Find(&teams).Error
	if err != nil {
		panic("Teams find error: " + err.Error())
	}

	// If the table is not empty, return.
	if len(teams) > 0 {
		return teams
	}

	// Generate 30 teams.
	for i := 0; i < 30; i++ {
		teams = append(teams, model.Team{
			Name: "Team " + strconv.Itoa(i+1),
			// Select a next group.
			GroupID: groups[i/5].ID,
		})
	}

	// Save the teams.
	err = postgre.DB().Create(&teams).Error
	if err != nil {
		panic("Teams create error: " + err.Error())
	}
	return teams
}
