package updateverified

import (
	"context"
	"fmt"

	"github.com/Weit145/Auth_golang/internal/domain"
	"github.com/Weit145/Auth_golang/internal/storage"
)

func UpdateVerifiedOp(ctx context.Context, runner storage.QueryRunner, user *domain.User) error {
	const op = "storage.postgresql.updateverified.UpdateVerifiedOp"

	stmt := `UPDATE auth SET is_verified = $1, refresh_token_hash = $2 WHERE id = $3`
	_, err := runner.Exec(ctx, stmt, user.IsVerified, user.RefreshTokenHash, user.Id)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}
