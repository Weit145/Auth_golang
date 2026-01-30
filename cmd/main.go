package main

import (
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/Weit145/Auth_golang/internal/config"
	"github.com/Weit145/Auth_golang/internal/grpc/gateway"
	"github.com/Weit145/Auth_golang/internal/lib/logger"
	"github.com/Weit145/Auth_golang/internal/service"
	// Import the registration service
)

func main() {
	//Init config
	cfg := config.MustLoad()

	//Init logger
	log := setupLogger(cfg.Env)
	log.Info("Start AUTH")

	// Init registration service
	Service := service.New(log)

	//Init grpc
	grpcServer, err := gateway.New(log, Service, cfg.GRPC.Address) // Pass the service
	if err != nil {
		log.Error("cannot create server", logger.Err(err))
		os.Exit(1)
	}

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop

	log.Info("Shutting down gRPC server...")
	grpcServer.GracefulStop()

}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case "local":
		log = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case "prod":
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	default:
		log = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	}

	return log
}
