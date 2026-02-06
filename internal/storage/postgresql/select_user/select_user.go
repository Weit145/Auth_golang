package select_user

import (
	"context"
	"fmt"

	"github.com/Weit145/Auth_golang/internal/domain"
	"github.com/Weit145/Auth_golang/internal/storage"
	"github.com/jackc/pgx/v5"
)

func GetUserByEmailOp(ctx context.Context, runner storage.QueryRunner, email string) (*domain.User, error) {
	const op = "storage.postgresql.select_user.GetUserByEmailOp"

	stmt := `SELECT id, login, email, password_hash, is_active, is_verified, role, refresh_token_hash FROM auth WHERE email = $1`
	var user domain.User
	err := runner.QueryRow(ctx, stmt, email).Scan(
		&user.Id,
		&user.Login,
		&user.Email,
		&user.PasswordHash,
		&user.IsActive,
		&user.IsVerified,
		&user.Role,
		&user.RefreshTokenHash,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("%s: user not found", op)
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &user, nil
}

func GetUserByLoginOp(ctx context.Context, runner storage.QueryRunner, login string) (*domain.User, error) {
	const op = "storage.postgresql.select_user.GetUserByLoginOp"

	stmt := `SELECT id, login, email, password_hash, is_active, is_verified, role, refresh_token_hash FROM auth WHERE login = $1`
	var user domain.User
	err := runner.QueryRow(ctx, stmt, login).Scan(
		&user.Id,
		&user.Login,
		&user.Email,
		&user.PasswordHash,
		&user.IsActive,
		&user.IsVerified,
		&user.Role,
		&user.RefreshTokenHash,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("%s: user not found", op)
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &user, nil
}
