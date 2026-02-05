package current

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/Weit145/Auth_golang/internal/config"
	"github.com/Weit145/Auth_golang/internal/domain"
	myjwt "github.com/Weit145/Auth_golang/internal/lib/jwt"
	"github.com/Weit145/Auth_golang/internal/lib/logger"
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
	Storage CurrentRepo
	Log     *slog.Logger
	Cfg     *config.Config
}

type CurrentRepo interface {
	BeginTx(ctx context.Context, opts pgx.TxOptions) (pgx.Tx, error)
	GetUserByLogin(ctx context.Context, tx pgx.Tx, login string) (*domain.User, error)
}

func (s *Current) Current(ctx context.Context, AssetToken string) (*User, error) {
	const op = "service.Current"

	login, err := myjwt.GetLogin(AssetToken, s.Cfg.JWT.Secret)
	if err != nil {
		s.Log.Error("failed to get login from token", slog.String("token", AssetToken), logger.Err(err))
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	tx, err := s.Storage.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return nil, fmt.Errorf("%s: failed to begin transaction: %w", op, err)
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
		return nil, fmt.Errorf("%s: failed to get user by login within transaction: %w", op, err)
	}

	if err = tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("%s: failed to commit transaction: %w", op, err)
	}

	resp := User{
		Id:         int(user.Id),
		Login:      user.Login,
		IsActive:   user.IsActive,
		IsVerified: user.IsVerified,
		Role:       user.Role,
	}
	s.Log.Info("Refresh method called", slog.String("AssetToken: ", AssetToken))
	return &resp, nil
}
