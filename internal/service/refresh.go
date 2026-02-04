package service

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"log/slog"

	myjwt "github.com/Weit145/Auth_golang/internal/lib/jwt"
	"github.com/Weit145/Auth_golang/internal/lib/logger"
	"github.com/jackc/pgx/v5"
)

func (s *Service) Refresh(ctx context.Context, RefreshToken string) (string, error) {
	const op = "service.Refresh"

	login, err := myjwt.GetLogin(RefreshToken, s.cfg.JWT.Secret)
	if err != nil {
		s.log.Error("failed to get login from token", slog.String("token", RefreshToken), logger.Err(err))
		return "", fmt.Errorf("%s: %w", op, err)
	}

	tx, err := s.storage.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return "", fmt.Errorf("%s: failed to begin transaction: %w", op, err)
	}

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback(ctx)
			panic(r)
		} else if err != nil {
			tx.Rollback(ctx)
		}
	}()

	user, err := s.storage.GetUserByLogin(ctx, tx, login)
	if err != nil {
		return "", fmt.Errorf("%s: failed to get user by login within transaction: %w", op, err)
	}

	if err = tx.Commit(ctx); err != nil {
		return "", fmt.Errorf("%s: failed to commit transaction: %w", op, err)
	}

	h := sha256.New()
	h.Write([]byte(RefreshToken))
	check := hex.EncodeToString(h.Sum(nil))

	if user.RefreshTokenHash != check {
		return "", fmt.Errorf("%s: failed refreshTokenHash to DB", op)
	}

	refreshToken, err := myjwt.CreateLoginJWT(s.cfg, s.log, user.Login)
	if err != nil {
		return "", fmt.Errorf("%s: failed to create login JWT: %w", op, err)
	}

	s.log.Info("Refresh method called", slog.String("RefreshToken: ", RefreshToken))
	return refreshToken, nil
}
