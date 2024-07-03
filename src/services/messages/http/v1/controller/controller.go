package controller

import (
	"time"

	"github.com/NoobforAl/real_time_chat_application/src/contract"
	"github.com/NoobforAl/real_time_chat_application/src/entity"
	"github.com/NoobforAl/real_time_chat_application/src/validation"
	"github.com/gofiber/fiber/v2"
)

func GetMessagesRoom(store contract.Store, log contract.Logger) func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		ctx := c.Context()

		roomId := c.Params("room_id", "")
		if roomId == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "room Id param not empty string!",
			})
		}

		messages, err := store.Messages(ctx, roomId, 100)
		if roomId == "" {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}

		return c.Status(fiber.StatusOK).JSON(messages)
	}
}

func CreateNewMessage(store contract.Store, log contract.Logger) func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		ctx := c.Context()
		token := c.Get("Access-Token", "")
		roomId := c.Params("room_id", "")

		user, err := store.Login(ctx, token)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
		}

		_, err = store.Room(ctx, roomId)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
		}

		var message entity.Message
		if err := c.BodyParser(&message); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "cannot parse JSON"})
		}

		if errors := validation.ValidateStruct(message); errors != nil {
			return c.Status(fiber.StatusBadRequest).JSON(errors)
		}

		message.RoomId = roomId
		message.SenderId = user.Id
		message.Timestamp = time.Now()

		message, err = store.CreateMessage(ctx, message)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		return c.Status(fiber.StatusOK).JSON(message)
	}
}
