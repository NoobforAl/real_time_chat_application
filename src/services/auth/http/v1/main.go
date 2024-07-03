package main

import (
	"context"
	"log"

	"github.com/NoobforAl/real_time_chat_application/src/config"
	"github.com/NoobforAl/real_time_chat_application/src/database"
	"github.com/NoobforAl/real_time_chat_application/src/logging"
	"github.com/NoobforAl/real_time_chat_application/src/services/auth/http/v1/middleware"
	"github.com/NoobforAl/real_time_chat_application/src/services/auth/http/v1/router"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

func main() {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	app := fiber.New()
	logger := logging.New()
	config.InitConfig(logger)

	store := database.New(ctx, logger)

	logLogrusType, ok := logger.(*logrus.Logger)
	if !ok {
		log.Fatal(
			"For this app we need init logrus logger," +
				" but your interface use other logging pkg!",
		)
	}

	app.Use(middleware.SetupLoggingMiddleware(logLogrusType))

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("<h1>Auth Service Is Working</h1>")
	})

	router.SetupAuthRoute(app, store, logger)

	addr := config.AuthServiceUri()
	logger.Fatal(app.Listen(addr))
}
