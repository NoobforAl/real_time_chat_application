package tasksMessage

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/NoobforAl/real_time_chat_application/src/contract"
	"github.com/NoobforAl/real_time_chat_application/src/entity"
	taskTypes "github.com/NoobforAl/real_time_chat_application/src/tasks/tasks_type"
	"github.com/hibiken/asynq"
)

type MessageSaveProcess struct {
	store contract.StoreMessageAndCache
	log   contract.Logger
}

func NewMessageSaveTask(message entity.Message) (*asynq.Task, error) {
	msg, err := json.Marshal(message)
	if err != nil {
		return nil, err
	}

	return asynq.NewTask(taskTypes.TypeMessageSave, msg), nil
}

func NewMessageSaveProcess(store contract.StoreMessageAndCache, log contract.Logger) *MessageSaveProcess {
	return &MessageSaveProcess{store: store, log: log}
}

func (msp *MessageSaveProcess) ProcessTask(ctx context.Context, t *asynq.Task) error {
	dataMessage := t.Payload()
	msp.log.Debugf("process save message, message content: %v", dataMessage)

	var message entity.Message
	err := json.Unmarshal(dataMessage, &message)
	if err != nil {
		return err
	}

	// set in queue for show with websocket to online user
	err = msp.store.SetMessage(ctx, message.RoomId, []*entity.Message{&message})
	if err != nil {
		return err
	}

	_, err = msp.store.CreateMessage(ctx, message)
	if err != nil {
		return err
	}

	notifContent := fmt.Sprintf("received message: %s", message.Content)
	_ = msp.store.SetNotifications(ctx, message.RoomId, []*entity.Notification{
		{
			SenderId: message.SenderId,
			RoomId:   message.RoomId,
			Content:  notifContent,
			CreateAt: time.Now(),
		},
	})

	return nil
}
