package contract

import (
	"context"
	"time"

	"github.com/NoobforAl/real_time_chat_application/src/entity"
)

type StoreMessage interface {
	Messages(ctx context.Context, roomId string, maxLen int) ([]*entity.Message, error)
	UserMessages(ctx context.Context, userId string, startTime, endTime time.Time) ([]*entity.Message, error)
	CreateMessage(ctx context.Context, message entity.Message) (entity.Message, error)
	SaveMessageToArchive(ctx context.Context, userId string, message []*entity.Message) error
}

type BrokerMessage interface {
	SendNewMessage(ctx context.Context, message entity.Message) error
	SendDailyReportOfMessage(ctx context.Context) error
	SendSignalCleanOldMessageAndArchive(ctx context.Context) error
}

type StoreMessageAndCache interface {
	Cache
	StoreMessage
}
