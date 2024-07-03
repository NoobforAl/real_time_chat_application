package controller

import (
	"time"

	"github.com/NoobforAl/real_time_chat_application/src/contract"
	"github.com/NoobforAl/real_time_chat_application/src/entity"
	"github.com/NoobforAl/real_time_chat_application/src/validation"
	"github.com/gofiber/fiber/v2"
)

func AllRooms(store contract.StoreRoom, log contract.Logger) func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		ctx := c.Context()
		rooms, err := store.Rooms(ctx, 100)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to get all room"})
		}

		return c.Status(fiber.StatusOK).JSON(rooms)
	}
}

func CreateNewRoom(store contract.StoreRoom, log contract.Logger) func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		var room entity.Room
		if err := c.BodyParser(&room); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "cannot parse JSON"})
		}

		if errors := validation.ValidateStruct(room); errors != nil {
			return c.Status(fiber.StatusBadRequest).JSON(errors)
		}

		room.CreateAt = time.Now()

		ctx := c.Context()
		room, err := store.CreateRoom(ctx, room)
		if err != nil {
			log.Error(err)
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "failed to create new room"})
		}

		return c.Status(fiber.StatusOK).JSON(room)
	}
}
