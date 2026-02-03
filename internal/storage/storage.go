package storage

import (
	"context"

	"github.com/Weit145/Auth_golang/internal/domain"
)

type Storage interface {
	RegistrationRepo(ctx context.Context, login, email, passwordHash string) error
	GetUserByEmail(ctx context.Context, email string) (*domain.User, error)
	ConfirmRepo(ctx context.Context, user *domain.User) error
}
