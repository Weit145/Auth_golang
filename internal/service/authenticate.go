package service

import (
	"context"
	"log/slog"
)

func (s *Service) Authenticate(ctx context.Context, login, password string) (string, error) {
	s.log.Info("Refresh method called", slog.String("RefreshToken: ", login))
	return "123123", nil
}
