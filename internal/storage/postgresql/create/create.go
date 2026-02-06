package create

import (
	"context"
	"fmt"

	"github.com/Weit145/Auth_golang/internal/storage"
	"github.com/jackc/pgx/v5/pgconn"
)

const (
	UniqueViolation = "23505"
)

func CreateUserOp(ctx context.Context, runner storage.QueryRunner, login, email, passwordHash string) error {
	const op = "storage.postgresql.create.CreateUserOp"

	stmt := `INSERT INTO auth (login, email, password_hash) VALUES ($1, $2, $3)`
	_, err := runner.Exec(ctx, stmt, login, email, passwordHash)
	if err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok && pgErr.Code == UniqueViolation {
			if pgErr.ConstraintName == "auth_login_key" {
				return fmt.Errorf("%s: login already exists", op)
			}
			if pgErr.ConstraintName == "auth_email_key" {
				return fmt.Errorf("%s: email already exists", op)
			}
		}
		return fmt.Errorf("%s: failed to insert user: %w", op, err)
	}

	return nil
}
