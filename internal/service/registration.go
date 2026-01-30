package service

import (
	"context"
	"log/slog"
)

func (s *Service) CreateUser(ctx context.Context, login, email, password, username string) error {
	s.log.Info("CreateUser method called", slog.String("email", email))
	return nil
}
