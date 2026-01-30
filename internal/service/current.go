package service

import (
	"context"
	"log/slog"
)

type User struct {
	Id         int
	Login      string
	IsActive   bool
	IsVerified bool
	Role       string
}

func (s *Service) Current(ctx context.Context, AssetToken string) (*User, error) {
	s.log.Info("Refresh method called", slog.String("AssetToken: ", AssetToken))
	resp := User{
		Id:         1,
		Login:      "Weit",
		IsActive:   true,
		IsVerified: true,
		Role:       "user",
	}
	return &resp, nil
}
