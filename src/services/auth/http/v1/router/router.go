package router

import (
	"context"

	"github.com/NoobforAl/real_time_chat_application/src/contract"
	"github.com/NoobforAl/real_time_chat_application/src/services/auth/http/v1/controller"
	"github.com/gofiber/fiber/v2"
)

func SetupAuthRoute(ctx context.Context, app *fiber.App, store contract.StoreUser, log contract.Logger) {
	// TODO: add refresh token!

	app.Post("/login/", controller.Login(ctx, store, log))

	app.Post("/register/", controller.Register(ctx, store, log))
}
