package router

import "github.com/gofiber/fiber/v2"

func SetupRoomRoute(app *fiber.App) {
	app.Get("/ws/:room_id")
}
