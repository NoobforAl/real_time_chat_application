package main

import (
	"context"
	"log"

	"github.com/NoobforAl/real_time_chat_application/src/config"
	"github.com/NoobforAl/real_time_chat_application/src/database"
	"github.com/NoobforAl/real_time_chat_application/src/logging"
	"github.com/NoobforAl/real_time_chat_application/src/services/auth/http/v1/router"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

func setupLoggingMiddleware(log *logrus.Logger) func(*fiber.Ctx) error {
	return func(ctx *fiber.Ctx) error {
		log.WithFields(logrus.Fields{
			"method": ctx.Method(),
			"path":   ctx.Path(),
			"ip":     ctx.IP(),
		}).Info("Request received")

		return ctx.Next()
	}
}

func main() {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	app := fiber.New()
	logger := logging.New()
	config.InitConfig(logger)

	mongodbUri := config.MongodbUri()
	redisUri := config.RedisUri()
	redisPassword := config.RedisPassword()
	store := database.New(ctx, mongodbUri, redisUri, redisPassword, logger)

	logLogrusType, ok := logger.(*logrus.Logger)
	if !ok {
		log.Fatal(
			"For this app we need init logrus logger," +
				" but your interface use other logging pkg!",
		)
	}

	app.Use(setupLoggingMiddleware(logLogrusType))

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("<h1>Auth Service Is Working</h1>")
	})

	router.SetupAuthRoute(ctx, app, store, logger)

	addr := config.AuthServiceUri()
	logger.Fatal(app.Listen(addr))
}
