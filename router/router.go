package router

import (
	"github.com/enspzr/nba-simulation/ws"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
)

// Setup This function register to router all routes.
func Setup(app fiber.Router) {
	app.Use("/ws", func(c *fiber.Ctx) error {
		// IsWebSocketUpgrade returns true if the client
		// requested upgrade to the WebSocket protocol.
		if websocket.IsWebSocketUpgrade(c) {
			c.Locals("allowed", true)
			return c.Next()
		}
		return fiber.ErrUpgradeRequired
	})
	// This endpoint is web socket endpoint.
	app.Get("/ws", websocket.New(ws.WsHandler))
	// This endpoint is a static file endpoint. Only return index.html.
	app.Static("/", "./views/index.html")
}
