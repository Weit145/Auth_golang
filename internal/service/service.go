package service

import (
	"log/slog"

	"github.com/Weit145/Auth_golang/internal/storage/postgresql"
)

type Service struct {
	db  *postgresql.Storage
	log *slog.Logger
}

func New(log *slog.Logger, db *postgresql.Storage) *Service {
	return &Service{log: log, db: db}
}
