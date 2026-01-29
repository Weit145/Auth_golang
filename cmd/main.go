package main

import (
	"log/slog"
	"os"

	"github.com/Weit145/Auth_golang/internal/config"
	"github.com/Weit145/Auth_golang/internal/grpc/gateway"
	"github.com/Weit145/Auth_golang/internal/lib/logger"
)

func main() {
	//Init config
	cfg := config.MustLoad()

	//Init logger
	log := setupLogger(cfg.Env)
	log.Info("Start AUTH")

	//Init grpc
	srv, err := gateway.New()
	if err != nil {
		log.Error("cannot create server: %v", logger.Err(err))
	}

	go func() {
		if err := srv.Start(); err != nil {
			log.Error("server failed: %v", logger.Err(err))
		}
	}()
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
