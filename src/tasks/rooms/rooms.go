package main

import (
	"context"

	"github.com/NoobforAl/real_time_chat_application/src/config"
	"github.com/NoobforAl/real_time_chat_application/src/database"
	"github.com/NoobforAl/real_time_chat_application/src/logging"
	taskRoom "github.com/NoobforAl/real_time_chat_application/src/tasks/rooms/tasks_rooms"
	taskTypes "github.com/NoobforAl/real_time_chat_application/src/tasks/tasks_type"
	"github.com/hibiken/asynq"
)

// TODO: need this service next time
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
	mux.Handle(taskTypes.TypeRoomSave, taskRoom.NewRoomSaveProcess(store, logger))

	if err := srv.Run(mux); err != nil {
		logger.Fatal(err)
	}
}
