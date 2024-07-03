package controller

import (
	"context"
	"time"

	"github.com/NoobforAl/real_time_chat_application/src/contract"
	"github.com/NoobforAl/real_time_chat_application/src/entity"
	"github.com/NoobforAl/real_time_chat_application/src/validation"
	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
)

func WebSocketMessages(ctx context.Context, store contract.Store, log contract.Logger) func(*fiber.Ctx) error {
	return websocket.New(func(c *websocket.Conn) {
		token := c.Headers("Access-Token", "")
		user, err := store.Login(ctx, token)
		if err != nil {
			log.Error(err)
			_ = c.WriteJSON(&fiber.Map{"error": err.Error()})
			return
		}

		roomId := c.Params("room_id")
		_, err = store.Room(ctx, roomId)
		if err != nil {
			log.Error(err)
			_ = c.WriteJSON(fiber.Map{"error": err.Error()})
			return
		}

		roomSes := getRoomSession(roomId, 0, log)
		roomSes.addUser(&userSession{
			conn:         c,
			userId:       user.Id,
			username:     user.Username,
			notification: user.Notification,
		})

		defer log.Debugf("wb messages: connection close remove user with id: %s", user.Id)
		defer roomSes.delUser(user.Id)
		defer c.Close()

		var message entity.Message

		for {
			if err = c.ReadJSON(&message); err != nil {
				err = c.WriteJSON(fiber.Map{"error": err.Error()})
				if err != nil {
					log.Error(err)
					return
				}

				continue
			}

			log.Debugf("receive new message: %v", message)

			if errors := validation.ValidateStruct(message); errors != nil {
				err = c.WriteJSON(errors)
				if err != nil {
					log.Error(err)
					return
				}

				log.Error(errors)
				continue
			}

			message.RoomId = roomId
			message.SenderId = user.Id
			message.Timestamp = time.Now()

			if err = store.SendNewMessage(ctx, message); err != nil {
				log.Error(err)

				err = c.WriteJSON(fiber.Map{"error": err.Error()})
				if err != nil {
					log.Error(err)
					return
				}

				continue
			}

			log.Debug("notify all users with notification")
			roomSes.notify(message, message.SenderId)

			log.Debug("send message again for user")
			if err = c.WriteJSON(message); err != nil {
				log.Error(err)
				return
			}
		}
	})
}

func WebSocketNotifications(ctx context.Context, store contract.Store, log contract.Logger) func(*fiber.Ctx) error {
	return websocket.New(func(c *websocket.Conn) {
		token := c.Headers("Access-Token", "")
		user, err := store.Login(ctx, token)
		if err != nil {
			log.Error(err)
			_ = c.WriteJSON(&fiber.Map{"error": err.Error()})
			return
		}

		roomId := c.Params("room_id")
		_, err = store.Room(ctx, roomId)
		if err != nil {
			log.Error(err)
			_ = c.WriteJSON(fiber.Map{"error": err.Error()})
			return
		}

		roomSes := getRoomSession(roomId, 1, log)
		roomSes.addUser(&userSession{
			conn:         c,
			userId:       user.Id,
			username:     user.Username,
			notification: user.Notification,
		})

		ctx, cancel := context.WithCancel(ctx)
		defer cancel()

		defer log.Debugf("wb notification: connection close remove user with id: %s", user.Id)
		defer roomSes.delUser(user.Id)
		defer c.Close()

		go roomSes.turnOffOrOnNotification(ctx, user.Id, roomId)

		for {
			notification, err := store.GetNotification(ctx, roomId)
			if err != nil {
				log.Error(err)
				continue
			}

			log.Debugf("get new notification room id=%s, message is: %v", roomId, notification)
			roomSes.notify(notification, notification.SenderId)
		}
	})
}
