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

func (s *Service) LoginUser(ctx context.Context, login, password string) (string, string, error) {
	const op = "service.LoginUser"

	tx, err := s.storage.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return "", "", fmt.Errorf("%s: failed to begin transaction: %w", op, err)
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
		return "", "", fmt.Errorf("%s: failed to get user by login within transaction: %w", op, err)
	}

	refreshToken, err := myjwt.CreateLoginJWT(s.cfg, s.log, user.Login)
	if err != nil {
		return "", "", fmt.Errorf("%s: failed to create login JWT: %w", op, err)
	}

	AssetToken, err := myjwt.CreateLoginJWT(s.cfg, s.log, user.Login)
	if err != nil {
		s.log.Error("failed to create login JWT", logger.Err(err))
		return "", "", fmt.Errorf("%s: %w", op, err)
	}
	h := sha256.New()
	h.Write([]byte(refreshToken))
	user.RefreshTokenHash = hex.EncodeToString(h.Sum(nil))

	if err = s.storage.AuthenticateRepo(ctx, tx, user); err != nil {
		return "", "", fmt.Errorf("%s: failed to authenticate user within transaction: %w", op, err)
	}

	if err = tx.Commit(ctx); err != nil {
		return "", "", fmt.Errorf("%s: failed to commit transaction: %w", op, err)
	}

	s.log.Info("Authenticate method called", slog.String("Login: ", login))
	return AssetToken, refreshToken, nil
}
