package data

import (
	"github.com/enspzr/nba-simulation/model"
	"math/rand"
	"sync"
	"time"
)

var random = rand.New(rand.NewSource(time.Now().UnixNano()))
var lck = &sync.Mutex{}

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

func selectRandomTeam(teams []model.Team) *model.Team {
	lck.Lock()
	defer lck.Unlock()
	return &teams[random.Intn(len(teams))]
}

func getRandomSecond() int {
	lck.Lock()
	defer lck.Unlock()
	return random.Intn(10) + 1
}

func getRandomEventType() model.EventType {
	lck.Lock()
	defer lck.Unlock()
	return model.EventType(random.Intn(8))
}
