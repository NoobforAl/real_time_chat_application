package config

import (
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/NoobforAl/real_time_chat_application/src/contract"
	"github.com/joho/godotenv"
)

type baseConfig struct {
	redisUri      string
	redisPassword string

	mongodbUri string

	authServiceUri         string
	roomsServiceUri        string
	messageServiceURi      string
	notificationServiceUri string

	grpcAuthUri string

	secretKey   string
	maxAgeToken time.Duration

	noncForHashPassword string
}

var onc sync.Once
var config *baseConfig

func checkIsExistEnvVariables(log contract.Logger) {
	allEnvKeyNeedHaveValue := [...]string{
		"REDIS_URI",
		"MONGODB_URI",
		"SECRET_KEY",
		"GRPC_AUTH_URI",
	}

	for _, envKey := range allEnvKeyNeedHaveValue {
		value := os.Getenv(envKey)
		if value == "" {
			log.Fatalf("%s must be have a value!", envKey)
		}
	}
}

func getMaxAgeToken(log contract.Logger) time.Duration {
	maxAgeTokenStr := os.Getenv("MAX_AGE_TOKEN")
	if maxAgeTokenStr == "" {
		maxAgeTokenStr = "10"
	}

	maxAgeTokenSecund, err := strconv.Atoi(maxAgeTokenStr)
	if err != nil {
		log.Fatal("not valid max age token!")
	}

	return time.Duration(maxAgeTokenSecund) * time.Second
}

func InitConfig(logger contract.Logger) {
	onc.Do(func() {
		err := godotenv.Load(".env")
		if err != nil {
			logger.Warning("not found .env file!")
		}

		checkIsExistEnvVariables(logger)

		config = &baseConfig{
			redisUri:      os.Getenv("REDIS_URI"),
			redisPassword: os.Getenv("REDIS_PASSWORD"),

			mongodbUri: os.Getenv("MONGODB_URI"),

			authServiceUri:         os.Getenv("AUTH_SERVICE_URI"),
			roomsServiceUri:        os.Getenv("ROOMS_SERVICE_URI"),
			messageServiceURi:      os.Getenv("MESSAGE_SERVICE_URI"),
			notificationServiceUri: os.Getenv("NOTIFICATION_SERVICE_URI"),

			grpcAuthUri: os.Getenv("GRPC_AUTH_URI"),

			secretKey:   os.Getenv("SECRET_KEY"),
			maxAgeToken: getMaxAgeToken(logger),

			noncForHashPassword: os.Getenv("NONC_FOR_HASH_PASSWORD"),
		}
	})
}

func RedisUri() string      { return config.redisUri }
func RedisPassword() string { return config.redisPassword }

func AuthServiceUri() string         { return config.authServiceUri }
func RoomsServiceUri() string        { return config.roomsServiceUri }
func MessageServiceURi() string      { return config.messageServiceURi }
func NotificationServiceUri() string { return config.notificationServiceUri }

func MongodbUri() string { return config.mongodbUri }

func GrpcAuthUri() string { return config.grpcAuthUri }

func SecretKey() string          { return config.secretKey }
func MaxAgeToken() time.Duration { return config.maxAgeToken }

func NoncForHashPassword() string { return config.noncForHashPassword }
