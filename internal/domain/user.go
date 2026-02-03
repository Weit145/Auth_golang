package domain

type User struct {
	Id               int64
	Login            string
	Email            string
	PasswordHash     string
	IsActive         bool
	IsVerified       bool
	Role             string
	RefreshTokenHash string
}
