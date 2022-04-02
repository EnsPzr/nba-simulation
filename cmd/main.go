package main

import (
	"encoding/json"
	"fmt"
	"github.com/enspzr/nba-simulation/data"
	"github.com/enspzr/nba-simulation/router"
	"github.com/enspzr/nba-simulation/storage/postgre"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	fmt.Println("Backend Starting")
	err := postgre.Connect()
	if err != nil {
		panic("Database connection failed => " + err.Error())
	}
	postgre.AutoMigrate()
	data.InitData()

	engine := html.New("./views/", ".html")
	engine.Reload(true)
	app := fiber.New(fiber.Config{
		ReadTimeout:    time.Duration(10) * time.Second,
		WriteTimeout:   time.Duration(10) * time.Second,
		IdleTimeout:    time.Duration(10) * time.Second,
		ReadBufferSize: fiber.DefaultReadBufferSize * 2, // Request Header Fields Too Large hatası için
		JSONEncoder:    json.Marshal,
		Views:          engine,
	})
	router.Setup(app)

	go func() {
		if errr := app.Listen(fmt.Sprintf(":%v", 8080)); errr != nil {
			fmt.Println("Fiber listen error => ", err.Error())
		}
	}()
	fmt.Println("Backend Started")

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	_ = <-c // block until interr
	fmt.Println("Gracefully shutting down")
	_ = app.Shutdown()
}
