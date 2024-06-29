package router

import "github.com/gofiber/fiber/v2"

/*
POST /register/ - Register a new user.
• POST /login/ - Authenticate a user and return a JWT token.
*/

func SetupNotificationsRoute(app *fiber.App) {
	app.Get("/ws/:user_id")
}
