package database

import (
	"context"
	"time"

	"github.com/NoobforAl/real_time_chat_application/src/entity"
	appErrors "github.com/NoobforAl/real_time_chat_application/src/errors"
	taskRoom "github.com/NoobforAl/real_time_chat_application/src/tasks/rooms/tasks_rooms"
	"github.com/hibiken/asynq"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Room struct {
	Id          primitive.ObjectID `bson:"_id,omitempty"`
	Name        string             `bson:"name"`
	Description string             `bson:"description"`

	AllowUsers []string // TODO:#feature

	UserId string // TODO:#feature, add how user create this room

	IsOpen bool // TODO: #feature

	CreateAt time.Time `json:"create_at"`
	CloseAt  time.Time // TODO: #feature
}

func entityRoomToModel(room entity.Room) Room {
	return Room{
		Name:        room.Name,
		Description: room.Description,
		CreateAt:    room.CreateAt,
	}
}

func modelRoomToEntity(room Room) entity.Room {
	return entity.Room{
		Id:          room.Id.Hex(),
		Name:        room.Name,
		Description: room.Description,
		CreateAt:    room.CreateAt,
	}
}

func (s *store) Room(ctx context.Context, roomId string) (entity.Room, error) {
	objectID, err := primitive.ObjectIDFromHex(roomId)
	if err != nil {
		s.log.Error(err)
		return entity.Room{}, appErrors.ErrNotValidId
	}

	var room entity.Room
	err = s.db.Collection("room").FindOne(ctx, bson.M{"_id": objectID}).Decode(&room)
	if err != nil {
		return entity.Room{}, err
	}

	return room, nil
}

func (s *store) Rooms(ctx context.Context, maxLen int) ([]*entity.Room, error) {
	filter := bson.M{}
	findOptions := options.Find()
	findOptions.SetLimit(10)

	cursor, err := s.db.Collection("room").Find(ctx, filter, findOptions)
	if err != nil {
		return nil, err
	}

	defer func() {
		if err = cursor.Close(ctx); err != nil {
			s.log.Error(err)
		}
	}()

	result := make([]*entity.Room, 0)
	for cursor.Next(ctx) {
		var room Room
		err := cursor.Decode(&room)
		if err != nil {
			return nil, err
		}

		entityRoom := modelRoomToEntity(room)
		result = append(result, &entityRoom)
	}

	return result, err
}

func (s *store) CreateRoom(ctx context.Context, room entity.Room) (entity.Room, error) {
	roomModel := entityRoomToModel(room)
	result, err := s.db.Collection("room").InsertOne(ctx, roomModel)
	if err == mongo.ErrNoDocuments {
		return entity.Room{}, appErrors.ErrNoDocuments
	} else if err != nil {
		s.log.Error(err)
		return entity.Room{}, err
	}

	id, _ := result.InsertedID.(primitive.ObjectID)
	roomModel.Id = id

	return modelRoomToEntity(roomModel), nil
}

func (s *store) SendNewRoom(ctx context.Context, room entity.Room) error {
	msgBroker := s.messageBrokerClient

	newTask, err := taskRoom.NewRoomSaveTask(room)
	if err != nil {
		return err
	}

	_, err = msgBroker.EnqueueContext(ctx, newTask, asynq.MaxRetry(10), asynq.Timeout(60*time.Second))
	return err
}
