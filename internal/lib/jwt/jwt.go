package myjwt

import (
	"fmt"
	"log/slog"
	"time"

	"github.com/Weit145/Auth_golang/internal/config"
	"github.com/Weit145/Auth_golang/internal/lib/logger"
	"github.com/golang-jwt/jwt/v5"
)

func CreateEmailJWT(cfg *config.Config, log *slog.Logger, email string) (string, error) {
	const op = "jwt.CreateEmailJWT"

	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["email"] = email
	claims["exp"] = time.Now().Add(time.Minute * 30).Unix()
	tokenString, err := token.SignedString([]byte(cfg.JWT.Secret))
	if err != nil {
		log.Error("failed to sign jwt token", logger.Err(err))
		return "", fmt.Errorf("%s: %w", op, err)
	}
	return tokenString, nil
}

func CreateLoginJWT(cfg *config.Config, log *slog.Logger, login string) (string, error) {
	const op = "jwt.CreateLoginJWT"

	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["login"] = login
	claims["exp"] = time.Now().Add(time.Minute * 30).Unix()
	tokenString, err := token.SignedString([]byte(cfg.JWT.Secret))
	if err != nil {
		log.Error("failed to sign jwt token", logger.Err(err))
		return "", fmt.Errorf("%s: %w", op, err)
	}
	return tokenString, nil
}

func GetEmail(tokenString string, secret string) (string, error) {
	const op = "jwt.GetEmail"

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secret), nil
	})

	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		if email, ok := claims["email"].(string); ok {
			return email, nil
		}
	}

	return "", fmt.Errorf("%s: invalid token", op)
}

func GetLogin(tokenString string, secret string) (string, error) {
	const op = "jwt.GetLogin"

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secret), nil
	})

	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		if login, ok := claims["login"].(string); ok {
			return login, nil
		}
	}

	return "", fmt.Errorf("%s: invalid token", op)
}
