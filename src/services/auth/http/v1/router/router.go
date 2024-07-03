package router

import (
	"github.com/NoobforAl/real_time_chat_application/src/contract"
	"github.com/NoobforAl/real_time_chat_application/src/services/auth/http/v1/controller"
	"github.com/gofiber/fiber/v2"
)

func SetupAuthRoute(app *fiber.App, store contract.StoreUser, log contract.Logger) {
	// TODO: add refresh token!

	app.Post("/login/", controller.Login(store, log))

	app.Post("/register/", controller.Register(store, log))
}
