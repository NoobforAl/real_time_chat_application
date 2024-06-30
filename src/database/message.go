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

type Message struct {
	Id        primitive.ObjectID `bson:"_id,omitempty"`
	Content   string             `bson:"content"`
	SenderId  primitive.ObjectID `bson:"sender_id"`
	RoomId    primitive.ObjectID `bson:"room_id"`
	Timestamp time.Time          `bson:"timestamp"`
}

func entityMessageToModel(message entity.Message) (Message, error) {
	senderId, err := primitive.ObjectIDFromHex(message.SenderId)
	if err != nil {
		return Message{}, appErrors.ErrNotValidId
	}

	roomId, err := primitive.ObjectIDFromHex(message.RoomId)
	if err != nil {
		return Message{}, appErrors.ErrNotValidId
	}

	return Message{
		Content:   message.Content,
		SenderId:  senderId,
		RoomId:    roomId,
		Timestamp: message.Timestamp,
	}, nil
}

func modelMessageToEntity(message Message) entity.Message {
	return entity.Message{
		Id:        message.Id.Hex(),
		Content:   message.Content,
		SenderId:  message.SenderId.Hex(),
		RoomId:    message.RoomId.Hex(),
		Timestamp: message.Timestamp,
	}
}

func (s *store) Messages(ctx context.Context, roomId string, maxLen int) ([]*entity.Message, error) {
	objectID, err := primitive.ObjectIDFromHex(roomId)
	if err != nil {
		s.log.Error(err)
		return nil, appErrors.ErrNotValidId
	}

	filter := bson.M{"room_id": objectID}
	findOptions := options.Find()
	findOptions.SetLimit(int64(maxLen))

	cursor, err := s.db.Collection("message").Find(ctx, filter, findOptions)
	if err != nil {
		return nil, err
	}

	defer func() {
		if err = cursor.Close(ctx); err != nil {
			s.log.Error(err)
		}
	}()

	result := make([]*entity.Message, 0)
	for cursor.Next(ctx) {
		var message Message
		err := cursor.Decode(&message)
		if err != nil {
			return nil, err
		}

		entityMessage := modelMessageToEntity(message)
		result = append(result, &entityMessage)
	}

	return result, err
}

func (s *store) CreateMessage(ctx context.Context, message entity.Message) (entity.Message, error) {
	messageModel, err := entityMessageToModel(message)
	if err != nil {
		return entity.Message{}, err
	}

	userColl := s.db.Collection("user")
	err = userColl.FindOne(ctx, bson.M{"_id": messageModel.SenderId}).Decode(&User{})
	if err == mongo.ErrNoDocuments {
		return entity.Message{}, appErrors.ErrNoDocuments
	} else if err != nil {
		s.log.Error(err)
		return entity.Message{}, err
	}

	roomColl := s.db.Collection("room")
	err = roomColl.FindOne(ctx, bson.M{"_id": messageModel.RoomId}).Decode(&Room{})
	if err == mongo.ErrNoDocuments {
		return entity.Message{}, appErrors.ErrNoDocuments
	} else if err != nil {
		s.log.Error(err)
		return entity.Message{}, err
	}

	result, err := s.db.Collection("message").InsertOne(ctx, messageModel)
	if err == mongo.ErrNoDocuments {
		return entity.Message{}, appErrors.ErrNoDocuments
	} else if err != nil {
		s.log.Error(err)
		return entity.Message{}, err
	}

	id, _ := result.InsertedID.(primitive.ObjectID)
	messageModel.Id = id

	return modelMessageToEntity(messageModel), nil
}
