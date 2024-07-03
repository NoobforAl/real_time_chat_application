package taskNotification

import (
	"context"
	"encoding/json"

	"github.com/NoobforAl/real_time_chat_application/src/contract"
	"github.com/NoobforAl/real_time_chat_application/src/entity"
	taskTypes "github.com/NoobforAl/real_time_chat_application/src/tasks/tasks_type"
	"github.com/hibiken/asynq"
)

type NotificationSaveProcess struct {
	store contract.StoreNotificationsAndCache
	log   contract.Logger
}

func NewNotificationSaveTask(message entity.Notification) (*asynq.Task, error) {
	msg, err := json.Marshal(message)
	if err != nil {
		return nil, err
	}

	return asynq.NewTask(taskTypes.TypeNotificationSave, msg), nil
}

func NewNotificationSaveProcess(store contract.StoreNotificationsAndCache, log contract.Logger) *NotificationSaveProcess {
	return &NotificationSaveProcess{store: store, log: log}
}

func (msp *NotificationSaveProcess) ProcessTask(ctx context.Context, t *asynq.Task) error {
	dataMessage := t.Payload()
	msp.log.Debugf("process save message, message content: %v", dataMessage)

	var notification entity.Notification
	err := json.Unmarshal(dataMessage, &notification)
	if err != nil {
		return err
	}

	// set in queue for show with websocket to online user
	err = msp.store.SetNotifications(ctx, notification.RoomId, []*entity.Notification{&notification})
	if err != nil {
		return err
	}

	_, err = msp.store.CreateNotification(ctx, notification)
	return err
}
