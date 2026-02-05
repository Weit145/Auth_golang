package registration

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/Weit145/Auth_golang/internal/config"
	myjwt "github.com/Weit145/Auth_golang/internal/lib/jwt"
	"github.com/Weit145/Auth_golang/internal/lib/logger"
	"golang.org/x/crypto/bcrypt"
)

type Registration struct {
	Storage RegistrationRepo
	Log     *slog.Logger
	Cfg     *config.Config
}

type RegistrationRepo interface {
	RegistrationRepo(ctx context.Context, login, email, passwordHash string) error
}

func (s *Registration) CreateUser(ctx context.Context, login, email, password string) error {
	const op = "service.CreateUser"

	s.Log.Info("CreateUser method called", slog.String("email", email), slog.String("login", login))

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		s.Log.Error("failed to generate password hash", logger.Err(err))
		return fmt.Errorf("%s: %w", op, err)
	}

	err = s.Storage.RegistrationRepo(ctx, login, email, string(passwordHash))
	if err != nil {
		s.Log.Error("failed to register user", logger.Err(err))
		return fmt.Errorf("%s: %w", op, err)
	}
	token, err := myjwt.CreateEmailJWT(s.Cfg, s.Log, email)
	if err != nil {
		s.Log.Error("failed to create email JWT", logger.Err(err))
		return fmt.Errorf("%s: %w", op, err)
	}

	s.Log.Info("user registered successfully", slog.String("email", email), slog.String("login", login))
	s.Log.Info("Token", slog.String("token", token))
	return nil
}
