package ws

import (
	"fmt"
	"github.com/enspzr/nba-simulation/model"
	"github.com/enspzr/nba-simulation/storage/postgre"
	"sort"
	"sync"
)

// This function initialize period data.
// Function calculate data for every 5 second.
// And write data in periodMap.
// Period map key is period number.
// Period map value is ResultVM array.
func initData() {
	var teams []model.Team
	var players []model.Player
	var games []model.Game

	// Find all teams.
	err := postgre.DB().Find(&teams).Error
	if err != nil {
		panic("Error loading teams => " + err.Error())
	}

	// Find all players.
	err = postgre.DB().Find(&players).Error
	if err != nil {
		panic("Error loading players => " + err.Error())
	}

	// Find all games.
	err = postgre.DB().Find(&games).Error
	if err != nil {
		panic("Error loading games => " + err.Error())
	}

	// Find all events.
	var events []model.Event
	err = postgre.DB().Find(&events).Error
	if err != nil {
		panic("Error loading events => " + err.Error())
	}

	// Teams and players convert to be map. Because for easy access to data.
	teamsMap := model.TeamsToMap(teams)
	playersMap := model.PlayersToMap(players)

	// Events convert to be map. Because for easy access to data.
	eventsMap := make(map[string][]model.Event)
	for _, event := range events {
		eventsMap[fmt.Sprintf("%d-%d", event.GameID, event.TeamID)] = append(eventsMap[fmt.Sprintf("%d-%d", event.GameID, event.TeamID)], event)
	}

	// Create period map for every period. (5 seconds)
	periodMap = make(map[int][]ResultVM)
	for i := 0; i < 48; i++ {
		vms := createResultForPeriod(i+1, teamsMap, playersMap, games, eventsMap)
		periodMap[i] = vms
	}
}

// This function create ResultVM array for every period.
func createResultForPeriod(period int, teams map[int]model.Team, players map[int]model.Player, games []model.Game, events map[string][]model.Event) []ResultVM {
	resultVM := make([]ResultVM, 0)
	// Game count is big. So we need to use goroutine to fast ready data.
	var wg sync.WaitGroup
	// Use lock. Because append is not thread safe.
	var gameLock sync.Mutex
	wg.Add(len(games))
	for _, game := range games {
		// Get home team.
		homeTeam := teams[game.HomeTeamID]
		// Get away team.
		awayTeam := teams[game.AwayTeamID]
		// Get home team events.
		homeTeamEvents := events[fmt.Sprintf("%d-%d", game.ID, game.HomeTeamID)]
		// Get away team events.
		awayTeamEvents := events[fmt.Sprintf("%d-%d", game.ID, game.AwayTeamID)]
		// Create ResultVM for period.
		go createResultForGame(period, game.ID, homeTeam, awayTeam, homeTeamEvents, awayTeamEvents, players, &wg, &gameLock, &resultVM)
	}
	// Wait for all goroutine finish.
	wg.Wait()
	return resultVM
}

