package database

import (
	"context"
	"sync"

	"github.com/NoobforAl/real_time_chat_application/src/config"
	"github.com/NoobforAl/real_time_chat_application/src/contract"
	"github.com/hibiken/asynq"
	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type store struct {
	db                  *mongo.Database
	cache               *redis.Client
	authClient          *authService
	messageBrokerClient *asynq.Client
	log                 contract.Logger
}

type Opts struct {
	NeedRedis       bool
	NeedMongodb     bool
	NeedAuthService bool

	NeedBrokerRoom         bool
	NeedBrokerMessage      bool
	NeedBrokerNotification bool
}

var onc sync.Once
var localStore *store
var DefaultOpts = Opts{
	NeedRedis:       true,
	NeedMongodb:     true,
	NeedAuthService: true,
}

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

func createNewClientMessageBroker(log contract.Logger) *asynq.Client {
	log.Debug("init redis client")
	return asynq.NewClient(asynq.RedisClientOpt{
		Addr:     config.RedisUri(),
		Password: config.RedisPassword(),
	})
}

func New(
	ctx context.Context,
	logger contract.Logger,
	opts ...Opts,
) contract.Store {
	onc.Do(func() {
		var (
			mongoUri      = config.MongodbUri()
			redisUri      = config.RedisUri()
			redisPassword = config.RedisPassword()
			grpcAddr      = config.GrpcAuthUri()

			messageBrokerClient *asynq.Client
			mongoDatabase       *mongo.Database
			redisClient         *redis.Client
			authClient          *authService
		)

		opt := Opts{}
		if opts == nil {
			opt = DefaultOpts
		} else {
			opt = opts[0]
		}

		if opt.NeedMongodb {
			serverAPI := options.ServerAPI(options.ServerAPIVersion1)
			optsMongodb := options.Client().ApplyURI(mongoUri).SetServerAPIOptions(serverAPI)

			mongodbClient, err := mongo.Connect(ctx, optsMongodb)
			if err != nil {
				logger.Fatal(err)
			}
			mongoDatabase = mongodbClient.Database("real_time_chat_app")
			initMongodbFiled(ctx, mongoDatabase, logger)
		}

		if opt.NeedRedis {
			redisClient = redis.NewClient(&redis.Options{
				Addr:     redisUri,
				Password: redisPassword,
				DB:       0,
			})

			statusRedis := redisClient.Ping(ctx)
			if statusRedis.Err() != nil {
				logger.Fatal(statusRedis.Err().Error())
			}

		}

		if opt.NeedAuthService {
			authClient = newAuthService(grpcAddr, logger)
		}

		if opt.NeedBrokerRoom || opt.NeedBrokerMessage || opt.NeedBrokerNotification {
			messageBrokerClient = createNewClientMessageBroker(logger)
		}

		localStore = &store{
			db:                  mongoDatabase,
			cache:               redisClient,
			authClient:          authClient,
			messageBrokerClient: messageBrokerClient,
			log:                 logger,
		}
	})

	return localStore
}
