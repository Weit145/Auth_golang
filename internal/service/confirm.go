package service

import (
	"context"
	"log/slog"
)

func (s *Service) Confirm(ctx context.Context, token string) (string, error) {
	s.log.Info("Confirm method called", slog.String("token: ", token))
	return "123123", nil
}
