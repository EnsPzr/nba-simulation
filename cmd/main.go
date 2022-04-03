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

// This function is called when the program is started.
// It initializes the database and starts the server.
func main() {
	fmt.Println("Backend Starting")
	//Connect database.
	err := postgre.Connect()
	if err != nil {
		panic("Database connection failed => " + err.Error())
	}
	//Migrate database.
	postgre.AutoMigrate()
	//Seeding database.
	data.InitData()

	engine := html.New("./views/", ".html")
	engine.Reload(true)
	// Create new Fiber web server instance.
	app := fiber.New(fiber.Config{
		ReadTimeout:    time.Duration(10) * time.Second,
		WriteTimeout:   time.Duration(10) * time.Second,
		IdleTimeout:    time.Duration(10) * time.Second,
		ReadBufferSize: fiber.DefaultReadBufferSize * 2, // Request Header Fields Too Large hatası için
		JSONEncoder:    json.Marshal,
		Views:          engine,
	})
	// Register routes.
	router.Setup(app)

	// Start server.
	go func() {
		if err = app.Listen(fmt.Sprintf(":%v", 8080)); err != nil {
			fmt.Println("Fiber listen error => ", err.Error())
		}
	}()
	fmt.Println("Backend Started")

	// Wait for interrupt signal to gracefully shutdown the server with
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	_ = <-c // block until interr
	fmt.Println("Gracefully shutting down")
	_ = app.Shutdown()
}
