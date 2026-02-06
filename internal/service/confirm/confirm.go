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
	"github.com/Weit145/Auth_golang/internal/storage"
	"github.com/jackc/pgx/v5"
)

type Confirm struct {
	Storage    ConfirmRepo
	TxProvider storage.TxProvider
	Cfg        *config.Config
	Log        *slog.Logger
}

type ConfirmRepo interface {
	GetUserByEmail(ctx context.Context, email string) (*domain.User, error)
	ConfirmRepo(ctx context.Context, user *domain.User) error
}

func (s *Confirm) Confirm(ctx context.Context, token string) (accessToken, refreshToken string, err error) {
	const op = "service.Confirm"

	email, err := myjwt.GetEmail(token, s.Cfg.JWT.Secret)
	if err != nil {
		s.Log.Error("failed to get email from token", slog.String("token", token), logger.Err(err))
		return "", "", fmt.Errorf("%s: %w", op, err)
	}

	err = s.TxProvider.WithTx(ctx, func(tx pgx.Tx) error {
		user, err := s.Storage.GetUserByEmail(ctx, email)
		if err != nil {
			return fmt.Errorf("%s: failed to get user by email within transaction: %w", op, err)
		}

		user.IsVerified = true

		refreshToken, err = myjwt.CreateLoginJWT(s.Cfg, s.Log, user.Login)
		if err != nil {
			return fmt.Errorf("%s: failed to create login JWT: %w", op, err)
		}

		accessToken, err = myjwt.CreateLoginJWT(s.Cfg, s.Log, user.Login)
		if err != nil {
			s.Log.Error("failed to create login JWT", logger.Err(err))
			return fmt.Errorf("%s: %w", op, err)
		}

		h := sha256.New()
		h.Write([]byte(refreshToken))
		user.RefreshTokenHash = hex.EncodeToString(h.Sum(nil))

		if err = s.Storage.ConfirmRepo(ctx, user); err != nil {
			return fmt.Errorf("%s: failed to update user within transaction: %w", op, err)
		}
		s.Log.Info("Confirm method called", slog.String("token: ", token))
		return nil
	})
	if err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}
