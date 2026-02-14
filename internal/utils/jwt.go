package utils

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/goutamkumar/golang_restapi_postgresql_test1/internal/config"
)

type JWTClaims struct {
	UserID string `json:"user_id"`
	Role   string `json:"role"` // "user"
	jwt.RegisteredClaims
}

func CreateToken(claims JWTClaims) (string, error) {
	env := config.LoadEnv()
	secret := []byte(env.JWTSecret)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString(secret)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func VerifyJWT(tokenString string) (*JWTClaims, error) {
	env := config.LoadEnv()
	secret := []byte(env.JWTSecret)

	token, err := jwt.ParseWithClaims(
		tokenString,
		&JWTClaims{},
		func(token *jwt.Token) (interface{}, error) {
			return secret, nil
		},
	)

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*JWTClaims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token")
	}

	return claims, nil
}

func PasswordResetToken(claims JWTClaims) (string, error) {
	env := config.LoadEnv()
	secret := []byte(env.JWTSecret)

	// Set a short expiration time for the reset token (e.g., 15 minutes)
	claims.ExpiresAt = jwt.NewNumericDate(time.Now().Add(15 * time.Minute))

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString(secret)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func VerifyPasswordResetToken(tokenString string) (*JWTClaims, error) {
	env := config.LoadEnv()
	secret := []byte(env.JWTSecret)

	token, err := jwt.ParseWithClaims(
		tokenString,
		&JWTClaims{},
		func(token *jwt.Token) (interface{}, error) {
			return secret, nil
		},
	)

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*JWTClaims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token")
	}

	return claims, nil
}
