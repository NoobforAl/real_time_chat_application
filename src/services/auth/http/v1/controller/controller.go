package controller

import (
	"context"
	"time"

	"github.com/NoobforAl/real_time_chat_application/src/config"
	"github.com/NoobforAl/real_time_chat_application/src/contract"
	"github.com/NoobforAl/real_time_chat_application/src/entity"
	"github.com/NoobforAl/real_time_chat_application/src/services/auth/jwt"
	"github.com/NoobforAl/real_time_chat_application/src/validation"
	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
)

func Login(ctx context.Context, store contract.StoreUser, log contract.Logger) func(*fiber.Ctx) error {
	secretKey := config.SecretKey()
	maxAgeToken := config.MaxAgeToken()

	salt := config.NoncForHashPassword()

	return func(c *fiber.Ctx) error {
		var loginInfo = struct {
			Username string `json:"username" validate:"required,max=64,min=8"`
			Password string `json:"password" validate:"required,max=64,min=8"`
		}{}

		if err := c.BodyParser(&loginInfo); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "cannot parse JSON"})
		}

		if errors := validation.ValidateStruct(loginInfo); errors != nil {
			return c.Status(fiber.StatusBadRequest).JSON(errors)
		}

		user, err := store.User(ctx, loginInfo.Username)
		if err != nil {
			log.Error(err)
			return c.Status(fiber.StatusUnauthorized).JSON(&fiber.Error{
				Code:    fiber.StatusUnauthorized,
				Message: "incorrect username & password",
			})
		}

		err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(salt+loginInfo.Password))
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "incorrect username & password"})
		}

		accessToken, refreshToken, err := jwt.GenerateTokens([]byte(secretKey), user.Id, user.Username, maxAgeToken, 24*time.Hour)

		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(&fiber.Error{
				Code:    fiber.StatusUnauthorized,
				Message: err.Error(),
			})
		}

		return c.JSON(fiber.Map{
			"access_token":  accessToken,
			"refresh_token": refreshToken,
		})
	}
}

func Register(ctx context.Context, store contract.StoreUser, log contract.Logger) func(*fiber.Ctx) error {
	salt := config.NoncForHashPassword()

	return func(c *fiber.Ctx) error {
		var user entity.User
		if err := c.BodyParser(&user); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "cannot parse JSON"})
		}

		if errors := validation.ValidateStruct(user); errors != nil {
			return c.Status(fiber.StatusBadRequest).JSON(errors)
		}

		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(salt+user.Password), bcrypt.DefaultCost)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to hash password"})
		}

		user.Password = string(hashedPassword)
		user.CreateAt = time.Now()
		user.Notification = true

		user, err = store.CreateUser(ctx, user)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(&fiber.Error{
				Code:    fiber.StatusInternalServerError,
				Message: err.Error(),
			})
		}

		return c.Status(fiber.StatusOK).JSON(user)
	}
}
