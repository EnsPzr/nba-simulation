package data

import (
	"github.com/enspzr/nba-simulation/model"
	"math/rand"
	"sync"
	"time"
)

// This file contains the random data generator.
// It is used to generate random data for the simulation.

var random = rand.New(rand.NewSource(time.Now().UnixNano()))

// This structure was used to avoid conflicts.
var lck = &sync.Mutex{}

// Select random player from the list.
// If selectedId is not nil, it will be used to avoid selecting the same player twice.
// If selectedId is nil, it will be ignored.
func selectRandomPlayerId(selectedId int, playersIDs []int) int {
	lck.Lock()
	defer lck.Unlock()
	if selectedId == 0 {
		return playersIDs[random.Intn(len(playersIDs))]
	}
	newPlayerId := selectedId
	for newPlayerId == selectedId {
		newPlayerId = playersIDs[random.Intn(len(playersIDs))]
	}
	return newPlayerId
}

// This function use select random playing team.
func selectRandomTeam(teams []model.Team) *model.Team {
	lck.Lock()
	defer lck.Unlock()
	return &teams[random.Intn(len(teams))]
}

// This function random generate second.
// This second used to calculate the time of the game.
func getRandomSecond() int {
	lck.Lock()
	defer lck.Unlock()
	return random.Intn(10) + 1
}

// This function random generate event(score, pass, catch the ball etc).
// This event used to generate simulation.
func getRandomEventType() model.EventType {
	lck.Lock()
	defer lck.Unlock()
	return model.EventType(random.Intn(8))
}
