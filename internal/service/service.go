package service

import (
	"log/slog"

	"github.com/Weit145/Auth_golang/internal/config"
	"github.com/Weit145/Auth_golang/internal/storage"
)

type Service struct {
	storage storage.Storage
	log     *slog.Logger
	cfg     *config.Config
}

func New(log *slog.Logger, storage storage.Storage, cfg *config.Config) *Service {
	return &Service{log: log, storage: storage, cfg: cfg}
}
