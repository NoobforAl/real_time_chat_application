// !Warn: before run this unit test need to run database test!
// for after run a test
package database

import (
	"context"
	"testing"
	"time"

	"github.com/NoobforAl/real_time_chat_application/src/entity"
	"github.com/NoobforAl/real_time_chat_application/src/logging"
	"github.com/stretchr/testify/assert"
)

var ctx = context.TODO()

func TestMain(m *testing.M) {
	mongoUri := "mongodb://localhost:27017"
	redisUri := "localhost:6379"
	redisPassword := ""
	logger := logging.New()

	_ = New(ctx, mongoUri, redisUri, redisPassword, logger)

	m.Run()
}

func TestUserOperations(t *testing.T) {
	store := New(ctx, "", "", "", nil)

	user := entity.User{
		Username: "testuser",
		Email:    "testuser@example.com",
		Password: "password",
		CreateAt: time.Now(),
	}

	createdUser, err := store.CreateUser(ctx, user)
	assert.Nil(t, err, "Error should be nil")
	assert.Equal(t, user.Username, createdUser.Username, "Username should match")
	assert.Equal(t, user.Email, createdUser.Email, "Email should match")

	retrievedUser, err := store.User(ctx, createdUser.Id)
	assert.Nil(t, err, "Error should be nil")
	assert.Equal(t, createdUser.Username, retrievedUser.Username, "Username should match")
	assert.Equal(t, createdUser.Email, retrievedUser.Email, "Email should match")

	updateData := entity.User{
		Username: "updateduser",
		Email:    "updateduser@example.com",
	}

	updatedUser, err := store.UpdateUser(ctx, createdUser.Id, updateData)
	assert.Nil(t, err, "Error should be nil")
	assert.Equal(t, updateData.Username, updatedUser.Username, "Username should match")
	assert.Equal(t, updateData.Email, updatedUser.Email, "Email should match")
}

func TestRoomOperations(t *testing.T) {
	store := New(ctx, "", "", "", nil)

	room := entity.Room{
		Name:        "testroom_random",
		Description: "This is a test room",
		CreateAt:    time.Now(),
	}

	createdRoom, err := store.CreateRoom(ctx, room)
	assert.Nil(t, err, "Error should be nil")
	assert.Equal(t, room.Name, createdRoom.Name, "Name should match")
	assert.Equal(t, room.Description, createdRoom.Description, "Description should match")

	rooms, err := store.Rooms(ctx, 10)
	assert.Nil(t, err, "Error should be nil")
	assert.NotEmpty(t, rooms, "Rooms should not be empty")
	assert.Equal(t, createdRoom.Name, rooms[0].Name, "Room name should match")
	assert.Equal(t, createdRoom.Description, rooms[0].Description, "Room description should match")
}

func TestMessageOperations(t *testing.T) {
	store := New(ctx, "", "", "", nil)

	user := entity.User{
		Username: "testuser_random",
		Email:    "testuser@example.com",
		Password: "password",
		CreateAt: time.Now(),
	}

	createdUser, err := store.CreateUser(ctx, user)
	assert.Nil(t, err, "Error should be nil")

	room := entity.Room{
		Name:        "testroom_room_random12",
		Description: "This is a test room",
		CreateAt:    time.Now(),
	}

	createdRoom, err := store.CreateRoom(ctx, room)
	assert.Nil(t, err, "Error should be nil")

	message := entity.Message{
		Content:   "Hello, World!",
		SenderId:  createdUser.Id,
		RoomId:    createdRoom.Id,
		Timestamp: time.Now(),
	}

	createdMessage, err := store.CreateMessage(ctx, message)
	assert.Nil(t, err, "Error should be nil")
	assert.Equal(t, message.Content, createdMessage.Content, "Content should match")
	assert.Equal(t, message.SenderId, createdMessage.SenderId, "SenderId should match")
	assert.Equal(t, message.RoomId, createdMessage.RoomId, "RoomId should match")

	// Test Messages retrieval
	messages, err := store.Messages(ctx, createdRoom.Id, 10)
	assert.Nil(t, err, "Error should be nil")
	assert.NotEmpty(t, messages, "Messages should not be empty")
	assert.Equal(t, createdMessage.Content, messages[0].Content, "Content should match")
	assert.Equal(t, createdMessage.SenderId, messages[0].SenderId, "SenderId should match")
	assert.Equal(t, createdMessage.RoomId, messages[0].RoomId, "RoomId should match")
}

func TestNotificationOperations(t *testing.T) {
	store := New(ctx, "", "", "", nil)

	user := entity.User{
		Username: "random",
		Email:    "random@example.com",
		Password: "password",
		CreateAt: time.Now(),
	}

	createdUser, err := store.CreateUser(ctx, user)
	if err != nil {
		t.Fatal(err)
	}

	room := entity.Room{
		Name:        "random_room",
		Description: "This is a test room",
		CreateAt:    time.Now(),
	}

	createdRoom, err := store.CreateRoom(ctx, room)
	if err != nil {
		t.Fatal(err)
	}

	notificationRoom := entity.NotificationRoom{
		UserId: createdUser.Id,
		RoomId: createdRoom.Id,
	}

	createdNotificationRoom, err := store.CreateNotificationRoom(ctx, notificationRoom)
	assert.Nil(t, err, "Error should be nil")

	notification := entity.Notification{
		SenderId:           createdUser.Id,
		RoomId:             createdRoom.Id,
		NotificationRoomId: createdNotificationRoom.Id,
		ReadNotification:   false,
		CreateAt:           time.Now(),
	}

	createdNotification, err := store.CreateNotification(ctx, notification)
	assert.Nil(t, err, "Error should be nil")
	assert.Equal(t, notification.SenderId, createdNotification.SenderId, "SenderId should match")
	assert.Equal(t, notification.RoomId, createdNotification.RoomId, "RoomId should match")

	notifications, err := store.Notifications(ctx, 10)
	assert.Nil(t, err, "Error should be nil")
	assert.NotEmpty(t, notifications, "Notifications should not be empty")
	assert.Equal(t, createdNotification.SenderId, notifications[0].SenderId, "SenderId should match")
	assert.Equal(t, createdNotification.RoomId, notifications[0].RoomId, "RoomId should match")

	notificationRoom = entity.NotificationRoom{
		Id:     createdNotificationRoom.Id,
		Enable: true,
	}

	updatedNotificationRoom, err := store.UpdateNotificationRoom(ctx, notificationRoom)
	assert.Nil(t, err, "Error should be nil")
	assert.Equal(t, notificationRoom.Enable, updatedNotificationRoom.Enable, "Enable should match")
}
