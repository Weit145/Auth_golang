package postgresql

import (
	"context"
	"fmt"

	"github.com/Weit145/Auth_golang/internal/domain"
	"github.com/jackc/pgx/v5"
)

func (s *Storage) AuthenticateRepo(ctx context.Context, tx pgx.Tx, user *domain.User) error {
	const op = "storage.postgresql.RefreshRepo"

	stmt := `UPDATE auth SET refresh_token_hash = $1 WHERE id = $2`
	_, err := tx.Exec(ctx, stmt, user.RefreshTokenHash, user.Id)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil

}
