package auth

import (
	"context"
	"log"
	"testing"
	"time"

	"github.com/NoobforAl/real_time_chat_application/src/config"
	"github.com/NoobforAl/real_time_chat_application/src/database"
	"github.com/NoobforAl/real_time_chat_application/src/entity"
	"github.com/NoobforAl/real_time_chat_application/src/logging"
	"github.com/NoobforAl/real_time_chat_application/src/services/auth/jwt"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
	grpc "google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var serverTest *server

func TestMain(m *testing.M) {
	err := godotenv.Load("./../../../.env")
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	logger := logging.New()
	ctx := context.Background()

	config.InitConfig(logger)

	mongodbUri := config.MongodbUri()

	redisUri := config.RedisUri()
	redisPassword := config.RedisPassword()

	store := database.New(ctx, mongodbUri, redisUri, redisPassword, logger)

	serverTest = New(store, logger)
	go serverTest.Run(config.GrpcAuthUri())

	// wait for server complete start
	time.Sleep(100 * time.Millisecond)

	m.Run()
}

func TestAuthGrpc(t *testing.T) {
	userData := entity.User{
		Id:       "sample_id",
		Username: "sample_user_name",
	}

	accessToken, _, err := jwt.GenerateTokens([]byte(config.SecretKey()), userData.Id, userData.Username, 10*time.Second, 10*time.Second)
	assert.NoError(t, err)

	conn, err := grpc.NewClient(config.GrpcAuthUri(), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	c := NewAuthServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.TODO(), time.Second)
	defer cancel()

	infoToken, err := c.Login(ctx, &LoginRequest{Token: accessToken})
	assert.NoError(t, err)
	assert.Equal(t, infoToken.GetId(), userData.Id)
	assert.Equal(t, infoToken.GetUsername(), userData.Username)
	assert.Equal(t, infoToken.GetNotification(), userData.Notification)

	_, err = serverTest.Login(context.TODO(), &LoginRequest{Token: accessToken + "some_string"})
	assert.Error(t, err)
}
