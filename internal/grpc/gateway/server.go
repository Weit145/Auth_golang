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
