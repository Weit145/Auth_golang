package logout

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/Weit145/Auth_golang/internal/config"
	"github.com/Weit145/Auth_golang/internal/domain"
	myjwt "github.com/Weit145/Auth_golang/internal/lib/jwt"
	"github.com/Weit145/Auth_golang/internal/lib/logger"
	"github.com/Weit145/Auth_golang/internal/storage"
	"github.com/jackc/pgx/v5"
)

type LogOut struct {
	Storage    LogOutRepo
	TxProvider storage.TxProvider
	Log        *slog.Logger
	Cfg        *config.Config
}

type LogOutRepo interface {
	GetUserByLogin(ctx context.Context, login string) (*domain.User, error)
	AuthenticateRepo(ctx context.Context, user *domain.User) error
}

func (s *LogOut) LogOutUser(ctx context.Context, AssetToken string) error {
	const op = "service.LogOutUser"

	login, err := myjwt.GetLogin(AssetToken, s.Cfg.JWT.Secret)
	if err != nil {
		s.Log.Error("failed to get login from token", slog.String("token", AssetToken), logger.Err(err))
		return fmt.Errorf("%s: %w", op, err)
	}

	err = s.TxProvider.WithTx(ctx, func(tx pgx.Tx) error {
		user, err := s.Storage.GetUserByLogin(ctx, login)
		if err != nil {
			return fmt.Errorf("%s: failed to get user by login within transaction: %w", op, err)
		}

		user.RefreshTokenHash = ""

		if err = s.Storage.AuthenticateRepo(ctx, user); err != nil {
			return fmt.Errorf("%s: failed to logout user within transaction: %w", op, err)
		}
		s.Log.Info("LogOut method called", slog.String("Token: ", AssetToken))
		return nil
	})
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
