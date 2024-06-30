package database

import (
	"context"
	"sync"

	"github.com/NoobforAl/real_time_chat_application/src/contract"
	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type store struct {
	db    *mongo.Database
	cache *redis.Client
	log   contract.Logger
}

var onc sync.Once
var localStore *store

func initMongodbFiled(
	ctx context.Context,
	client *mongo.Database,
	log contract.Logger,
) {
	// init user coll
	userColl := client.Collection("user")

	indexUserModelFiled := mongo.IndexModel{
		Keys: bson.D{
			{Key: "username", Value: 1},
			{Key: "email", Value: 1},
		},
		Options: options.Index().SetUnique(true),
	}

	_, err := userColl.Indexes().CreateOne(ctx, indexUserModelFiled)
	if err != nil {
		log.Fatal(err)
	}

	roomColl := client.Collection("room")

	indexRoomModelFiled := mongo.IndexModel{
		Keys: bson.D{
			{Key: "name", Value: 1},
		},
		Options: options.Index().SetUnique(true),
	}

	_, err = roomColl.Indexes().CreateOne(ctx, indexRoomModelFiled)
	if err != nil {
		log.Fatal(err)
	}
}

func New(
	ctx context.Context,
	mongoUri, redisUri, redisPassword string,
	logger contract.Logger,
) contract.Store {
	onc.Do(func() {
		serverAPI := options.ServerAPI(options.ServerAPIVersion1)
		optsMongodb := options.Client().ApplyURI(mongoUri).SetServerAPIOptions(serverAPI)

		mongodbClient, err := mongo.Connect(ctx, optsMongodb)
		if err != nil {
			logger.Fatal(err)
		}

		initMongodbFiled(ctx, mongodbClient.Database("real_time_chat_app"), logger)

		redisClient := redis.NewClient(&redis.Options{
			Addr:     redisUri,
			Password: redisPassword,
			DB:       0,
		})

		statusRedis := redisClient.Ping(ctx)
		if statusRedis.Err() != nil {
			logger.Fatal(statusRedis.Err().Error())
		}

		localStore = &store{
			db:    mongodbClient.Database("real_time_chat_app"),
			cache: redisClient,
			log:   logger,
		}
	})

	return localStore
}
