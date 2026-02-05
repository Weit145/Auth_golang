package gateway

import (
	"context"
	"log/slog"
	"net"
	"os"

	"github.com/Weit145/Auth_golang/internal/lib/logger"
	"github.com/Weit145/Auth_golang/internal/service"
	pb "github.com/Weit145/proto-repo/auth"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Server struct {
	pb.UnimplementedAuthServer
	Service service.ServiceAuth
	Log     *slog.Logger
}

func New(Log *slog.Logger, serv service.ServiceAuth, lis net.Listener) (*grpc.Server, error) {

	s := grpc.NewServer()

	pb.RegisterAuthServer(s, &Server{Service: serv, Log: Log})

	go func() {
		Log.Info("gRPC server started", slog.String("addr", lis.Addr().String()))
		if err := s.Serve(lis); err != nil {
			Log.Error("gRPC server failed", logger.Err(err))
			os.Exit(1)
		}
	}()
	return s, nil
}

func (s *Server) CreateUser(ctx context.Context, req *pb.UserCreateRequest) (*pb.Okey, error) {
	login := req.GetLogin()
	email := req.GetEmail()
	password := req.GetPassword()

	if login == "" {
		return nil, status.Error(codes.InvalidArgument, "login is required")
	}
	if email == "" {
		return nil, status.Error(codes.InvalidArgument, "email is required")
	}
	if password == "" {
		return nil, status.Error(codes.InvalidArgument, "password is required")
	}

	s.Log.Info("Calling Service.CreateUser", slog.String("login", login), slog.String("email", email))
	err := s.Service.CreateUser(ctx, login, email, password)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to create user")
	}
	resp := pb.Okey{Success: true}
	return &resp, nil
}

func (s *Server) RegistrationUser(ctx context.Context, req *pb.TokenRequest) (*pb.CookieResponse, error) {
	token := req.GetTokenPod()
	if token == "" {
		return nil, status.Error(codes.InvalidArgument, "token is required")
	}
	AssetToken, RefreshToken, err := s.Service.Confirm(ctx, token)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to confirm user")
	}
	resp := pb.CookieResponse{
		AccessToken: AssetToken,
		Cookie: &pb.Cookie{
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

func (s *Server) RefreshToken(ctx context.Context, req *pb.CookieRequest) (*pb.AccessTokenResponse, error) {
	RefreshToken := req.GetRefreshToken()
	if RefreshToken == "" {
		return nil, status.Error(codes.InvalidArgument, "RefreshToken is required")
	}
	AssetToken, err := s.Service.Refresh(ctx, RefreshToken)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to refresh token")
	}
	resp := pb.AccessTokenResponse{
		AccessToken: AssetToken,
	}
	return &resp, nil
}

func (s *Server) Authenticate(ctx context.Context, req *pb.UserLoginRequest) (*pb.CookieResponse, error) {
	login := req.GetLogin()
	password := req.GetPassword()
	if login == "" {
		return nil, status.Error(codes.InvalidArgument, "Login is required")
	}
	if password == "" {
		return nil, status.Error(codes.InvalidArgument, "Password is required")
	}
	AccessToken, RefreshToken, err := s.Service.LoginUser(ctx, login, password)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to authenticate user")
	}
	resp := pb.CookieResponse{
		AccessToken: AccessToken,
		Cookie: &pb.Cookie{
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

func (s *Server) CurrentUser(ctx context.Context, req *pb.UserCurrentRequest) (*pb.CurrentUserResponse, error) {
	AssetToken := req.GetAccessToken()
	if AssetToken == "" {
		return nil, status.Error(codes.InvalidArgument, "AssetToken is required")
	}
	user, err := s.Service.Current(ctx, AssetToken)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to get current user")
	}
	resp := pb.CurrentUserResponse{
		Id:         int32(user.Id),
		Login:      user.Login,
		IsActive:   user.IsActive,
		IsVerified: user.IsVerified,
		Role:       user.Role,
	}
	return &resp, nil
}

func (s *Server) LogOutUser(ctx context.Context, req *pb.TokenRequest) (*pb.Empty, error) {
	AssetToken := req.GetTokenPod()
	if AssetToken == "" {
		return nil, status.Error(codes.InvalidArgument, "AssetToken is required")
	}

	err := s.Service.LogOutUser(ctx, AssetToken)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to Log out user")
	}

	resp := pb.Empty{}
	return &resp, nil
}
