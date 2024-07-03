package database

import (
	"context"
	"time"

	"github.com/NoobforAl/real_time_chat_application/src/entity"
	appErrors "github.com/NoobforAl/real_time_chat_application/src/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Notification struct {
	Id                 primitive.ObjectID `bson:"_id,omitempty"`
	SenderId           primitive.ObjectID `bson:"sender_id,omitempty"`
	RoomId             primitive.ObjectID `bson:"room_id,omitempty"`
	NotificationRoomId primitive.ObjectID `bson:"notification_room_id,omitempty"`
	Content            string             // TODO
	CreateAt           time.Time          `bson:"create_at"`
	ReadNotification   bool               `bson:"read_notification"`
}

type NotificationRoom struct {
	Id       primitive.ObjectID `bson:"_id,omitempty"`
	UserId   primitive.ObjectID `bson:"user_id,omitempty"`
	RoomId   primitive.ObjectID `bson:"room_id,omitempty"`
	CreateAt time.Time          `bson:"create_at"`
	Enable   bool               `bson:"enable"`
}

func entityNotificationToModel(notification entity.Notification) (Notification, error) {
	senderId, err := primitive.ObjectIDFromHex(notification.SenderId)
	if err != nil {
		return Notification{}, appErrors.ErrNotValidId
	}

	roomId, err := primitive.ObjectIDFromHex(notification.RoomId)
	if err != nil {
		return Notification{}, appErrors.ErrNotValidId
	}

	notificationRoomId, err := primitive.ObjectIDFromHex(notification.NotificationRoomId)
	if err != nil {
		return Notification{}, appErrors.ErrNotValidId
	}

	return Notification{
		SenderId:           senderId,
		RoomId:             roomId,
		NotificationRoomId: notificationRoomId,
		ReadNotification:   notification.ReadNotification,
		CreateAt:           notification.CreateAt,
	}, nil
}

func modelNotificationToEntity(notification Notification) entity.Notification {
	return entity.Notification{
		Id:                 notification.Id.Hex(),
		SenderId:           notification.SenderId.Hex(),
		RoomId:             notification.RoomId.Hex(),
		NotificationRoomId: notification.NotificationRoomId.Hex(),
		ReadNotification:   notification.ReadNotification,
		CreateAt:           notification.CreateAt,
	}
}

func entityNotificationRoomToModel(notificationRoom entity.NotificationRoom) (NotificationRoom, error) {
	userId, err := primitive.ObjectIDFromHex(notificationRoom.UserId)
	if err != nil {
		return NotificationRoom{}, appErrors.ErrNotValidId
	}

	roomId, err := primitive.ObjectIDFromHex(notificationRoom.RoomId)
	if err != nil {
		return NotificationRoom{}, appErrors.ErrNotValidId
	}

	return NotificationRoom{
		UserId:   userId,
		RoomId:   roomId,
		Enable:   notificationRoom.Enable,
		CreateAt: notificationRoom.CreateAt,
	}, nil
}

func modelNotificationRoomToEntity(notificationRoom NotificationRoom) entity.NotificationRoom {
	return entity.NotificationRoom{
		Id:       notificationRoom.Id.Hex(),
		UserId:   notificationRoom.UserId.Hex(),
		RoomId:   notificationRoom.RoomId.Hex(),
		Enable:   notificationRoom.Enable,
		CreateAt: notificationRoom.CreateAt,
	}
}

func (s *store) Notifications(ctx context.Context, maxLen int) ([]*entity.Notification, error) {
	filter := bson.M{}
	findOptions := options.Find()
	findOptions.SetLimit(int64(maxLen))

	cursor, err := s.db.Collection("notification").Find(ctx, filter, findOptions)
	if err != nil {
		return nil, err
	}

	defer func() {
		if err = cursor.Close(ctx); err != nil {
			s.log.Error(err)
		}
	}()

	result := make([]*entity.Notification, 0)
	for cursor.Next(ctx) {
		var notification Notification
		err := cursor.Decode(&notification)
		if err != nil {
			return nil, err
		}

		entityNotification := modelNotificationToEntity(notification)
		result = append(result, &entityNotification)
	}

	return result, err
}

