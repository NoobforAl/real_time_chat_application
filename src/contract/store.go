package contract

import (
	"context"

	"github.com/NoobforAl/real_time_chat_application/src/entity"
)

type Store interface {
	User(ctx context.Context, id string) entity.User
	CreateUser(ctx context.Context, userData entity.User) entity.User

	Messages(ctx context.Context, roomId string, maxLen int) []*entity.Message
	CreateMessage(ctx context.Context, message entity.Message) entity.Message

	Rooms(ctx context.Context, maxLen int) []*entity.Room
	CreateRoom(ctx context.Context, message entity.Room) entity.Room
}
