package service

import (
	"context"
	"fmt"
	"log/slog"

	myjwt "github.com/Weit145/Auth_golang/internal/lib/jwt"
	"github.com/Weit145/Auth_golang/internal/lib/logger"
	"golang.org/x/crypto/bcrypt"
)

func (s *Service) CreateUser(ctx context.Context, login, email, password string) error {
	const op = "service.CreateUser"

	s.log.Info("CreateUser method called", slog.String("email", email), slog.String("login", login))

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		s.log.Error("failed to generate password hash", logger.Err(err))
		return fmt.Errorf("%s: %w", op, err)
	}

	err = s.storage.RegistrationRepo(ctx, login, email, string(passwordHash))
	if err != nil {
		s.log.Error("failed to register user", logger.Err(err))
		return fmt.Errorf("%s: %w", op, err)
	}
	token, err := myjwt.CreateEmailJWT(s.cfg, s.log, email)
	if err != nil {
		s.log.Error("failed to create email JWT", logger.Err(err))
		return fmt.Errorf("%s: %w", op, err)
	}

	s.log.Info("user registered successfully", slog.String("email", email), slog.String("login", login))
	s.log.Info("Token", slog.String("token", token))
	return nil
}
