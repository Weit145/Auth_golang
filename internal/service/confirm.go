package service

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"log/slog"

	myjwt "github.com/Weit145/Auth_golang/internal/lib/jwt"
	"github.com/Weit145/Auth_golang/internal/lib/logger"
)

func (s *Service) Confirm(ctx context.Context, token string) (string, string, error) {
	const op = "service.Confirm"
	s.log.Info("Confirm method called", slog.String("token: ", token))

	email, err := myjwt.GetEmail(token, s.cfg.JWT.Secret)
	if err != nil {
		s.log.Error("failed to get email from token", slog.String("token", token), logger.Err(err))
		return "", "", fmt.Errorf("%s: %w", op, err)
	}

	user, err := s.storage.GetUserByEmail(ctx, email)
	if err != nil {
		s.log.Error("failed to get user by email", logger.Err(err))
		return "", "", fmt.Errorf("%s: %w", op, err)
	}

	user.IsVerified = true

	refreshToken, err := myjwt.CreateLoginJWT(s.cfg, s.log, user.Login)
	if err != nil {
		s.log.Error("failed to create login JWT", logger.Err(err))
		return "", "", fmt.Errorf("%s: %w", op, err)
	}
	AssetToken, err := myjwt.CreateLoginJWT(s.cfg, s.log, user.Login)
	if err != nil {
		s.log.Error("failed to create login JWT", logger.Err(err))
		return "", "", fmt.Errorf("%s: %w", op, err)
	}

	h := sha256.New()
	h.Write([]byte(refreshToken))
	user.RefreshTokenHash = hex.EncodeToString(h.Sum(nil))

	if err := s.storage.ConfirmRepo(ctx, user); err != nil {
		s.log.Error("failed to update user", logger.Err(err))
		return "", "", fmt.Errorf("%s: %w", op, err)
	}

	return AssetToken, refreshToken, nil
}
