package router

import "github.com/gofiber/fiber/v2"

/*
GET /messages/ - List all messages.
• POST /messages/ - Create a new message.
*/

func SetupAuthRoute(app *fiber.App) {
	app.Get("/messages/")
	app.Post("/messages/")
}
