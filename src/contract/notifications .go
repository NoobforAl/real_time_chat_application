package contract

import (
	"context"

	"github.com/NoobforAl/real_time_chat_application/src/entity"
)

type StoreNotifications interface {
	Notifications(ctx context.Context, maxLen int) ([]*entity.Notification, error)
	CreateNotification(ctx context.Context, notifications entity.Notification) (entity.Notification, error)

	CreateNotificationRoom(ctx context.Context, notificationRoomData entity.NotificationRoom) (entity.NotificationRoom, error)
	UpdateNotificationRoom(ctx context.Context, message entity.NotificationRoom) (entity.NotificationRoom, error)
}

type BrokerNotification interface {
	SendNewNotification(ctx context.Context, message entity.Notification) error
}

type StoreNotificationsAndCache interface {
	Cache
	StoreNotifications
}
