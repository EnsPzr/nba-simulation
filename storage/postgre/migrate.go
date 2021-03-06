package postgre

import (
	"github.com/enspzr/nba-simulation/model"
)

// AutoMigrate This function use to migrate database.
func AutoMigrate() {
	migrate(&model.Group{})
	migrate(&model.Player{})
	migrate(&model.Team{})
	migrate(&model.Game{})
	migrate(&model.GamePlayer{})
	migrate(&model.Event{})
}

func migrate[T any](value T) {
	if err := db.AutoMigrate(&value); err != nil {
		panic("Failed to migrate => " + err.Error())
	}
}
