package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

func SetupLoggingMiddleware(log *logrus.Logger) func(*fiber.Ctx) error {
	return func(ctx *fiber.Ctx) error {
		log.WithFields(logrus.Fields{
			"time":   ctx.Context().Time(),
			"method": ctx.Method(),
			"path":   ctx.Path(),
			"ip":     ctx.IP(),
		}).Info("Request received")

		return ctx.Next()
	}
}
