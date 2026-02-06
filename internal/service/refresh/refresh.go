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
	"github.com/Weit145/Auth_golang/internal/storage"
	"github.com/jackc/pgx/v5"
)

type Refresh struct {
	Storage    RefreshRepo
	TxProvider storage.TxProvider
	Log        *slog.Logger
	Cfg        *config.Config
}

type RefreshRepo interface {
	GetUserByLogin(ctx context.Context, login string) (*domain.User, error)
}

func (s *Refresh) Refresh(ctx context.Context, RefreshToken string) (newRefreshToken string, err error) {
	const op = "service.Refresh"

	login, err := myjwt.GetLogin(RefreshToken, s.Cfg.JWT.Secret)
	if err != nil {
		s.Log.Error("failed to get login from token", slog.String("token", RefreshToken), logger.Err(err))
		return "", fmt.Errorf("%s: %w", op, err)
	}

	err = s.TxProvider.WithTx(ctx, func(tx pgx.Tx) error {
		user, err := s.Storage.GetUserByLogin(ctx, login)
		if err != nil {
			return fmt.Errorf("%s: failed to get user by login within transaction: %w", op, err)
		}

		h := sha256.New()
		h.Write([]byte(RefreshToken))
		check := hex.EncodeToString(h.Sum(nil))

		if user.RefreshTokenHash != check {
			return fmt.Errorf("%s: failed refreshTokenHash to DB", op)
		}

		newRefreshToken, err = myjwt.CreateLoginJWT(s.Cfg, s.Log, user.Login)
		if err != nil {
			return fmt.Errorf("%s: failed to create login JWT: %w", op, err)
		}
		s.Log.Info("Refresh method called", slog.String("RefreshToken: ", RefreshToken))
		return nil
	})
	if err != nil {
		return "", err
	}

	return newRefreshToken, nil
}
