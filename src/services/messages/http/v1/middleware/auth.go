package middleware

import (
	"github.com/NoobforAl/real_time_chat_application/src/contract"
	"github.com/gofiber/fiber/v2"
)

func CheckJwtToken(store contract.AuthenticationService, log contract.Logger) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		token := c.Get("Access-Token", "")
		_, err := store.Login(c.Context(), token)
		if err != nil {
			log.Error(err)
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"message": "can't invalid token!",
			})
		}
		return c.Next()
	}
}
