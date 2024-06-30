package router

import (
	"bytes"
	"context"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/NoobforAl/real_time_chat_application/src/config"
	"github.com/NoobforAl/real_time_chat_application/src/database"
	"github.com/NoobforAl/real_time_chat_application/src/logging"
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

	mongodbUri := config.MongodbUri()
	redisUri := config.RedisUri()
	redisPassword := config.RedisPassword()
	store := database.New(ctx, mongodbUri, redisUri, redisPassword, logger)

	SetupAuthRoute(ctx, app, store, logger)

	m.Run()
}

func TestAuthApp(t *testing.T) {
	t.Run("Register", func(t *testing.T) {
		payload := map[string]string{
			"username": "new_user_auth_test",
			"password": "password_auth_test",
			"email":    "test@test.com",
		}
		body, _ := json.Marshal(payload)
		req := httptest.NewRequest("POST", "/register/", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")

		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})

	t.Run("Login - Successful", func(t *testing.T) {
		payload := map[string]string{
			"username": "new_user_auth_test",
			"password": "password_auth_test",
		}
		body, _ := json.Marshal(payload)
		req := httptest.NewRequest("POST", "/login/", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")

		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})

	t.Run("Login - Invalid Password", func(t *testing.T) {
		payload := map[string]string{
			"username": "new_user_auth_test",
			"password": "wrong_password",
		}
		body, _ := json.Marshal(payload)
		req := httptest.NewRequest("POST", "/login/", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")

		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	})

	t.Run("Login - User Not Found", func(t *testing.T) {
		payload := map[string]string{
			"username": "non_existing_user",
			"password": "password",
		}
		body, _ := json.Marshal(payload)
		req := httptest.NewRequest("POST", "/login/", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")

		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	})
}