func createResultForGame(period, gameId int,
	homeTeam model.Team, awayTeam model.Team,
	homeTeamEvents []model.Event, awayTeamEvents []model.Event,
	players map[int]model.Player,
	s *sync.WaitGroup, gameLock *sync.Mutex, vm *[]ResultVM) {
	defer s.Done()
	resultVM := ResultVM{
		GameID:               gameId,
		HomeTeamName:         homeTeam.Name,
		AwayTeamName:         awayTeam.Name,
		HomeTeamPlayerEvents: make(map[int]*PlayerEventVM),
		AwayTeamPlayerEvents: make(map[int]*PlayerEventVM),
	}
	// Sort events by time. Because we need to get events in order.
	sort.Slice(homeTeamEvents, func(i, j int) bool {
		return homeTeamEvents[i].Time < homeTeamEvents[j].Time
	})
	// Sort events by time. Because we need to get events in order.
	sort.Slice(awayTeamEvents, func(i, j int) bool {
		return awayTeamEvents[i].Time < awayTeamEvents[j].Time
	})

	for _, event := range homeTeamEvents {
		// If event time bigger than period time break loop.
		if event.Time > period*60 {
			break
		}
		// If event is attack, attack count + 1.
		if event.IsAttack {
			resultVM.HomeTeamAttackCount++
		}
		switch event.Type {
		// If event type successful three point shoot, team score + 3. And player statics updated.
		case model.EventTypeSuccessfulThreePointShoot:
			resultVM.HomeTeamScore += 3
			if val, ok := resultVM.HomeTeamPlayerEvents[event.PlayerID]; ok {
				val.SuccessfulThreePointShootCount++
			} else {
				resultVM.HomeTeamPlayerEvents[event.PlayerID] = &PlayerEventVM{
					PlayerName:                     players[event.PlayerID].Name,
					SuccessfulThreePointShootCount: 1,
				}
			}
			// If event type successful two point shoot, team score + 2. And player statics updated.
		case model.EventTypeSuccessfulTwoPointShoot:
			resultVM.HomeTeamScore += 2
			if val, ok := resultVM.HomeTeamPlayerEvents[event.PlayerID]; ok {
				val.SuccessfulTwoPointShootCount++
			} else {
				resultVM.HomeTeamPlayerEvents[event.PlayerID] = &PlayerEventVM{
					PlayerName:                   players[event.PlayerID].Name,
					SuccessfulTwoPointShootCount: 1,
				}
			}
			// If event type pass and event is a assist, player statics updated. Assist count + 1.
		case model.EventTypeToPass:
			if event.IsAsist {
				if val, ok := resultVM.HomeTeamPlayerEvents[event.PlayerID]; ok {
					val.AssistCount++
				} else {
					resultVM.HomeTeamPlayerEvents[event.PlayerID] = &PlayerEventVM{
						PlayerName:  players[event.PlayerID].Name,
						AssistCount: 1,
					}
				}
			}
		}
	}

	for _, event := range awayTeamEvents {
		// If event time bigger than period time break loop.
		if event.Time > period*60 {
			break
		}
		// If event is attack, attack count + 1.
		if event.IsAttack {
			resultVM.AwayTeamAttackCount++
		}
		switch event.Type {
		// If event type successful three point shoot, team score + 3. And player statics updated.
		case model.EventTypeSuccessfulThreePointShoot:
			resultVM.AwayTeamScore += 3
			if val, ok := resultVM.AwayTeamPlayerEvents[event.PlayerID]; ok {
				val.SuccessfulThreePointShootCount++
			} else {
				resultVM.AwayTeamPlayerEvents[event.PlayerID] = &PlayerEventVM{
					PlayerName:                     players[event.PlayerID].Name,
					SuccessfulThreePointShootCount: 1,
				}
			}
			// If event type successful two point shoot, team score + 2. And player statics updated.
		case model.EventTypeSuccessfulTwoPointShoot:
			resultVM.AwayTeamScore += 2
			if val, ok := resultVM.AwayTeamPlayerEvents[event.PlayerID]; ok {
				val.SuccessfulTwoPointShootCount++
			} else {
				resultVM.AwayTeamPlayerEvents[event.PlayerID] = &PlayerEventVM{
					PlayerName:                   players[event.PlayerID].Name,
					SuccessfulTwoPointShootCount: 1,
				}
			}
			// If event type pass and event is a assist, player statics updated. Assist count + 1.
		case model.EventTypeToPass:
			if event.IsAsist {
				if val, ok := resultVM.AwayTeamPlayerEvents[event.PlayerID]; ok {
					val.AssistCount++
				} else {
					resultVM.AwayTeamPlayerEvents[event.PlayerID] = &PlayerEventVM{
						PlayerName:  players[event.PlayerID].Name,
						AssistCount: 1,
					}
				}
			}
		}
	}
	// Use lock. Because append is not thread safe.
	gameLock.Lock()
	*vm = append(*vm, resultVM)
	gameLock.Unlock()
}
