package storage

import (
	"context"

	"github.com/jackc/pgx/v5"

	"github.com/Weit145/Auth_golang/internal/domain"
)

type Storage interface {
	RegistrationRepo(ctx context.Context, login, email, passwordHash string) error
	AuthenticateRepo(ctx context.Context, tx pgx.Tx, user *domain.User) error
	GetUserByEmail(ctx context.Context, tx pgx.Tx, email string) (*domain.User, error)
	GetUserByLogin(ctx context.Context, tx pgx.Tx, login string) (*domain.User, error)
	ConfirmRepo(ctx context.Context, tx pgx.Tx, user *domain.User) error

	BeginTx(ctx context.Context, opts pgx.TxOptions) (pgx.Tx, error)
}
