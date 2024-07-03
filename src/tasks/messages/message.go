package main

import (
	"context"

	"github.com/NoobforAl/real_time_chat_application/src/config"
	"github.com/NoobforAl/real_time_chat_application/src/database"
	"github.com/NoobforAl/real_time_chat_application/src/logging"
	tasksMessage "github.com/NoobforAl/real_time_chat_application/src/tasks/messages/tasks_message"
	taskTypes "github.com/NoobforAl/real_time_chat_application/src/tasks/tasks_type"
	"github.com/hibiken/asynq"
)

func main() {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	logger := logging.New()
	config.InitConfig(logger)

	store := database.New(ctx, logger, database.Opts{
		NeedRedis:   true,
		NeedMongodb: true,
	})

	srv := asynq.NewServer(
		asynq.RedisClientOpt{
			Addr:     config.RedisUri(),
			Password: config.RedisPassword(),
		},
		asynq.Config{Concurrency: 100},
	)

	mux := asynq.NewServeMux()
	mux.Handle(taskTypes.TypeMessageSave, tasksMessage.NewMessageSaveProcess(store, logger))
	mux.Handle(taskTypes.TypeMessageReportDaily, tasksMessage.NewMessageReportProcess(store, logger))
	mux.Handle(taskTypes.TypeMessageArchive, tasksMessage.NewMessageArchiveProcess(store, logger))

	if err := srv.Run(mux); err != nil {
		logger.Fatal(err)
	}
}
