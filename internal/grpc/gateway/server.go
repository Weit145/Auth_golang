package gateway

import (
	"context"
	"log/slog"
	"net"
	"os"

	"github.com/Weit145/Auth_golang/internal/lib/logger"
	"github.com/Weit145/Auth_golang/internal/service"
	GRPCauth "github.com/Weit145/proto-repo/auth"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type server struct {
	GRPCauth.UnimplementedAuthServer
	service *service.Service
	log     *slog.Logger
}

func New(log *slog.Logger, serv *service.Service, addr string) (*grpc.Server, error) {

	lis, err := net.Listen("tcp", addr)
	if err != nil {
		return nil, err
	}

	s := grpc.NewServer()

	GRPCauth.RegisterAuthServer(s, &server{service: serv, log: log})

	go func() {
		log.Info("gRPC server started", slog.String("addr", addr))
		if err := s.Serve(lis); err != nil {
			log.Error("gRPC server failed", logger.Err(err))
			os.Exit(1)
		}
	}()
	return s, nil
}

func (s *server) CreateUser(ctx context.Context, req *GRPCauth.UserCreateRequest) (*GRPCauth.Okey, error) {
	login := req.GetLogin()
	email := req.GetEmail()
	password := req.GetPassword()
	username := req.GetUsername()

	if login == "" {
		return nil, status.Error(codes.InvalidArgument, "login is required")
	}
	if email == "" {
		return nil, status.Error(codes.InvalidArgument, "email is required")
	}
	if password == "" {
		return nil, status.Error(codes.InvalidArgument, "password is required")
	}
	if username == "" {
		return nil, status.Error(codes.InvalidArgument, "username is required")
	}

	err := s.service.CreateUser(ctx, login, email, password, username)
	if err != nil {
		return nil, err
	}
	ok := GRPCauth.Okey{Success: true}
	return &ok, nil
}

func (s *server) RegistrationUser(ctx context.Context, req *GRPCauth.TokenRequest) (*GRPCauth.CookieResponse, error) {
	token := req.GetTokenPod()
	if token == "" {
		return nil, status.Error(codes.InvalidArgument, "token is required")
	}
	RefreshToken, err := s.service.Confirm(ctx, token)
	if err != nil {
		return nil, err
	}
	resp := GRPCauth.CookieResponse{
		AccessToken: "123123",
		Cookie: &GRPCauth.Cookie{
			Key:      "refresh_token",
			Value:    RefreshToken,
			Httponly: true,
			Secure:   true,
			Samesite: "lax",
			MaxAge:   24,
		},
	}
	return &resp, nil
}

func (s *server) RefreshToken(ctx context.Context, req *GRPCauth.CookieRequest) (*GRPCauth.AccessTokenResponse, error) {
	RefreshToken := req.GetRefreshToken()
	if RefreshToken == "" {
		return nil, status.Error(codes.InvalidArgument, "RefreshToken is required")
	}
	AssetToken, err := s.service.Refresh(ctx, RefreshToken)
	if err != nil {
		return nil, err
	}
	resp := GRPCauth.AccessTokenResponse{
		AccessToken: AssetToken,
	}
	return &resp, nil
}

func (s *server) Authenticate(ctx context.Context, req *GRPCauth.UserLoginRequest) (*GRPCauth.CookieResponse, error) {
	login := req.GetLogin()
	password := req.GetPassword()
	if login == "" {
		return nil, status.Error(codes.InvalidArgument, "Login is required")
	}
	if password == "" {
		return nil, status.Error(codes.InvalidArgument, "Password is required")
	}
	RefreshToken, err := s.service.LoginUser(ctx, login, password)
	if err != nil {
		return nil, err
	}
	resp := GRPCauth.CookieResponse{
		AccessToken: "123123",
		Cookie: &GRPCauth.Cookie{
			Key:      "refresh_token",
			Value:    RefreshToken,
			Httponly: true,
			Secure:   true,
			Samesite: "lax",
			MaxAge:   24,
		},
	}
	return &resp, nil
}

func (s *server) CurrentUser(ctx context.Context, req *GRPCauth.UserCurrentRequest) (*GRPCauth.CurrentUserResponse, error) {
	AssetToken := req.GetAccessToken()
	if AssetToken == "" {
		return nil, status.Error(codes.InvalidArgument, "AssetToken is required")
	}
	user, err := s.service.Current(ctx, AssetToken)
	if err != nil {
		return nil, err
	}
	resp := GRPCauth.CurrentUserResponse{
		Id:         int32(user.Id),
		Login:      user.Login,
		IsActive:   user.IsActive,
		IsVerified: user.IsVerified,
		Role:       user.Role,
	}
	return &resp, nil
}

func (s *server) LogOutUser(ctx context.Context, req *GRPCauth.TokenRequest) (*GRPCauth.Empty, error) {
	AssetToken := req.GetTokenPod()
	if AssetToken == "" {
		return nil, status.Error(codes.InvalidArgument, "AssetToken is required")
	}

	err := s.service.LogOutUser(ctx, AssetToken)
	if err != nil {
		return nil, err
	}

	resp := GRPCauth.Empty{}
	return &resp, nil
}
