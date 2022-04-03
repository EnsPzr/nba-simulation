package data

import (
	"github.com/enspzr/nba-simulation/model"
	"github.com/enspzr/nba-simulation/storage/postgre"
	"sync"
)

// PlayGame Event is a struct that contains all the event for basketball match.
// All pass, score, catch the ball etc contains event struct.
// This function generate game events.
func PlayGame(games []model.Game) {
	var eventCounts int64
	// Check event count. If event count is 0, then we can start to play game.
	err := postgre.DB().Model(&model.Event{}).Count(&eventCounts).Error
	if err != nil {
		panic("Events count error: " + err.Error())
	}

	// If event count is not 0, return.
	if eventCounts > 0 {
		return
	}

	// If event count is 0, play all games.
	// Game count is big. So we need to use goroutine to fast generate game data.
	var wg sync.WaitGroup
	wg.Add(len(games))
	for _, game := range games {
		go playGameWorker(game, &wg)
	}
	wg.Wait()
}

// This function generate game events only single game.
func playGameWorker(game model.Game, wg *sync.WaitGroup) {
	defer wg.Done()
	// Get home and away team.
	homeTeam := getTeam(game.HomeTeamID)
	awayTeam := getTeam(game.AwayTeamID)

	// Get home and away team main players.
	playerIds := make(map[int][]int)
	playerIds[game.HomeTeamID] = getTeamMainPlayersIds(game.HomeTeamID, game.ID)
	playerIds[game.AwayTeamID] = getTeamMainPlayersIds(game.AwayTeamID, game.ID)

	teams := []model.Team{homeTeam, awayTeam}

	// Store total passing time.
	// If total passing time bigger than 2880 second, match is finished.
	totalPassingTime := 0
	// Store team passing time.
	// This is attack time.
	// If team passing time bigger than 24 second, attack is finished. The other team take the ball.
	teamPassingTime := 0

	// Randomly select the team with the ball.
	teamPlaying := selectRandomTeam(teams)
	// Randomly select the player with the ball in the team with the ball.
	playerIDPlaying := selectRandomPlayerId(0, playerIds[teamPlaying.ID])

	events := make([]model.Event, 0)

	// Continue until the time runs out.
	for totalPassingTime <= 48*60 {
		// Randomly select event. (pass, score, catch the ball etc)
		event := getRandomEventType()

		// Randomly select second for passing time.
		passingTime := getRandomSecond()

		// Added passing time to total passing time.
		totalPassingTime += passingTime

		// Added passing time to team passing time.
		teamPassingTime += passingTime

		// If team passing time bigger than 24 second, attack is finished. The other team take the ball.
		if teamPassingTime > 24 {
			// If team passing time bigger than 24 second, we subtracted the excess time from the elapsed time.
			diff := teamPassingTime - 24
			totalPassingTime -= diff

			// Team passing time reset.
			teamPassingTime = 0

			// Change playing team.
			teamPlaying = changePlayingTeam(teams, teamPlaying)
			// Change playing player.
			playerIDPlaying = selectRandomPlayerId(0, playerIds[teamPlaying.ID])
			continue
		}

		// Create event.
		addedToBeEvent := model.Event{
			GameID:     game.ID,
			PlayerID:   playerIDPlaying,
			TeamID:     teamPlaying.ID,
			Time:       totalPassingTime,
			AttackTime: teamPassingTime,
			Type:       event,
		}
		// Appent event to events.
		events = append(events, addedToBeEvent)

		// If event is score or catch the ball, we need to change the player with the ball.
		if event.DidTheBallChangeSides() {
			// Team passing time reset.
			teamPassingTime = 0

			// Change playing team.
			teamPlaying = changePlayingTeam(teams, teamPlaying)
			// Change playing player.
			playerIDPlaying = selectRandomPlayerId(0, playerIds[teamPlaying.ID])
			continue
		}

		// If event is failed shoot, one of the two teams may have snatched the ball.
		// So randomly select playing team.
		if event.FailedShoot() {
			// Team passing time reset.
			teamPassingTime = 0

			// Randomly select playing team.
			teamPlaying = selectRandomTeam(teams)

			// Randomly select playing player.
			playerIDPlaying = selectRandomPlayerId(0, playerIds[teamPlaying.ID])
			continue
		}

		// If event is pass, we need to change the player with the ball.
		// But we don't select same player.
		if event == model.EventTypeToPass {
			playerIDPlaying = selectRandomPlayerId(playerIDPlaying, playerIds[teamPlaying.ID])
			continue
		}

		if event == model.EventTypeDribblingTheBall {
			continue
		}
	}

	// Firstly event is attack. So we need mark first event is attack.
	events[0].IsAttack = true

	for i := 1; i < len(events); i++ {
		// If event is successful shoot and previous event is pass and previous event team same event team, this is assist.
		// Mark this event is assist.
		if events[i].Type.SuccessfulShoot() &&
			events[i-1].Type == model.EventTypeToPass &&
			events[i].TeamID == events[i-1].TeamID {
			events[i-1].IsAsist = true
		}
		// If event team don't same previous event team or previous event type is failed shoot or previous event type catch the ball this is a new attack.
		// Mark this event is attack.
		if events[i].TeamID != events[i-1].TeamID || events[i-1].Type.FailedShoot() || events[i-1].Type == model.EventTypeCatchTheBall {
			events[i].IsAttack = true
		}
	}

	// Create events.
	err := postgre.DB().Create(&events).Error
	if err != nil {
		panic("Events create error: " + err.Error())
	}
}

// Get team by id in database.
func getTeam(id int) model.Team {
	var team model.Team
	err := postgre.DB().Model(&model.Team{}).Where("id = ?", id).First(&team).Error
	if err != nil {
		panic("Team find error: " + err.Error())
	}
	return team
}

// Get main players ids by team and game id in database.
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

// Change playing team.
func changePlayingTeam(teams []model.Team, teamPlaying *model.Team) *model.Team {
	// We cannot choose the same team. So we check team ids.
	if teams[0].ID == teamPlaying.ID {
		return &teams[1]
	}
	return &teams[0]
}
