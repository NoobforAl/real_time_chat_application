package router

import (
	"github.com/NoobforAl/real_time_chat_application/src/contract"
	"github.com/NoobforAl/real_time_chat_application/src/services/rooms/http/v1/controller"
	"github.com/NoobforAl/real_time_chat_application/src/services/rooms/http/v1/middleware"
	"github.com/gofiber/fiber/v2"
)

func SetupRoomRoute(app *fiber.App, store contract.Store, log contract.Logger) {
	app.Use(middleware.CheckJwtToken(store, log))
	app.Get("/rooms/", controller.AllRooms(store, log))
	app.Post("/rooms/", controller.CreateNewRoom(store, log))
}
