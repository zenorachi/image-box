package model

import "time"

type RefreshToken struct {
	ID        uint
	UserID    uint
	Token     string
	ExpiresAt time.Time
}

func CreateToken(userID uint, token string, expiresAt time.Time) RefreshToken {
	return RefreshToken{
		UserID:    userID,
		Token:     token,
		ExpiresAt: expiresAt,
	}
}
