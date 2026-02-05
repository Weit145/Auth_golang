package service

import (
	"context"
	"log/slog"

	"github.com/Weit145/Auth_golang/internal/config"
	"github.com/Weit145/Auth_golang/internal/service/authenticate"
	"github.com/Weit145/Auth_golang/internal/service/confirm"
	"github.com/Weit145/Auth_golang/internal/service/current"
	"github.com/Weit145/Auth_golang/internal/service/logout"
	"github.com/Weit145/Auth_golang/internal/service/refresh"
	"github.com/Weit145/Auth_golang/internal/service/registration"
	"github.com/Weit145/Auth_golang/internal/storage"
)

type Service struct {
	Auth         authenticate.Login
	ConfirmUser  confirm.Confirm
	CurrentUser  current.Current
	LogOut       logout.LogOut
	RefreshUser  refresh.Refresh
	Registration registration.Registration
}

//go:generate go run github.com/vektra/mockery/v2@v2.53.5 --name=ServiceAuth
type ServiceAuth interface {
	LoginUser(ctx context.Context, login, password string) (string, string, error)
	Confirm(ctx context.Context, token string) (string, string, error)
	Current(ctx context.Context, AssetToken string) (*current.User, error)
	LogOutUser(ctx context.Context, AssetToken string) error
	Refresh(ctx context.Context, RefreshToken string) (string, error)
	CreateUser(ctx context.Context, login, email, password string) error
}

func New(log *slog.Logger, storage storage.Storage, cfg *config.Config) *Service {
	return &Service{
		Auth: authenticate.Login{
			Storage: storage,
			Cfg:     cfg,
			Log:     log,
		},
		ConfirmUser: confirm.Confirm{
			Storage: storage,
			Cfg:     cfg,
			Log:     log,
		},
		CurrentUser: current.Current{
			Storage: storage,
			Cfg:     cfg,
			Log:     log,
		},
		LogOut: logout.LogOut{
			Storage: storage,
			Cfg:     cfg,
			Log:     log,
		},
		RefreshUser: refresh.Refresh{
			Storage: storage,
			Cfg:     cfg,
			Log:     log,
		},
		Registration: registration.Registration{
			Storage: storage,
			Cfg:     cfg,
			Log:     log,
		},
	}
}

func (s *Service) LoginUser(ctx context.Context, login, password string) (string, string, error) {
	return s.Auth.LoginUser(ctx, login, password)
}

func (s *Service) Confirm(ctx context.Context, token string) (string, string, error) {
	return s.ConfirmUser.Confirm(ctx, token)
}

func (s *Service) Current(ctx context.Context, AssetToken string) (*current.User, error) {
	return s.CurrentUser.Current(ctx, AssetToken)
}

func (s *Service) LogOutUser(ctx context.Context, AssetToken string) error {
	return s.LogOut.LogOutUser(ctx, AssetToken)
}

func (s *Service) Refresh(ctx context.Context, RefreshToken string) (string, error) {
	return s.RefreshUser.Refresh(ctx, RefreshToken)
}

func (s *Service) CreateUser(ctx context.Context, login, email, password string) error {
	return s.Registration.CreateUser(ctx, login, email, password)
}
