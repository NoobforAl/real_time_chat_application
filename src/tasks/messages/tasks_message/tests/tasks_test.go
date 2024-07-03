package tasksMessage

import (
	"context"
	"log"
	"testing"
	"time"

	"github.com/NoobforAl/real_time_chat_application/src/config"
	"github.com/NoobforAl/real_time_chat_application/src/contract"
	"github.com/NoobforAl/real_time_chat_application/src/database"
	"github.com/NoobforAl/real_time_chat_application/src/entity"
	"github.com/NoobforAl/real_time_chat_application/src/logging"
	tasksMessage "github.com/NoobforAl/real_time_chat_application/src/tasks/messages/tasks_message"
	taskTypes "github.com/NoobforAl/real_time_chat_application/src/tasks/tasks_type"
	"github.com/hibiken/asynq"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
)

var ctx = context.TODO()
var store contract.Store

func TestMain(m *testing.M) {
	logger := logging.New()

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	if err := godotenv.Load("./../../../../../.env"); err != nil {
		log.Fatal("not found .env file !!")
	}

	config.InitConfig(logger)

	store = database.New(ctx, logger, database.Opts{
		NeedRedis:      true,
		NeedMongodb:    true,
		NeedBrokerRoom: true,
	})

	go store.SendDailyReportOfMessage(ctx, "@every 5s")
	go store.SendSignalCleanOldMessageAndArchive(ctx, "@every 5s")

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

	go func() { logger.Fatal(srv.Run(mux)) }()
	m.Run()
}

func TestMessageTask(t *testing.T) {
	room := entity.Room{
		Name:        "testroom_random_msgBroker",
		Description: "This is a test room",
		CreateAt:    time.Now(),
	}

	createdRoom, err := store.CreateRoom(ctx, room)
	assert.Nil(t, err, "Error should be nil")
	assert.Equal(t, room.Name, createdRoom.Name, "Name should match")
	assert.Equal(t, room.Description, createdRoom.Description, "Description should match")

	user := entity.User{
		Username: "testuser_random_msgBroker",
		Email:    "testuser@example.com",
		Password: "password",
		CreateAt: time.Now(),
	}

	createdUser, err := store.CreateUser(ctx, user)
	assert.Nil(t, err, "Error should be nil")
	assert.Equal(t, user.Username, createdUser.Username, "Username should match")
	assert.Equal(t, user.Email, createdUser.Email, "Email should match")

	message := entity.Message{
		Content:   "Hello, World!",
		SenderId:  createdUser.Id,
		RoomId:    createdRoom.Id,
		Timestamp: time.Now(),
	}

	t.Run("send new message in Queue & save it & archive message after 10s", func(t *testing.T) {
		for range 3 {
			err = store.SendNewMessage(ctx, message)
			assert.Nil(t, err, "Error should be nil")
		}

		time.Sleep(3 * time.Second)

		messages, err := store.Messages(ctx, createdRoom.Id, 10)
		assert.Nil(t, err, "Error should be nil")
		data := messages[0]

		assert.Equal(t, data.RoomId, createdRoom.Id)
		assert.Equal(t, data.SenderId, createdUser.Id)
		assert.Equal(t, data.Content, message.Content)

		// archive
		time.Sleep(10 * time.Second)
		messages, err = store.AllArchiveMessage(ctx, createdUser.Id)
		assert.Nil(t, err, "Error should be nil")
		if len(messages) < 1 {
			t.Error("not any message archived!!")
		}
	})

}

func TestReport(t *testing.T) {
	// TODO: what's is must be report ??
}
