package service

import (
	"context"
	"fmt"
	"log/slog"

	myjwt "github.com/Weit145/Auth_golang/internal/lib/jwt"
	"github.com/Weit145/Auth_golang/internal/lib/logger"
	"github.com/jackc/pgx/v5"
)

func (s *Service) LogOutUser(ctx context.Context, AssetToken string) error {
	const op = "service.LogOutUser"

	login, err := myjwt.GetLogin(AssetToken, s.cfg.JWT.Secret)
	if err != nil {
		s.log.Error("failed to get login from token", slog.String("token", AssetToken), logger.Err(err))
		return fmt.Errorf("%s: %w", op, err)
	}

	tx, err := s.storage.BeginTx(ctx, pgx.TxOptions{})
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

	user, err := s.storage.GetUserByLogin(ctx, tx, login)
	if err != nil {
		return fmt.Errorf("%s: failed to get user by login within transaction: %w", op, err)
	}

	user.RefreshTokenHash = ""

	if err = s.storage.AuthenticateRepo(ctx, tx, user); err != nil {
		return fmt.Errorf("%s: failed to lgout user within transaction: %w", op, err)
	}

	if err = tx.Commit(ctx); err != nil {
		return fmt.Errorf("%s: failed to commit transaction: %w", op, err)
	}

	s.log.Info("LogOut method called", slog.String("Token: ", AssetToken))
	return nil
}
