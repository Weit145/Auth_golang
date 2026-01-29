package gateway

import (
	"net"
	"os"

	"log/slog"

	"github.com/Weit145/Auth_golang/internal/lib/logger"
	GRPCauth "github.com/Weit145/proto-repo/auth"
	"google.golang.org/grpc"
)

type server struct {
	GRPCauth.UnimplementedAuthServer
}

func New(log *slog.Logger) (*grpc.Server, error) {
	port := ":50051"
	lis, err := net.Listen("tcp", port)
	if err != nil {
		return nil, err
	}

	s := grpc.NewServer()

	GRPCauth.RegisterAuthServer(s, &server{})

	go func() {
		if err := s.Serve(lis); err != nil {
			log.Error("grpc serve failed", logger.Err(err))
			os.Exit(1)
		}
	}()
	return s, nil
}
