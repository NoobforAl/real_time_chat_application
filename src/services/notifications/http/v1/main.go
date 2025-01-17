package main

import (
	"github.com/NoobforAl/real_time_chat_application/src/config"
	"github.com/NoobforAl/real_time_chat_application/src/logging"
	"github.com/NoobforAl/real_time_chat_application/src/services/notifications/http/v1/router"
	"github.com/gofiber/fiber/v2"
)

// TODO: fix it later, I implementation it in messages service!
func main() {
	app := fiber.New()
	logger := logging.New()
	config.InitConfig(logger)

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("<h1>Notifications Service Is Working</h1>")
	})

	router.SetupNotificationsRoute(app)

	addr := config.AuthServiceUri()
	logger.Fatal(app.Listen(addr))
}