func (s *store) CreateNotification(ctx context.Context, notificationData entity.Notification) (entity.Notification, error) {
	notification, err := entityNotificationToModel(notificationData)
	if err != nil {
		return entity.Notification{}, err
	}

	result, err := s.db.Collection("notification").InsertOne(ctx, notification)
	if err != nil {
		s.log.Error(err)
		return entity.Notification{}, err
	}

	id, _ := result.InsertedID.(primitive.ObjectID)
	notification.Id = id

	return modelNotificationToEntity(notification), nil
}

func (s *store) CreateNotificationRoom(ctx context.Context, notificationRoomData entity.NotificationRoom) (entity.NotificationRoom, error) {
	notificationRoomDataModel, err := entityNotificationRoomToModel(notificationRoomData)
	if err == nil {
		return entity.NotificationRoom{}, err
	}

	existingNotificationRoom := NotificationRoom{}
	err = s.db.Collection("notification_room").FindOne(ctx, bson.M{"user_id": notificationRoomDataModel.UserId, "room_id": notificationRoomDataModel.RoomId}).Decode(&existingNotificationRoom)
	if err == nil {
		return entity.NotificationRoom{}, appErrors.ErrDatabaseIndex
	} else if err != mongo.ErrNoDocuments {
		s.log.Error(err)
		return entity.NotificationRoom{}, err
	}

	notificationRoom := NotificationRoom{
		UserId:   notificationRoomDataModel.UserId,
		RoomId:   notificationRoomDataModel.RoomId,
		Enable:   notificationRoomData.Enable,
		CreateAt: notificationRoomData.CreateAt,
	}

	result, err := s.db.Collection("notification_room").InsertOne(ctx, notificationRoom)
	if err != nil {
		s.log.Error(err)
		return entity.NotificationRoom{}, err
	}

	id, _ := result.InsertedID.(primitive.ObjectID)
	notificationRoom.Id = id

	return modelNotificationRoomToEntity(notificationRoom), nil
}

func (s *store) UpdateNotificationRoom(ctx context.Context, notificationRoomData entity.NotificationRoom) (entity.NotificationRoom, error) {
	objectID, err := primitive.ObjectIDFromHex(notificationRoomData.Id)
	if err != nil {
		s.log.Error(err)
		return entity.NotificationRoom{}, appErrors.ErrNotValidId
	}

	updateFields := bson.M{}
	if notificationRoomData.UserId != "" {
		updateFields["user_id"] = notificationRoomData.UserId
	}

	if notificationRoomData.RoomId != "" {
		updateFields["room_id"] = notificationRoomData.RoomId
	}

	updateFields["enable"] = notificationRoomData.Enable

	if !notificationRoomData.CreateAt.IsZero() {
		updateFields["create_at"] = notificationRoomData.CreateAt
	}

	if len(updateFields) == 0 {
		return entity.NotificationRoom{}, appErrors.ErrNoFieldsToUpdate
	}

	update := bson.M{
		"$set": updateFields,
	}

	_, err = s.db.Collection("notification_room").UpdateOne(ctx, bson.M{"_id": objectID}, update)
	if err == mongo.ErrNoDocuments {
		return entity.NotificationRoom{}, appErrors.ErrNoDocuments
	} else if err != nil {
		s.log.Error(err)
		return entity.NotificationRoom{}, err
	}

	var updatedNotificationRoom NotificationRoom
	err = s.db.Collection("notification_room").FindOne(ctx, bson.M{"_id": objectID}).Decode(&updatedNotificationRoom)
	if err == mongo.ErrNoDocuments {
		return entity.NotificationRoom{}, appErrors.ErrNoDocuments
	} else if err != nil {
		s.log.Error(err)
		return entity.NotificationRoom{}, err
	}

	return modelNotificationRoomToEntity(updatedNotificationRoom), nil
}

// TODO: feature
func (s *store) NotificationRoom(ctx context.Context, notificationRoomData entity.NotificationRoom) (entity.NotificationRoom, error) {
	return entity.NotificationRoom{}, nil
}

func (s *store) SendNewNotification(ctx context.Context, message entity.Notification) error {
	return nil
}
