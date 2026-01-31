package service

import (
	"context"
	"log/slog"
)

func (s *Service) LoginUser(ctx context.Context, login, password string) (string, error) {
	s.log.Info("Refresh method called", slog.String("Login: ", login))
	return "123123", nil
}
