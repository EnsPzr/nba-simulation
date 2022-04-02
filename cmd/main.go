package main

import (
	"fmt"
	"github.com/enspzr/nba-simulation/data"
	"github.com/enspzr/nba-simulation/storage/postgre"
)

func main() {
	fmt.Println("Backend Starting")
	err := postgre.Connect()
	if err != nil {
		panic("Database connection failed => " + err.Error())
	}
	postgre.AutoMigrate()
	data.InitData()
	fmt.Println("Backend Started")
}
