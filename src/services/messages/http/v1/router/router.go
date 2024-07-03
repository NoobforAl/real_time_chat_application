package router

import (
	"context"

	"github.com/NoobforAl/real_time_chat_application/src/contract"
	"github.com/NoobforAl/real_time_chat_application/src/services/messages/http/v1/controller"
	"github.com/NoobforAl/real_time_chat_application/src/services/messages/http/v1/middleware"
	"github.com/gofiber/fiber/v2"
)

func SetupMessageRoute(ctx context.Context, app *fiber.App, store contract.Store, log contract.Logger) {
	app.Get("/messages/ws/:room_id", controller.WebSocketMessages(ctx, store, log))
	app.Get("/notification/ws/:room_id", controller.WebSocketNotifications(ctx, store, log))

	app.Post("/messages/:room_id", controller.CreateNewMessage(store, log))

	app.Use(middleware.CheckJwtToken(store, log))
	app.Get("/messages/:room_id", controller.GetMessagesRoom(store, log))
}
