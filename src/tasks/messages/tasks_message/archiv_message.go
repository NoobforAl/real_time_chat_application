package tasksMessage

import (
	"context"
	"time"

	"github.com/NoobforAl/real_time_chat_application/src/contract"
	taskTypes "github.com/NoobforAl/real_time_chat_application/src/tasks/tasks_type"
	"github.com/hibiken/asynq"
)

type MessageArchiveProcess struct {
	store contract.Store
	log   contract.Logger
}

func NewMessageArchiveTask() (*asynq.Task, error) {
	return asynq.NewTask(taskTypes.TypeMessageArchive, []byte("")), nil
}

func NewMessageArchiveProcess(store contract.Store, log contract.Logger) *MessageArchiveProcess {
	return &MessageArchiveProcess{store: store, log: log}
}

func (msp *MessageArchiveProcess) ProcessTask(ctx context.Context, t *asynq.Task) error {
	userIds, err := msp.store.UserIds(ctx)
	if err != nil {
		return err
	}

	// last 30 day message
	endTime := time.Now()
	startTime := time.Now().Add(-24 * 30 * time.Hour)

	// Some Process Report // TODO
	for _, userId := range userIds {
		msp.log.Debugf("process user id: %s, save archive message to other collection", userId)
		messageUser, err := msp.store.UserMessages(ctx, userId, startTime, endTime)
		if err != nil {
			msp.log.Error(err)
			return err
		}

		err = msp.store.SaveMessageToArchive(ctx, userId, messageUser)
		if err != nil {
			msp.log.Error(err)
			return err
		}
	}

	return nil
}
