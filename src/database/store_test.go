// !Warn: before run this unit test need to run database test!
// for after run a test
package database

import (
	"context"
	"testing"
	"time"

	"github.com/NoobforAl/real_time_chat_application/src/config"
	"github.com/NoobforAl/real_time_chat_application/src/entity"
	"github.com/NoobforAl/real_time_chat_application/src/logging"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var ctx = context.TODO()

func TestMain(m *testing.M) {
	logger := logging.New()

	if err := godotenv.Load("./../../.env"); err != nil {
		logger.Fatal("not found .env file!")
	}

	config.InitConfig(logger)

	// created new store one time
	_ = New(ctx, logger)

	m.Run()
}

func TestCache(t *testing.T) {
	store := New(ctx, nil)
	timeNow := time.Now()

	roomSample := &entity.Room{
		Id:          "random_test_room",
		Name:        "sample room",
		Description: "sample description",
		CreateAt:    timeNow,
	}

	messageSample := &entity.Message{
		Id:        "sample_id",
		Content:   "this is a test message",
		SenderId:  "sample sender id",
		RoomId:    "random_test_room",
		Timestamp: timeNow,
	}

	notificationSample := &entity.Notification{
		Id:               "sample_notification_id",
		Content:          "this is a test notification",
		SenderId:         messageSample.SenderId,
		RoomId:           roomSample.Id,
		CreateAt:         timeNow,
		ReadNotification: false,
	}

	notificationRoomSample := &entity.NotificationRoom{
		Id:       "sample_notification_room_id",
		UserId:   messageSample.SenderId,
		RoomId:   roomSample.Id,
		CreateAt: timeNow,
		Enable:   true,
	}

	t.Run("test set new Item in redis store", func(t *testing.T) {
		t.Run("test set new rooms", func(t *testing.T) {
			roomsSample := []*entity.Room{roomSample}
			err := store.SetRooms(ctx, roomsSample)
			assert.NoError(t, err)
		})

		t.Run("set up messages", func(t *testing.T) {
			messagesSample := []*entity.Message{messageSample}
			err := store.SetMessage(ctx, roomSample.Id, messagesSample)
			assert.NoError(t, err)
		})

		t.Run("set up notifications", func(t *testing.T) {
			notificationsSample := []*entity.Notification{notificationSample}
			err := store.SetNotifications(ctx, roomSample.Id, notificationsSample)
			assert.NoError(t, err)
		})

		t.Run("set up notifications room", func(t *testing.T) {
			notificationsRoomSample := []*entity.NotificationRoom{notificationRoomSample}
			err := store.SetNotificationsRoom(ctx, notificationRoomSample.UserId, notificationsRoomSample)
			assert.NoError(t, err)
		})
	})

	t.Run("test get all item form store", func(t *testing.T) {
		t.Run("test get all rooms", func(t *testing.T) {
			dataRooms, err := store.GetRooms(ctx)
			assert.NoError(t, err)
			require.Equal(t, len(dataRooms), 1)
			assert.True(t, roomsEqual(*roomSample, *dataRooms[0]))
		})

		t.Run("get message", func(t *testing.T) {
			dataMessages, err := store.GetMessages(ctx, roomSample.Id)
			assert.NoError(t, err)
			require.Equal(t, len(dataMessages), 1)
			assert.True(t, messagesEqual(*messageSample, *dataMessages[0]))
		})

		t.Run("get notifications", func(t *testing.T) {
			dataNotifications, err := store.GetNotifications(ctx, roomSample.Id)
			assert.NoError(t, err)
			require.Equal(t, len(dataNotifications), 1)
			assert.True(t, notificationsEqual(*notificationSample, *dataNotifications[0]))
		})

		t.Run("get notifications room", func(t *testing.T) {
			dataNotificationsRoom, err := store.GetNotificationsRoom(ctx, notificationRoomSample.UserId)
			assert.NoError(t, err)
			require.Equal(t, len(dataNotificationsRoom), 1)
			assert.True(t, notificationRoomsEqual(*notificationRoomSample, *dataNotificationsRoom[0]))
		})

		t.Run("get notification", func(t *testing.T) {
			dataNotification, err := store.GetNotification(ctx, roomSample.Id)
			require.NoError(t, err)
			assert.True(t, notificationsEqual(*notificationSample, dataNotification))
		})
	})

	t.Run("test logic queue", func(t *testing.T) {
		data := make([]*entity.Message, 22)
		for i := range make([]struct{}, 22) {
			data[i] = messageSample
		}

		err := store.SetMessage(ctx, roomSample.Id, data)
		assert.NoError(t, err)

		dataMessages, err := store.GetMessages(ctx, roomSample.Id)
		assert.NoError(t, err)
		require.Equal(t, len(dataMessages), 20)

		// clear data and new one
		data = data[:0]
		data = append(data, &entity.Message{RoomId: roomSample.Id})
		require.Equal(t, len(data), 1)

		err = store.SetMessage(ctx, roomSample.Id, data)
		assert.NoError(t, err)

		dataMessages, err = store.GetMessages(ctx, roomSample.Id)
		assert.NoError(t, err)
		require.Equal(t, len(dataMessages), 20)
		assert.True(t, messagesEqual(*data[0], *dataMessages[0]))
	})
}

func TestUserOperations(t *testing.T) {
	var retrievedUser entity.User
	var createdUser entity.User
	var err error

	store := New(ctx, nil)

	t.Run("test create user", func(t *testing.T) {
		user := entity.User{
			Username: "testuser",
			Email:    "testuser@example.com",
			Password: "password",
			CreateAt: time.Now(),
		}

		createdUser, err = store.CreateUser(ctx, user)
		assert.Nil(t, err, "Error should be nil")
		assert.Equal(t, user.Username, createdUser.Username, "Username should match")
		assert.Equal(t, user.Email, createdUser.Email, "Email should match")
	})

	t.Run("test get user", func(t *testing.T) {
		retrievedUser, err = store.User(ctx, createdUser.Username)
		assert.Nil(t, err, "Error should be nil")
		assert.Equal(t, createdUser.Username, retrievedUser.Username, "Username should match")
		assert.Equal(t, createdUser.Email, retrievedUser.Email, "Email should match")
	})

	t.Run("test update user", func(t *testing.T) {
		updateData := entity.User{
			Username: "updateduser",
			Email:    "updateduser@example.com",
		}

		updatedUser, err := store.UpdateUser(ctx, createdUser.Id, updateData)
		assert.Nil(t, err, "Error should be nil")
		assert.Equal(t, updateData.Username, updatedUser.Username, "Username should match")
		assert.Equal(t, updateData.Email, updatedUser.Email, "Email should match")
	})
}

func TestRoomOperations(t *testing.T) {
	var createdRoom entity.Room
	var rooms []*entity.Room
	var err error

	store := New(ctx, nil)

	t.Run("test create new room", func(t *testing.T) {
		room := entity.Room{
			Name:        "testroom_random",
			Description: "This is a test room",
			CreateAt:    time.Now(),
		}

		createdRoom, err = store.CreateRoom(ctx, room)
		assert.Nil(t, err, "Error should be nil")
		assert.Equal(t, room.Name, createdRoom.Name, "Name should match")
		assert.Equal(t, room.Description, createdRoom.Description, "Description should match")
	})

	t.Run("test get all room", func(t *testing.T) {
		rooms, err = store.Rooms(ctx, 10)
		assert.Nil(t, err, "Error should be nil")
		assert.NotEmpty(t, rooms, "Rooms should not be empty")
		assert.Equal(t, createdRoom.Name, rooms[0].Name, "Room name should match")
		assert.Equal(t, createdRoom.Description, rooms[0].Description, "Room description should match")
	})
}

func TestMessageOperations(t *testing.T) {
	var createdMessage entity.Message
	var messages []*entity.Message

	var createdRoom entity.Room
	var createdUser entity.User
	var err error

	store := New(ctx, nil)

	t.Run("create user & room for test messages", func(t *testing.T) {
		user := entity.User{
			Username: "testuser_random",
			Email:    "testuser@example.com",
			Password: "password",
			CreateAt: time.Now(),
		}

		createdUser, err = store.CreateUser(ctx, user)
		assert.Nil(t, err, "Error should be nil")

		room := entity.Room{
			Name:        "testroom_room_random12",
			Description: "This is a test room",
			CreateAt:    time.Now(),
		}

		createdRoom, err = store.CreateRoom(ctx, room)
		assert.Nil(t, err, "Error should be nil")
	})

	t.Run("create new message", func(t *testing.T) {
		message := entity.Message{
			Content:   "Hello, World!",
			SenderId:  createdUser.Id,
			RoomId:    createdRoom.Id,
			Timestamp: time.Now(),
		}

		createdMessage, err = store.CreateMessage(ctx, message)
		assert.Nil(t, err, "Error should be nil")
		assert.Equal(t, message.Content, createdMessage.Content, "Content should match")
		assert.Equal(t, message.SenderId, createdMessage.SenderId, "SenderId should match")
		assert.Equal(t, message.RoomId, createdMessage.RoomId, "RoomId should match")
	})

	t.Run("get all messages", func(t *testing.T) {
		messages, err = store.Messages(ctx, createdRoom.Id, 10)
		assert.Nil(t, err, "Error should be nil")
		assert.NotEmpty(t, messages, "Messages should not be empty")
		assert.Equal(t, createdMessage.Content, messages[0].Content, "Content should match")
		assert.Equal(t, createdMessage.SenderId, messages[0].SenderId, "SenderId should match")
		assert.Equal(t, createdMessage.RoomId, messages[0].RoomId, "RoomId should match")
	})
}

func TestNotificationOperations(t *testing.T) {
	var createdNotificationRoom entity.NotificationRoom
	var createdNotification entity.Notification
	var notifications []*entity.Notification

	var createdRoom entity.Room
	var createdUser entity.User
	var err error

	store := New(ctx, nil)

	t.Run("create new user & room", func(t *testing.T) {
		user := entity.User{
			Username: "random",
			Email:    "random@example.com",
			Password: "password",
			CreateAt: time.Now(),
		}

		createdUser, err = store.CreateUser(ctx, user)
		if err != nil {
			t.Fatal(err)
		}

		room := entity.Room{
			Name:        "random_room",
			Description: "This is a test room",
			CreateAt:    time.Now(),
		}

		createdRoom, err = store.CreateRoom(ctx, room)
		if err != nil {
			t.Fatal(err)
		}
	})

	t.Run("create new notification room", func(t *testing.T) {
		notificationRoom := entity.NotificationRoom{
			UserId: createdUser.Id,
			RoomId: createdRoom.Id,
		}

		createdNotificationRoom, err = store.CreateNotificationRoom(ctx, notificationRoom)
		assert.Nil(t, err, "Error should be nil")
		assert.Equal(t, createdNotificationRoom.UserId, notificationRoom.UserId, "SenderId should match")
		assert.Equal(t, createdNotificationRoom.RoomId, notificationRoom.RoomId, "RoomId should match")
	})

	t.Run("create new notification", func(t *testing.T) {
		notification := entity.Notification{
			SenderId:           createdUser.Id,
			RoomId:             createdRoom.Id,
			NotificationRoomId: createdNotificationRoom.Id,
			ReadNotification:   false,
			CreateAt:           time.Now(),
		}

		createdNotification, err = store.CreateNotification(ctx, notification)
		assert.Nil(t, err, "Error should be nil")
		assert.Equal(t, notification.SenderId, createdNotification.SenderId, "SenderId should match")
		assert.Equal(t, notification.RoomId, createdNotification.RoomId, "RoomId should match")
	})

	t.Run("get notification", func(t *testing.T) {
		notifications, err = store.Notifications(ctx, 10)
		assert.Nil(t, err, "Error should be nil")
		assert.NotEmpty(t, notifications, "Notifications should not be empty")
		assert.Equal(t, createdNotification.SenderId, notifications[0].SenderId, "SenderId should match")
		assert.Equal(t, createdNotification.RoomId, notifications[0].RoomId, "RoomId should match")
	})

	t.Run("update notification room", func(t *testing.T) {
		notificationRoom := entity.NotificationRoom{
			Id:     createdNotificationRoom.Id,
			Enable: true,
		}

		updatedNotificationRoom, err := store.UpdateNotificationRoom(ctx, notificationRoom)
		assert.Nil(t, err, "Error should be nil")
		assert.Equal(t, notificationRoom.Enable, updatedNotificationRoom.Enable, "Enable should match")
	})
}
