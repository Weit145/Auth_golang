package confirm

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

type Confirm struct {
	Storage ConfirmRepo
	Cfg     *config.Config
	Log     *slog.Logger
}

type ConfirmRepo interface {
	BeginTx(ctx context.Context, opts pgx.TxOptions) (pgx.Tx, error)
	GetUserByEmail(ctx context.Context, tx pgx.Tx, email string) (*domain.User, error)
	ConfirmRepo(ctx context.Context, tx pgx.Tx, user *domain.User) error
}

func (s *Confirm) Confirm(ctx context.Context, token string) (string, string, error) {
	const op = "service.Confirm"

	email, err := myjwt.GetEmail(token, s.Cfg.JWT.Secret)
	if err != nil {
		s.Log.Error("failed to get email from token", slog.String("token", token), logger.Err(err))
		return "", "", fmt.Errorf("%s: %w", op, err)
	}

	tx, err := s.Storage.BeginTx(ctx, pgx.TxOptions{})
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

	user, err := s.Storage.GetUserByEmail(ctx, tx, email)
	if err != nil {
		return "", "", fmt.Errorf("%s: failed to get user by email within transaction: %w", op, err)
	}

	user.IsVerified = true

	refreshToken, err := myjwt.CreateLoginJWT(s.Cfg, s.Log, user.Login)
	if err != nil {
		return "", "", fmt.Errorf("%s: failed to create login JWT: %w", op, err)
	}

	AssetToken, err := myjwt.CreateLoginJWT(s.Cfg, s.Log, user.Login)
	if err != nil {
		s.Log.Error("failed to create login JWT", logger.Err(err))
		return "", "", fmt.Errorf("%s: %w", op, err)
	}

	h := sha256.New()
	h.Write([]byte(refreshToken))
	user.RefreshTokenHash = hex.EncodeToString(h.Sum(nil))

	if err = s.Storage.ConfirmRepo(ctx, tx, user); err != nil {
		return "", "", fmt.Errorf("%s: failed to update user within transaction: %w", op, err)
	}

	if err = tx.Commit(ctx); err != nil {
		return "", "", fmt.Errorf("%s: failed to commit transaction: %w", op, err)
	}

	s.Log.Info("Confirm method called", slog.String("token: ", token))
	return AssetToken, refreshToken, nil
}
