package jwt

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt"
)

type TokenClaims struct {
	Username      string `json:"username"`
	Notifications bool   `json:"notifications"`
	jwt.StandardClaims
}

func GenerateTokens(
	secretKey []byte,
	id string,
	username string,
	accessTokenDuration time.Duration,
	refreshTokenDuration time.Duration,
) (string, string, error) {
	accessTokenClaims := TokenClaims{
		Username: username,
		StandardClaims: jwt.StandardClaims{
			Id:        id,
			ExpiresAt: time.Now().Add(accessTokenDuration).Unix(),
		},
	}
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessTokenClaims)
	accessTokenString, err := accessToken.SignedString(secretKey)
	if err != nil {
		return "", "", err
	}

	refreshTokenClaims := TokenClaims{
		Username: username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(refreshTokenDuration).Unix(),
		},
	}

	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshTokenClaims)
	refreshTokenString, err := refreshToken.SignedString(secretKey)
	if err != nil {
		return "", "", err
	}

	return accessTokenString, refreshTokenString, nil
}

func ValidateToken(tokenString string, secretKey []byte) (*TokenClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &TokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		return secretKey, nil
	})
	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*TokenClaims); ok && token.Valid {
		if claims.ExpiresAt < time.Now().Unix() {
			return nil, errors.New("token is expired")
		}
		return claims, nil
	}
	return nil, errors.New("invalid token")
}

func RefreshAccessToken(
	refreshTokenString string,
	secretKey []byte,
	newAccessTokenDuration time.Duration,
) (string, error) {
	claims, err := ValidateToken(refreshTokenString, secretKey)
	if err != nil {
		return "", errors.New("invalid or expired refresh token")
	}

	accessTokenClaims := TokenClaims{
		Username: claims.Username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(newAccessTokenDuration).Unix(),
		},
	}

	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessTokenClaims)
	accessTokenString, err := accessToken.SignedString(secretKey)
	if err != nil {
		return "", err
	}

	return accessTokenString, nil
}
