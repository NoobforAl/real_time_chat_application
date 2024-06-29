package main

import (
	"github.com/NoobforAl/real_time_chat_application/src/config"
	"github.com/NoobforAl/real_time_chat_application/src/logging"
	"github.com/NoobforAl/real_time_chat_application/src/services/rooms/http/v1/router"
	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New()
	logger := logging.New()
	config.InitConfig(logger)

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("<h1>Rooms Service Is Working</h1>")
	})

	router.SetupRoomRoute(app)

	addr := config.AuthServiceUri()
	logger.Fatal(app.Listen(addr))
}
