package storage

import (
	"context"

	"github.com/Weit145/Auth_golang/internal/domain"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type QueryRunner interface {
	Exec(ctx context.Context, sql string, arguments ...any) (pgconn.CommandTag, error)
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
}

type Storage interface {
	RegistrationRepo(ctx context.Context, login, email, passwordHash string) error
	AuthenticateRepo(ctx context.Context, user *domain.User) error
	GetUserByEmail(ctx context.Context, email string) (*domain.User, error)
	GetUserByLogin(ctx context.Context, login string) (*domain.User, error)
	ConfirmRepo(ctx context.Context, user *domain.User) error
	UpdateRefreshToken(ctx context.Context, user *domain.User) error
}

type TxProvider interface {
	WithTx(ctx context.Context, fn func(tx pgx.Tx) error) error
}
