package gateway

import (
	"flag"
	"fmt"
	"log"
	"net"

	GRPCauth "github.com/Weit145/proto-repo/auth"
	"google.golang.org/grpc"
)

type server struct {
	GRPCauth.UnimplementedAuthServer
	grpcServer *grpc.Server
	listener   net.Listener
}

var (
	port = flag.Int("port", 50051, "The server port")
)

func New() (*server, error) {
	flag.Parse()

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		return nil, fmt.Errorf("failed to listen: %w", err)
	}

	s := grpc.NewServer()
	srv := &server{
		grpcServer: s,
		listener:   lis,
	}

	GRPCauth.RegisterAuthServer(s, srv)

	return srv, nil
}

func (s *server) Start() error {
	log.Printf("server listening at %v", s.listener.Addr())
	if err := s.grpcServer.Serve(s.listener); err != nil {
		return fmt.Errorf("failed to serve: %w", err)
	}
	return nil
}

func (s *server) Stop() {
	s.grpcServer.GracefulStop()
}
