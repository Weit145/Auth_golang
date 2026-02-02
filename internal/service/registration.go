package service

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/Weit145/Auth_golang/internal/lib/logger"
)

func (s *Service) CreateUser(ctx context.Context, login, email, passwordHash string) error {
	const op = "service.CreateUser"

	s.log.Info("CreateUser method called", slog.String("email", email), slog.String("login", login))

	err := s.db.RegistrationRepo(ctx, login, email, passwordHash)
	if err != nil {
		s.log.Error("failed to register user", logger.Err(err))
		return fmt.Errorf("%s: %w", op, err)
	}

	s.log.Info("user registered successfully", slog.String("email", email), slog.String("login", login))
	return nil
}
