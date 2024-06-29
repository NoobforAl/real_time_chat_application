package jwt

import (
	"testing"
	"time"
)

func TestJwtToken(t *testing.T) {
	secretKey := []byte("random_secretKey:) jdlafjirjo2efs")
	id := "sdfaaefdmsaeilf1"
	username := "Abas"
	accessTokenDuration := 500 * time.Millisecond
	refreshTokenDuration := 1500 * time.Millisecond

	var accessToken, refreshToken string

	t.Run("Test Create New Token", func(t *testing.T) {
		var err error
		accessToken, refreshToken, err = GenerateTokens(secretKey, id, username, accessTokenDuration, refreshTokenDuration)
		if err != nil {
			t.Fatalf("Failed to generate tokens: %v", err)
		}

		if accessToken == "" || refreshToken == "" {
			t.Fatal("Generated tokens should not be empty")
		}
	})

	t.Run("Test Validate Access Token", func(t *testing.T) {
		time.Sleep(500 * time.Millisecond)

		claims, err := ValidateToken(accessToken, secretKey)
		if err != nil {
			t.Fatalf("Failed to validate access token: %v", err)
		}

		if claims.Username != username {
			t.Errorf("Expected username %v, got %v", username, claims.Username)
		}
	})

	t.Run("Test Expired Access Token", func(t *testing.T) {
		time.Sleep(700 * time.Millisecond)

		_, err := ValidateToken(accessToken, secretKey)
		if err == nil {
			t.Fatal("Expected error for expired token, but got none")
		}
	})

	t.Run("Test Refresh Access Token", func(t *testing.T) {
		newAccessToken, err := RefreshAccessToken(refreshToken, secretKey, accessTokenDuration)
		if err != nil {
			t.Fatalf("Failed to refresh access token: %v", err)
		}

		if newAccessToken == "" {
			t.Fatal("New access token should not be empty")
		}

		claims, err := ValidateToken(newAccessToken, secretKey)
		if err != nil {
			t.Fatalf("Failed to validate new access token: %v", err)
		}

		if claims.Username != username {
			t.Errorf("Expected username %v, got %v", username, claims.Username)
		}
	})

	t.Run("Test Expired Refresh Token", func(t *testing.T) {
		time.Sleep(refreshTokenDuration)

		_, err := RefreshAccessToken(refreshToken, secretKey, accessTokenDuration)
		if err == nil {
			t.Fatal("Expected error for expired refresh token, but got none")
		}
	})
}
