package postgresql

import (
	"context"
	"fmt"

	"github.com/Weit145/Auth_golang/internal/domain"
	"github.com/jackc/pgx/v5"
)

func (s *Storage) ConfirmRepo(ctx context.Context, tx pgx.Tx, user *domain.User) error {
	const op = "storage.postgresql.ConfirmRepo"

	stmt := `UPDATE auth SET is_verified = $1, refresh_token_hash = $2 WHERE id = $3`
	_, err := tx.Exec(ctx, stmt, user.IsVerified, user.RefreshTokenHash, user.Id)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil

}
