package router

import "github.com/gofiber/fiber/v2"

/*
POST /register/ - Register a new user.
• POST /login/ - Authenticate a user and return a JWT token.
*/

func SetupAuthRoute(app *fiber.App) {
	app.Post("/login/")
	app.Post("/register/")
}
