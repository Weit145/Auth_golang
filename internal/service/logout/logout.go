package logout

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

type LogOut struct {
	Storage LogOutRepo
	Log     *slog.Logger
	Cfg     *config.Config
}

type LogOutRepo interface {
	BeginTx(ctx context.Context, opts pgx.TxOptions) (pgx.Tx, error)
	GetUserByLogin(ctx context.Context, tx pgx.Tx, login string) (*domain.User, error)
	AuthenticateRepo(ctx context.Context, tx pgx.Tx, user *domain.User) error
}

func (s *LogOut) LogOutUser(ctx context.Context, AssetToken string) error {
	const op = "service.LogOutUser"

	login, err := myjwt.GetLogin(AssetToken, s.Cfg.JWT.Secret)
	if err != nil {
		s.Log.Error("failed to get login from token", slog.String("token", AssetToken), logger.Err(err))
		return fmt.Errorf("%s: %w", op, err)
	}

	tx, err := s.Storage.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return fmt.Errorf("%s: failed to begin transaction: %w", op, err)
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
		return fmt.Errorf("%s: failed to get user by login within transaction: %w", op, err)
	}

	user.RefreshTokenHash = ""

	if err = s.Storage.AuthenticateRepo(ctx, tx, user); err != nil {
		return fmt.Errorf("%s: failed to lgout user within transaction: %w", op, err)
	}

	if err = tx.Commit(ctx); err != nil {
		return fmt.Errorf("%s: failed to commit transaction: %w", op, err)
	}

	s.Log.Info("LogOut method called", slog.String("Token: ", AssetToken))
	return nil
}
