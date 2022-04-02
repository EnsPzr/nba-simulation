package data

import (
	"github.com/enspzr/nba-simulation/model"
	"github.com/enspzr/nba-simulation/storage/postgre"
	"strconv"
)

func CreateGroups() []model.Group {
	groups := make([]model.Group, 0)
	err := postgre.DB().Find(&groups).Error
	if err != nil {
		panic("Groups find error: " + err.Error())
	}
	if len(groups) > 0 {
		return groups
	}

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
