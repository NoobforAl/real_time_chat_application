package database

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/NoobforAl/real_time_chat_application/src/entity"
	"github.com/redis/go-redis/v9"
)

const (
	roomKeyPrefix             = "rooms="
	messagesKeyPrefix         = "messages-room-id="
	notificationKeyPrefix     = "Notification-room-id="
	notificationRoomKeyPrefix = "Notification-room-user-id="

	dequeSize = 20
)

func setDataInRedisQueue(ctx context.Context, rdb *redis.Client, key string, data ...string) error {
	err := rdb.LPush(ctx, key, data).Err()
	if err != nil {
		return err
	}

	err = rdb.LTrim(ctx, key, 0, int64(dequeSize-1)).Err()
	if err != nil {
		return err
	}
	return nil
}

func getDataInRedisQueue(ctx context.Context, rdb *redis.Client, key string) ([]string, error) {
	return rdb.LRange(ctx, key, 0, -1).Result()
}

func decodeDataSliceToByte(data []any) ([]string, error) {
	var dataEncoded []string
	for _, v := range data {
		room, err := json.Marshal(v)
		if err != nil {
			return nil, err
		}

		dataEncoded = append(dataEncoded, string(room))
	}

	return dataEncoded, nil
}

func (s *store) GetRooms(ctx context.Context) ([]*entity.Room, error) {
	data, err := getDataInRedisQueue(ctx, s.cache, roomKeyPrefix)
	if err != nil {
		return nil, err
	}

	var result []*entity.Room
	for _, v := range data {
		room := &entity.Room{}
		err = json.Unmarshal([]byte(v), room)
		if err != nil {
			return nil, err
		}
		result = append(result, room)
	}

	return result, nil
}

func (s *store) SetRooms(ctx context.Context, data []*entity.Room) error {
	rooms := make([]any, len(data))
	for i, v := range data {
		rooms[i] = v
	}

	dataEncoded, err := decodeDataSliceToByte(rooms)
	if err != nil {
		return err
	}

	return setDataInRedisQueue(ctx, s.cache, roomKeyPrefix, dataEncoded...)
}

func (s *store) GetMessages(ctx context.Context, room_id string) ([]*entity.Message, error) {
	key := messagesKeyPrefix + room_id
	data, err := getDataInRedisQueue(ctx, s.cache, key)
	if err != nil {
		return nil, err
	}

	var result []*entity.Message
	for _, v := range data {
		message := &entity.Message{}
		err = json.Unmarshal([]byte(v), message)
		if err != nil {
			return nil, err
		}
		result = append(result, message)
	}

	return result, nil
}

func (s *store) SetMessage(ctx context.Context, room_id string, data []*entity.Message) error {
	key := messagesKeyPrefix + room_id
	messages := make([]any, len(data))
	for i, v := range data {
		messages[i] = v
	}

	dataEncoded, err := decodeDataSliceToByte(messages)
	if err != nil {
		return err
	}

	return setDataInRedisQueue(ctx, s.cache, key, dataEncoded...)
}

func (s *store) GetNotification(ctx context.Context, room_id string) (entity.Notification, error) {
	key := notificationKeyPrefix + room_id
	notif, err := s.cache.BLPop(ctx, 10*time.Second, key).Result()
	if err != nil {
		return entity.Notification{}, err
	}

	if len(notif) != 2 {
		return entity.Notification{}, fmt.Errorf("bad data from queue")
	}

	notification := entity.Notification{}
	err = json.Unmarshal([]byte(notif[1]), &notification)
	if err != nil {
		return entity.Notification{}, err
	}

	return notification, nil
}

func (s *store) GetNotifications(ctx context.Context, room_id string) ([]*entity.Notification, error) {
	key := notificationKeyPrefix + room_id
	data, err := getDataInRedisQueue(ctx, s.cache, key)
	if err != nil {
		return nil, err
	}

	var result []*entity.Notification
	for _, v := range data {
		notification := &entity.Notification{}
		err = json.Unmarshal([]byte(v), notification)
		if err != nil {
			return nil, err
		}
		result = append(result, notification)
	}

	return result, nil
}

func (s *store) SetNotifications(ctx context.Context, room_id string, data []*entity.Notification) error {
	key := notificationKeyPrefix + room_id
	notifications := make([]any, len(data))
	for i, v := range data {
		notifications[i] = v
	}

	dataEncoded, err := decodeDataSliceToByte(notifications)
	if err != nil {
		return err
	}

	return setDataInRedisQueue(ctx, s.cache, key, dataEncoded...)
}

func (s *store) GetNotificationsRoom(ctx context.Context, user_id string) ([]*entity.NotificationRoom, error) {
	key := notificationRoomKeyPrefix + user_id
	data, err := getDataInRedisQueue(ctx, s.cache, key)
	if err != nil {
		return nil, err
	}

	var result []*entity.NotificationRoom
	for _, v := range data {
		notificationRoom := &entity.NotificationRoom{}
		err = json.Unmarshal([]byte(v), notificationRoom)
		if err != nil {
			return nil, err
		}
		result = append(result, notificationRoom)
	}

	return result, nil
}

func (s *store) SetNotificationsRoom(ctx context.Context, user_id string, data []*entity.NotificationRoom) error {
	key := notificationRoomKeyPrefix + user_id
	notifications := make([]any, len(data))
	for i, v := range data {
		notifications[i] = v
	}

	dataEncoded, err := decodeDataSliceToByte(notifications)
	if err != nil {
		return err
	}

	return setDataInRedisQueue(ctx, s.cache, key, dataEncoded...)
}
