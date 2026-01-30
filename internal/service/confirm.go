package service

import (
	"context"
	"log/slog"

	GRPCauth "github.com/Weit145/proto-repo/auth"
)

func (s *Service) Confirm(ctx context.Context, req *GRPCauth.UserCreateRequest) (*GRPCauth.Okey, error) {
	s.log.Info("CreateUser method called", slog.String("email", req.GetEmail()))
	return &GRPCauth.Okey{Success: true}, nil
}
