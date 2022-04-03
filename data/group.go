package data

import (
	"github.com/enspzr/nba-simulation/model"
	"github.com/enspzr/nba-simulation/storage/postgre"
	"strconv"
)

// CreateGroups This function creates groups.
// Just create 6 groups.
// Each group is at a conference.
// The NBA has two conferences.
func CreateGroups() []model.Group {
	groups := make([]model.Group, 0)
	// Check if the table is empty.
	err := postgre.DB().Find(&groups).Error
	if err != nil {
		panic("Groups find error: " + err.Error())
	}
	// If the table is not empty, return.
	if len(groups) > 0 {
		return groups
	}

	// If the table is empty, create groups.
	for i := 0; i < 6; i++ {
		conference := model.ConferenceEastern
		if i%2 == 0 {
			conference = model.ConferenceWestern
		}
		groups = append(groups, model.Group{
			Name:       "Group " + strconv.Itoa(i+1),
			Conference: conference,
		})
	}
	err = postgre.DB().Create(&groups).Error
	if err != nil {
		panic("Groups create error: " + err.Error())
	}
	return groups
}
