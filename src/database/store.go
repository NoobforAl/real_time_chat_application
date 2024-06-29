package database

import (
	"context"
	"sync"

	"github.com/NoobforAl/real_time_chat_application/src/contract"
	"github.com/NoobforAl/real_time_chat_application/src/entity"
	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type store struct {
	db    *mongo.Client
	cache *redis.Client
	log   contract.Logger
}

var onc sync.Once
var localStore *store

func New(ctx context.Context, mongoUri, redisUri, redisPassword string, logger contract.Logger) contract.Store {
	onc.Do(func() {
		serverAPI := options.ServerAPI(options.ServerAPIVersion1)
		optsMongodb := options.Client().ApplyURI(mongoUri).SetServerAPIOptions(serverAPI)

		mongodbClient, err := mongo.Connect(ctx, optsMongodb)
		if err != nil {
			logger.Fatal(err)
		}

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
			db:    mongodbClient,
			cache: redisClient,
			log:   logger,
		}
	})

	return localStore
}

func (s *store) User(ctx context.Context, id string) entity.User
func (s *store) CreateUser(ctx context.Context, userData entity.User) entity.User

func (s *store) Messages(ctx context.Context, roomId string, maxLen int) []*entity.Message
func (s *store) CreateMessage(ctx context.Context, message entity.Message) entity.Message

func (s *store) Rooms(ctx context.Context, maxLen int) []*entity.Room
func (s *store) CreateRoom(ctx context.Context, message entity.Room) entity.Room
