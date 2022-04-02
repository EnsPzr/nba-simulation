package data

import (
	"fmt"
	"github.com/enspzr/nba-simulation/model"
	"github.com/enspzr/nba-simulation/storage/postgre"
	"sort"
	"sync"
)

func InitData() {
	groups := CreateGroups()
	teams := CreateTeams(groups)
	players := CreatePlayers(teams)
	games := CreateGames(teams)
	CreateGamesPlayers(games)
	PlayGame(games)

	teamsMap := model.TeamsToMap(teams)
	playersMap := model.PlayersToMap(players)
	eventsMap := make(map[string][]model.Event)
	var events []model.Event
	err := postgre.DB().Find(&events).Error
	if err != nil {
		panic(err)
	}
	for _, event := range events {
		eventsMap[fmt.Sprintf("%d-%d", event.GameID, event.TeamID)] = append(eventsMap[fmt.Sprintf("%d-%d", event.GameID, event.TeamID)], event)
	}

	periodMap := make(map[int][]ResultVM)
	for i := 1; i <= 48; i++ {
		vms := createResultForPeriod(i, teamsMap, playersMap, games, eventsMap)
		sort.Slice(vms, func(i, j int) bool {
			return vms[i].GameID < vms[j].GameID
		})
		periodMap[i] = vms
	}
	fmt.Println("Elhamdülillah")
}

type ResultVM struct {
	GameID               int                    `json:"game_id"`
	HomeTeamName         string                 `json:"home_team_name"`
	AwayTeamName         string                 `json:"away_team_name"`
	HomeTeamScore        int                    `json:"home_team_score"`
	AwayTeamScore        int                    `json:"away_team_score"`
	HomeTeamAttackCount  int                    `json:"home_team_attack_count"`
	AwayTeamAttackCount  int                    `json:"away_team_attack_count"`
	HomeTeamPlayerEvents map[int]*PlayerEventVM `json:"home_team_player_events"`
	AwayTeamPlayerEvents map[int]*PlayerEventVM `json:"away_team_player_events"`
}
type PlayerEventVM struct {
	PlayerName                     string `json:"player_name"`
	AssistCount                    int    `json:"assist_count"`
	SuccessfulThreePointShootCount int    `json:"successful_three_point_shoot_count"`
	SuccessfulTwoPointShootCount   int    `json:"successful_two_point_shoot_count"`
}

func createResultForPeriod(period int, teams map[int]model.Team, players map[int]model.Player, games []model.Game, events map[string][]model.Event) []ResultVM {
	resultVM := make([]ResultVM, 0)
	var wg sync.WaitGroup
	wg.Add(len(games))
	for _, game := range games {
		homeTeam := teams[game.HomeTeamID]
		awayTeam := teams[game.AwayTeamID]
		homeTeamEvents := events[fmt.Sprintf("%d-%d", game.ID, game.HomeTeamID)]
		awayTeamEvents := events[fmt.Sprintf("%d-%d", game.ID, game.AwayTeamID)]
		go createResultForGame(period, game.ID, homeTeam, awayTeam, homeTeamEvents, awayTeamEvents, players, &wg, &resultVM)
	}
	wg.Wait()
	return resultVM
}

func createResultForGame(period, gameId int,
	homeTeam model.Team, awayTeam model.Team,
	homeTeamEvents []model.Event, awayTeamEvents []model.Event,
	players map[int]model.Player,
	s *sync.WaitGroup, vm *[]ResultVM) {
	defer s.Done()
	resultVM := ResultVM{
		GameID:               gameId,
		HomeTeamName:         homeTeam.Name,
		AwayTeamName:         awayTeam.Name,
		HomeTeamPlayerEvents: make(map[int]*PlayerEventVM),
		AwayTeamPlayerEvents: make(map[int]*PlayerEventVM),
	}
	sort.Slice(homeTeamEvents, func(i, j int) bool {
		return homeTeamEvents[i].Time < homeTeamEvents[j].Time
	})
	sort.Slice(awayTeamEvents, func(i, j int) bool {
		return awayTeamEvents[i].Time < awayTeamEvents[j].Time
	})
	for _, event := range homeTeamEvents {
		if event.Time > period*60 {
			break
		}
		if event.IsAttack {
			resultVM.HomeTeamAttackCount++
		}
		switch event.Type {
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
		if event.Time > period*60 {
			break
		}
		if event.IsAttack {
			resultVM.AwayTeamAttackCount++
		}
		switch event.Type {
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
		case model.EventTypeSuccessfulTwoPointShoot:
			resultVM.HomeTeamScore += 2
			if val, ok := resultVM.AwayTeamPlayerEvents[event.PlayerID]; ok {
				val.SuccessfulTwoPointShootCount++
			} else {
				resultVM.AwayTeamPlayerEvents[event.PlayerID] = &PlayerEventVM{
					PlayerName:                   players[event.PlayerID].Name,
					SuccessfulTwoPointShootCount: 1,
				}
			}
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
	*vm = append(*vm, resultVM)
}
