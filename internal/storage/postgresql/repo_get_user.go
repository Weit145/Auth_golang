package postgresql

import (
	"context"
	"fmt"

	"github.com/Weit145/Auth_golang/internal/domain"
	"github.com/jackc/pgx/v5"
)

func (s *Storage) GetUserByEmail(ctx context.Context, email string) (*domain.User, error) {
	const op = "storage.postgresql.GetUserByEmail"

	stmt := `SELECT id, login, email, password_hash, is_active, is_verified, role, refresh_token_hash FROM auth WHERE email = $1`
	var user domain.User
	err := s.db.QueryRow(ctx, stmt, email).Scan(
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
