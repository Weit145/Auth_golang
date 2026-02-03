package postgresql

import (
	"context"
	"fmt"

	"github.com/Weit145/Auth_golang/internal/domain"
)

func (s *Storage) ConfirmRepo(ctx context.Context, user *domain.User) error {
	const op = "storage.postgresql.ConfirmRepo"

	stmt := `UPDATE auth SET is_verified = $1, refresh_token_hash = $2 WHERE id = $3`
	_, err := s.db.Exec(ctx, stmt, user.IsVerified, user.RefreshTokenHash, user.Id)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil

}
