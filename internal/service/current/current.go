package current

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

type User struct {
	Id         int
	Login      string
	IsActive   bool
	IsVerified bool
	Role       string
}

type Current struct {
	Storage    CurrentRepo
	TxProvider storage.TxProvider
	Log        *slog.Logger
	Cfg        *config.Config
}

type CurrentRepo interface {
	GetUserByLogin(ctx context.Context, login string) (*domain.User, error)
}

func (s *Current) Current(ctx context.Context, AssetToken string) (*User, error) {
	const op = "service.Current"

	login, err := myjwt.GetLogin(AssetToken, s.Cfg.JWT.Secret)
	if err != nil {
		s.Log.Error("failed to get login from token", slog.String("token", AssetToken), logger.Err(err))
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	var resp User
	err = s.TxProvider.WithTx(ctx, func(tx pgx.Tx) error {
		user, err := s.Storage.GetUserByLogin(ctx, login)
		if err != nil {
			return fmt.Errorf("%s: failed to get user by login within transaction: %w", op, err)
		}

		resp = User{
			Id:         int(user.Id),
			Login:      user.Login,
			IsActive:   user.IsActive,
			IsVerified: user.IsVerified,
			Role:       user.Role,
		}
		s.Log.Info("Current method called", slog.String("AssetToken: ", AssetToken))
		return nil
	})
	if err != nil {
		return nil, err
	}

	return &resp, nil
}
