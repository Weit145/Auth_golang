package postgresql

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgconn"
)

const (
	UniqueViolation = "23505" // PostgreSQL error code for unique_violation
)

func (s *Storage) RegistrationRepo(ctx context.Context, login, email, passwordHash string) error {
	const op = "storage.postgresql.RegistrationRepo"

	tx, err := s.db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("%s: failed to begin transaction: %w", op, err)
	}
	defer tx.Rollback(ctx)

	stmt := `INSERT INTO auth (login, email, password_hash) VALUES ($1, $2, $3)`
	_, err = tx.Exec(ctx, stmt, login, email, passwordHash)
	if err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok && pgErr.Code == UniqueViolation {
			if pgErr.ConstraintName == "auth_login_key" {

				return fmt.Errorf("%s: login already exists: %w", op, err)
			}
			if pgErr.ConstraintName == "auth_email_key" {
				return fmt.Errorf("%s: email already exists: %w", op, err)
			}
		}
		return fmt.Errorf("%s: failed to insert user: %w", op, err)
	}

	return tx.Commit(ctx)
}
