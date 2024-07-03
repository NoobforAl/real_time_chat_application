package contract

import (
	"context"

	"github.com/NoobforAl/real_time_chat_application/src/entity"
)

type StoreRoom interface {
	Room(ctx context.Context, roomId string) (entity.Room, error)
	Rooms(ctx context.Context, maxLen int) ([]*entity.Room, error)
	CreateRoom(ctx context.Context, room entity.Room) (entity.Room, error)
}

type BrokerRoom interface {
	SendNewRoom(ctx context.Context, room entity.Room) error
}

type StoreRoomsAndCache interface {
	Cache
	StoreRoom
}
