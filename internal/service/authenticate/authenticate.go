package authenticate

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

type Login struct {
	Storage AuthRepo
	Cfg     *config.Config
	Log     *slog.Logger
}

type AuthRepo interface {
	BeginTx(ctx context.Context, opts pgx.TxOptions) (pgx.Tx, error)
	GetUserByLogin(ctx context.Context, tx pgx.Tx, login string) (*domain.User, error)
	AuthenticateRepo(ctx context.Context, tx pgx.Tx, user *domain.User) error
}

func (s Login) LoginUser(ctx context.Context, login, password string) (string, string, error) {
	const op = "service.LoginUser"

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

	user, err := s.Storage.GetUserByLogin(ctx, tx, login)
	if err != nil {
		return "", "", fmt.Errorf("%s: failed to get user by login within transaction: %w", op, err)
	}

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

	if err = s.Storage.AuthenticateRepo(ctx, tx, user); err != nil {
		return "", "", fmt.Errorf("%s: failed to authenticate user within transaction: %w", op, err)
	}

	if err = tx.Commit(ctx); err != nil {
		return "", "", fmt.Errorf("%s: failed to commit transaction: %w", op, err)
	}

	s.Log.Info("Authenticate method called", slog.String("Login: ", login))
	return AssetToken, refreshToken, nil
}
