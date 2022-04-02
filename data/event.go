package data

import (
	"github.com/enspzr/nba-simulation/model"
	"github.com/enspzr/nba-simulation/storage/postgre"
	"sync"
)

func PlayGame(games []model.Game) {
	var eventCouns int64
	err := postgre.DB().Model(&model.Event{}).Count(&eventCouns).Error
	if err != nil {
		panic("Events count error: " + err.Error())
	}
	if eventCouns > 0 {
		return
	}

	var wg sync.WaitGroup
	wg.Add(len(games))
	for _, game := range games {
		go playGameWorker(game, &wg)
	}
	wg.Wait()
}

func playGameWorker(game model.Game, wg *sync.WaitGroup) {
	defer wg.Done()
	homeTeam := getTeam(game.HomeTeamID)
	awayTeam := getTeam(game.AwayTeamID)

	playerIds := make(map[int][]int)
	playerIds[game.HomeTeamID] = getTeamMainPlayersIds(game.HomeTeamID, game.ID)
	playerIds[game.AwayTeamID] = getTeamMainPlayersIds(game.AwayTeamID, game.ID)

	teams := []model.Team{homeTeam, awayTeam}

	totalPassingTime := 0
	teamPassingTime := 0
	teamPlaying := selectRandomTeam(teams)
	playerIDPlaying := selectRandomPlayerId(0, playerIds[teamPlaying.ID])

	events := make([]model.Event, 0)
	for totalPassingTime < 48*60 {
		event := getRandomEventType()
		passingTime := getRandomSecond()
		totalPassingTime += passingTime
		teamPassingTime += passingTime
		if teamPassingTime > 24 {
			// fazlalık zamanı geri verdik
			diff := teamPassingTime - 24
			totalPassingTime -= diff
			teamPassingTime = 0
			teamPlaying = changePlayingTeam(teams, teamPlaying)
			playerIDPlaying = selectRandomPlayerId(0, playerIds[teamPlaying.ID])
			continue
		}

		if totalPassingTime >= 48*60 {
			// oyun bitmiş
			break
		}

		addedToBeEvent := model.Event{
			GameID:     game.ID,
			PlayerID:   playerIDPlaying,
			TeamID:     teamPlaying.ID,
			Time:       totalPassingTime,
			AttackTime: teamPassingTime,
			Type:       event,
		}
		events = append(events, addedToBeEvent)

		if event.DidTheBallChangeSides() {
			teamPassingTime = 0
			teamPlaying = changePlayingTeam(teams, teamPlaying)
			playerIDPlaying = selectRandomPlayerId(0, playerIds[teamPlaying.ID])
			continue
		}

		if event.FailedShoot() {
			teamPassingTime = 0
			teamPlaying = selectRandomTeam(teams)
			playerIDPlaying = selectRandomPlayerId(0, playerIds[teamPlaying.ID])
			continue
		}

		if event == model.EventTypeToPass {
			playerIDPlaying = selectRandomPlayerId(playerIDPlaying, playerIds[teamPlaying.ID])
			continue
		}

		if event == model.EventTypeDribblingTheBall {
			continue
		}
	}

	events[0].IsAttack = true
	for i := 1; i < len(events); i++ {
		if events[i].Type.SuccessfulShoot() &&
			events[i-1].Type == model.EventTypeToPass &&
			events[i].TeamID == events[i-1].TeamID {
			events[i-1].IsAsist = true
		}
		if events[i].TeamID != events[i-1].TeamID || events[i].Type.FailedShoot() || events[i].Type == model.EventTypeCatchTheBall {
			events[i].IsAttack = true
		}
	}

	err := postgre.DB().Create(&events).Error
	if err != nil {
		panic("Events create error: " + err.Error())
	}
}

func getTeam(id int) model.Team {
	var team model.Team
	err := postgre.DB().Model(&model.Team{}).Where("id = ?", id).First(&team).Error
	if err != nil {
		panic("Team find error: " + err.Error())
	}
	return team
}

func getTeamMainPlayersIds(teamId, gameId int) []int {
	var playersIds []int
	err := postgre.DB().Model(&model.GamePlayer{}).Where("game_id = ? AND player_type = ? AND player_id IN "+
		"(SELECT id FROM players WHERE team_id = ?)",
		gameId, model.GamePlayerTypeMain, teamId).
		Pluck("player_id", &playersIds).Error
	if err != nil {
		panic("GamePlayers find error: " + err.Error())
	}
	return playersIds
}

func changePlayingTeam(teams []model.Team, teamPlaying *model.Team) *model.Team {
	if teams[0].ID == teamPlaying.ID {
		return &teams[1]
	}
	return &teams[0]
}
