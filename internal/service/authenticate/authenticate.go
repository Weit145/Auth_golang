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
	"github.com/Weit145/Auth_golang/internal/storage"
	"github.com/jackc/pgx/v5"
)

type Login struct {
	Storage    AuthRepo
	TxProvider storage.TxProvider
	Cfg        *config.Config
	Log        *slog.Logger
}

type AuthRepo interface {
	GetUserByLogin(ctx context.Context, login string) (*domain.User, error)
	AuthenticateRepo(ctx context.Context, user *domain.User) error
}

func (s Login) LoginUser(ctx context.Context, login, password string) (accessToken, refreshToken string, err error) {
	const op = "service.LoginUser"

	err = s.TxProvider.WithTx(ctx, func(tx pgx.Tx) error {
		user, err := s.Storage.GetUserByLogin(ctx, login)
		if err != nil {
			return fmt.Errorf("%s: failed to get user by login within transaction: %w", op, err)
		}

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

		if err = s.Storage.AuthenticateRepo(ctx, user); err != nil {
			return fmt.Errorf("%s: failed to authenticate user within transaction: %w", op, err)
		}

		s.Log.Info("Authenticate method called", slog.String("Login: ", login))
		return nil
	})
	if err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}
