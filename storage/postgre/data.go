package postgre

import (
	"github.com/enspzr/nba-simulation/model"
	"github.com/enspzr/nba-simulation/utils"
	"math/rand"
	"strconv"
)

func InitData() {
	groups := createGroups()
	teams := createTeams(groups)
	_ = createPlayers(teams)
	createGames(teams)
}

func createGroups() []model.Group {
	groups := make([]model.Group, 0)
	err := db.Find(&groups).Error
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
	err = db.Create(&groups).Error
	if err != nil {
		panic("Groups create error: " + err.Error())
	}
	return groups
}

func createTeams(groups []model.Group) []model.Team {
	teams := make([]model.Team, 0)
	err := db.Find(&teams).Error
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

	err = db.Create(&teams).Error
	if err != nil {
		panic("Teams create error: " + err.Error())
	}
	return teams
}

func createPlayers(teams []model.Team) []model.Player {
	players := make([]model.Player, 0)
	err := db.Find(&players).Error
	if err != nil {
		panic("Players find error: " + err.Error())
	}
	if len(players) > 0 {
		return players
	}
	for i := 0; i < 15*30; i++ {
		players = append(players, model.Player{
			Name:   "Player " + strconv.Itoa(i+1),
			TeamID: teams[i%15].ID,
		})
	}
	err = db.Create(&players).Error
	if err != nil {
		panic("Players create error: " + err.Error())
	}
	return players
}

func createGames(teams []model.Team) {
	var gameCount int64
	err := db.Model(&model.Game{}).Count(&gameCount).Error
	if err != nil {
		panic("Games find error: " + err.Error())
	}
	if gameCount > 0 {
		return
	}

	games := make([]model.Game, 0)
	avaibleTeamsIDs := getTeamsIds(teams)
	for i := 0; i < 15; i++ {
		homeTeamID := avaibleTeamsIDs[rand.Intn(len(avaibleTeamsIDs))]
		avaibleTeamsIDs = utils.Remove(avaibleTeamsIDs, homeTeamID)
		awayTeamID := avaibleTeamsIDs[rand.Intn(len(avaibleTeamsIDs))]
		avaibleTeamsIDs = utils.Remove(avaibleTeamsIDs, awayTeamID)
		games = append(games, model.Game{
			HomeTeamID: homeTeamID,
			AwayTeamID: awayTeamID,
		})
	}
	err = db.Create(&games).Error
	if err != nil {
		panic("Games create error: " + err.Error())
	}
}

func getTeamsIds(teams []model.Team) []int {
	ids := make([]int, 0)
	for _, team := range teams {
		ids = append(ids, team.ID)
	}
	return ids
}
