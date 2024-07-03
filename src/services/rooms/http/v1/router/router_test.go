package router

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/NoobforAl/real_time_chat_application/src/config"
	"github.com/NoobforAl/real_time_chat_application/src/database"
	"github.com/NoobforAl/real_time_chat_application/src/entity"
	"github.com/NoobforAl/real_time_chat_application/src/grpc/auth"
	"github.com/NoobforAl/real_time_chat_application/src/logging"
	"github.com/NoobforAl/real_time_chat_application/src/services/auth/jwt"
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
)

var app *fiber.App

func TestMain(m *testing.M) {
	err := godotenv.Load("./../../../../../../.env")
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	app = fiber.New()
	logger := logging.New()
	config.InitConfig(logger)

	store := database.New(ctx, logger)

	SetupRoomRoute(app, store, logger)

	authService := auth.New(store, logger)
	go authService.Run(config.GrpcAuthUri())

	m.Run()
}

func TestRoutApp(t *testing.T) {
	secretKey := config.SecretKey()

	sampleUser := entity.User{
		Id:       "sample-test-id",
		Username: "test-user-name",
	}

	accessToken, _, err := jwt.GenerateTokens([]byte(secretKey), sampleUser.Id, sampleUser.Username, time.Hour, time.Hour)
	if err != nil {
		t.Fatal(err)
	}

	payload := map[string]any{
		"name":        "new_user_room_test",
		"description": "description_room_test",
	}

	t.Run("send broken token for api", func(t *testing.T) {
		body, _ := json.Marshal(payload)
		req := httptest.NewRequest("GET", "/rooms/", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Add("Access-Token", accessToken+"bad string, this make it pass test:)")

		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	})

	t.Run("test create new room", func(t *testing.T) {
		body, _ := json.Marshal(payload)
		req := httptest.NewRequest("POST", "/rooms/", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Add("Access-Token", accessToken)

		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})

	t.Run("test create new room again with same name", func(t *testing.T) {
		body, _ := json.Marshal(payload)
		req := httptest.NewRequest("POST", "/rooms/", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Add("Access-Token", accessToken)

		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("test get all rooms", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/rooms/", nil)
		req.Header.Add("Access-Token", accessToken)

		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		defer resp.Body.Close()
		bodyBytes, err := io.ReadAll(resp.Body)
		assert.NoError(t, err)

		var rooms []entity.Room
		err = json.Unmarshal(bodyBytes, &rooms)
		assert.NoError(t, err)

		assert.NotEmpty(t, rooms)
		for _, room := range rooms {
			assert.NotEmpty(t, room.Name)
			assert.NotEmpty(t, room.Description)

			assert.Equal(t, room.Name, payload["name"])
			assert.Equal(t, room.Description, payload["description"])
		}
	})
}
