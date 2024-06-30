package contract

import (
	"context"

	"github.com/NoobforAl/real_time_chat_application/src/entity"
)

type StoreUser interface {
	User(ctx context.Context, username string) (entity.User, error)
	CreateUser(ctx context.Context, userData entity.User) (entity.User, error)
	UpdateUser(ctx context.Context, id string, updatedData entity.User) (entity.User, error)
}

type StoreRoom interface {
	Messages(ctx context.Context, roomId string, maxLen int) ([]*entity.Message, error)
	CreateMessage(ctx context.Context, message entity.Message) (entity.Message, error)
}
type StoreMessage interface {
	Rooms(ctx context.Context, maxLen int) ([]*entity.Room, error)
	CreateRoom(ctx context.Context, message entity.Room) (entity.Room, error)
}

type StoreNotifications interface {
	Notifications(ctx context.Context, maxLen int) ([]*entity.Notification, error)
	CreateNotification(ctx context.Context, message entity.Notification) (entity.Notification, error)

	CreateNotificationRoom(ctx context.Context, notificationRoomData entity.NotificationRoom) (entity.NotificationRoom, error)
	UpdateNotificationRoom(ctx context.Context, message entity.NotificationRoom) (entity.NotificationRoom, error)
}

type Store interface {
	StoreUser
	StoreRoom
	StoreMessage
	StoreNotifications
}
