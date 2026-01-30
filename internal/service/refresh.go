package service

import (
	"context"
	"log/slog"
)

func (s *Service) Refresh(ctx context.Context, RefreshToken string) (string, error) {
	s.log.Info("Refresh method called", slog.String("RefreshToken: ", RefreshToken))
	return "123123", nil
}
