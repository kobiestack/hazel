package auth

import (
	"errors"
	"log/slog"
	"os"
	"time"

	"github.com/freekobie/hazel/models"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

var (
	ErrInvalidToken = errors.New("the provided token is no valid")
)

type TokenType string
type CustomClaims struct {
	TokenType string `json:"token_type"`
	jwt.RegisteredClaims
}

const (
	TokenTypeAccess  TokenType = "ACCESS"
	TokenTypeRefresh TokenType = "REFRESH"
)

type UserSession struct {
	User         models.User `json:"user"`
	RefreshToken string      `json:"refreshToken"`
	ExpiresAt    time.Time   `json:"expiresAt"`
}

type UserAccess struct {
	AccessToken string    `json:"accessToken"`
	ExpiresAt   time.Time `json:"expiresAt"`
}

func GenerateToken(userID uuid.UUID, duration time.Duration, tokenType TokenType) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"iat":        time.Now().UTC().UnixNano(),
		"exp":        time.Now().Add(duration).UnixNano(),
		"sub":        userID.String(),
		"token_type": tokenType,
	})

	tokenString, err := token.SignedString([]byte(os.Getenv("TOKEN_SECRET")))
	if err != nil {
		slog.Error("failed to sign access token", "error", err.Error())
		return "", err
	}

	return tokenString, nil
}

func ValidateToken(tokenStr string, tokenType TokenType) (uuid.UUID, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &CustomClaims{}, func(t *jwt.Token) (any, error) {
		return []byte(os.Getenv("TOKEN_SECRET")), nil
	})
	if err != nil || !token.Valid {
		return uuid.Nil, err
	}

	if claims, ok := token.Claims.(*CustomClaims); ok {
		if claims.TokenType != string(tokenType) {
			return uuid.Nil, ErrInvalidToken
		}
	}
	userIDString, err := token.Claims.GetSubject()
	if err != nil {
		slog.Error("failed to fetch 'sub' claim", "error", err.Error())
		return uuid.Nil, err
	}

	userID := uuid.MustParse(userIDString)

	return userID, nil
}
