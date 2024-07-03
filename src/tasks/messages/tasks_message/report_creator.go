package tasksMessage

import (
	"context"

	"github.com/NoobforAl/real_time_chat_application/src/contract"
	taskTypes "github.com/NoobforAl/real_time_chat_application/src/tasks/tasks_type"
	"github.com/hibiken/asynq"
)

type MessageReportProcess struct {
	store contract.Store
	log   contract.Logger
}

func NewMessageReportTask() (*asynq.Task, error) {
	return asynq.NewTask(taskTypes.TypeMessageReportDaily, []byte("")), nil
}

func NewMessageReportProcess(store contract.Store, log contract.Logger) *MessageReportProcess {
	return &MessageReportProcess{store: store, log: log}
}

func (msp *MessageReportProcess) ProcessTask(ctx context.Context, t *asynq.Task) error {
	userIds, err := msp.store.UserIds(ctx)
	if err != nil {
		return err
	}

	// Some Process Report // TODO
	for _, v := range userIds {
		msp.log.Debugf("do something data process for user id: %s", v)
		// do it ...
		// _, _ = msp.store.UserMessagesInOneDay(ctx, v)
		// process data message on day

		// send notification
		// msp.store.SetNotifications(/* data */)
	}

	return nil
}
