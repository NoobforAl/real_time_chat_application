package database

import (
	"context"
	"time"

	"github.com/NoobforAl/real_time_chat_application/src/entity"
	appErrors "github.com/NoobforAl/real_time_chat_application/src/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type User struct {
	Id       primitive.ObjectID `bson:"_id,omitempty"`
	Username string             `bson:"username"`
	Email    string             `bson:"email"`
	Password string             `bson:"password"`

	Notification bool `bson:"notification"`

	CreateAt time.Time `bson:"create_at"`
}

func entityUserToModel(user entity.User) User {
	return User{
		Username: user.Username,
		Email:    user.Email,
		Password: user.Password,

		Notification: user.Notification,

		CreateAt: user.CreateAt,
	}
}

func modelUserToEntity(user User) entity.User {
	return entity.User{
		Id: user.Id.Hex(),

		Username: user.Username,
		Email:    user.Email,
		Password: user.Password,

		Notification: user.Notification,

		CreateAt: user.CreateAt,
	}
}

func (s *store) User(ctx context.Context, username string) (entity.User, error) {
	var user User
	err := s.db.Collection("user").FindOne(ctx, bson.M{"username": username}).Decode(&user)
	if err == mongo.ErrNoDocuments {
		return entity.User{}, appErrors.ErrNoDocuments
	} else if err != nil {
		s.log.Error(err)
		return entity.User{}, err
	}

	return modelUserToEntity(user), nil
}

func (s *store) CreateUser(ctx context.Context, userData entity.User) (entity.User, error) {
	user := entityUserToModel(userData)

	result, err := s.db.Collection("user").InsertOne(ctx, user)
	if err == mongo.ErrMultipleIndexDrop {
		return entity.User{}, appErrors.ErrDatabaseIndex
	} else if err != nil {
		s.log.Error(err)
		return entity.User{}, err
	}

	id, _ := result.InsertedID.(primitive.ObjectID)
	user.Id = id

	return modelUserToEntity(user), nil
}

func (s *store) UpdateUser(ctx context.Context, id string, updatedData entity.User) (entity.User, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		s.log.Error(err)
		return entity.User{}, appErrors.ErrNotValidId
	}

	updateFields := bson.M{}
	if updatedData.Username != "" {
		updateFields["username"] = updatedData.Username
	}

	if updatedData.Email != "" {
		updateFields["email"] = updatedData.Email
	}

	if updatedData.Password != "" {
		updateFields["password"] = updatedData.Password
	}

	if !updatedData.CreateAt.IsZero() {
		updateFields["create_at"] = updatedData.CreateAt
	}

	updateFields["notification"] = updatedData.Notification

	if len(updateFields) == 0 {
		return entity.User{}, appErrors.ErrNoFieldsToUpdate
	}

	update := bson.M{
		"$set": updateFields,
	}

	result, err := s.db.Collection("user").UpdateOne(ctx, bson.M{"_id": objectID}, update)
	if err == mongo.ErrNoDocuments {
		return entity.User{}, appErrors.ErrNoDocuments
	} else if err != nil {
		s.log.Error(err)
		return entity.User{}, err
	}

	if result.MatchedCount == 0 {
		return entity.User{}, appErrors.ErrNoDocuments
	}

	var updatedUser User
	err = s.db.Collection("user").FindOne(ctx, bson.M{"_id": objectID}).Decode(&updatedUser)
	if err != nil {
		s.log.Error(err)
		return entity.User{}, err
	}

	return modelUserToEntity(updatedUser), nil
}
