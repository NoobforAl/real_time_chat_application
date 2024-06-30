package contract

import (
	"context"
	"time"

	"github.com/NoobforAl/real_time_chat_application/src/entity"
)

type Cache interface {
	GetRooms(ctx context.Context) ([]*entity.Room, error)
	SetRooms(ctx context.Context, data []*entity.Room) error

	GetMessages(ctx context.Context, room_id string) ([]*entity.Message, error)
	SetMessage(ctx context.Context, room_id string) ([]*entity.Message, error)

	GetNotification(ctx context.Context, room_id string) (entity.Notification, error)
	GetNotifications(ctx context.Context, room_id string) ([]*entity.Notification, error)
	GetNotificationsRoom(ctx context.Context, user_id string) ([]*entity.Notification, error)

	SetNotifications(ctx context.Context, room_id string, data []*entity.Notification) error
	SetNotificationsRoom(ctx context.Context, user_id string, data []*entity.NotificationRoom) error
}

type StoreUser interface {
	User(ctx context.Context, username string) (entity.User, error)
	CreateUser(ctx context.Context, userData entity.User) (entity.User, error)
	UpdateUser(ctx context.Context, id string, updatedData entity.User) (entity.User, error)
}

type StoreRoom interface {
	Rooms(ctx context.Context, maxLen int) ([]*entity.Room, error)
	CreateRoom(ctx context.Context, message entity.Room) (entity.Room, error)
}

type StoreMessage interface {
	Messages(ctx context.Context, roomId string, maxLen int) ([]*entity.Message, error)
	CreateMessage(ctx context.Context, message entity.Message) (entity.Message, error)
}

type StoreNotifications interface {
	Notifications(ctx context.Context, maxLen int) ([]*entity.Notification, error)
	CreateNotification(ctx context.Context, message entity.Notification) (entity.Notification, error)

	NotificationRoom(ctx context.Context, notificationRoomData entity.NotificationRoom) (entity.NotificationRoom, error)
	CreateNotificationRoom(ctx context.Context, notificationRoomData entity.NotificationRoom) (entity.NotificationRoom, error)
	UpdateNotificationRoom(ctx context.Context, message entity.NotificationRoom) (entity.NotificationRoom, error)
}

type AuthenticationService interface {
	Login(ctx context.Context, token string) entity.User
}

type BrokerRoom interface {
	SendNewRoom(ctx context.Context, message entity.Room) error
}

type BrokerMessage interface {
	SendNewMessage(ctx context.Context, message entity.Message) error
	SendDailyReportOfMessage(ctx context.Context, user_id string) error
	SendSignalCleanOldMessageAndArchive(ctx context.Context, timeBefore time.Time) error
}

type BrokerNotification interface {
	SendNewNotification(ctx context.Context, message entity.Notification) error
}

type Store interface {
	Cache

	StoreUser
	StoreRoom
	StoreMessage
	StoreNotifications

	AuthenticationService

	BrokerRoom
	BrokerMessage
	BrokerNotification
}
