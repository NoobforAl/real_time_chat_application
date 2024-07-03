package main

import (
	"context"

	"github.com/NoobforAl/real_time_chat_application/src/config"
	"github.com/NoobforAl/real_time_chat_application/src/database"
	"github.com/NoobforAl/real_time_chat_application/src/grpc/auth"
	"github.com/NoobforAl/real_time_chat_application/src/logging"
)

func main() {
	logger := logging.New()
	ctx := context.Background()

	config.InitConfig(logger)
	store := database.New(ctx, logger)

	server := auth.New(store, logger)
	server.Run(config.GrpcAuthUri())
}
