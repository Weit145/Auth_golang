package service

import (
	"context"
	"log/slog"
)

func (s *Service) LogOutUser(ctx context.Context, AssetToken string) error {
	s.log.Info("Refresh method called", slog.String("Token: ", AssetToken))
	return nil
}
