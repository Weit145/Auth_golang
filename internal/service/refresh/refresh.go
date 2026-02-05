package refresh

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"log/slog"

	"github.com/Weit145/Auth_golang/internal/config"
	"github.com/Weit145/Auth_golang/internal/domain"
	myjwt "github.com/Weit145/Auth_golang/internal/lib/jwt"
	"github.com/Weit145/Auth_golang/internal/lib/logger"
	"github.com/jackc/pgx/v5"
)

type Refresh struct {
	Storage RefreshRepo
	Log     *slog.Logger
	Cfg     *config.Config
}

type RefreshRepo interface {
	BeginTx(ctx context.Context, opts pgx.TxOptions) (pgx.Tx, error)
	GetUserByLogin(ctx context.Context, tx pgx.Tx, login string) (*domain.User, error)
}

func (s *Refresh) Refresh(ctx context.Context, RefreshToken string) (string, error) {
	const op = "service.Refresh"

	login, err := myjwt.GetLogin(RefreshToken, s.Cfg.JWT.Secret)
	if err != nil {
		s.Log.Error("failed to get login from token", slog.String("token", RefreshToken), logger.Err(err))
		return "", fmt.Errorf("%s: %w", op, err)
	}

	tx, err := s.Storage.BeginTx(ctx, pgx.TxOptions{})
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

	user, err := s.Storage.GetUserByLogin(ctx, tx, login)
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

	refreshToken, err := myjwt.CreateLoginJWT(s.Cfg, s.Log, user.Login)
	if err != nil {
		return "", fmt.Errorf("%s: failed to create login JWT: %w", op, err)
	}

	s.Log.Info("Refresh method called", slog.String("RefreshToken: ", RefreshToken))
	return refreshToken, nil
}
