package contract

import (
	"context"

	"github.com/NoobforAl/real_time_chat_application/src/entity"
)

type Cache interface {
	GetRooms(ctx context.Context) ([]*entity.Room, error)
	SetRooms(ctx context.Context, data []*entity.Room) error

	GetMessages(ctx context.Context, roomId string) ([]*entity.Message, error)
	SetMessage(ctx context.Context, roomId string, data []*entity.Message) error

	GetNotification(ctx context.Context, roomId string) (entity.Notification, error)
	GetNotifications(ctx context.Context, roomId string) ([]*entity.Notification, error)
	GetNotificationsRoom(ctx context.Context, userId string) ([]*entity.NotificationRoom, error)

	SetNotifications(ctx context.Context, roomId string, data []*entity.Notification) error
	SetNotificationsRoom(ctx context.Context, userId string, data []*entity.NotificationRoom) error
}
