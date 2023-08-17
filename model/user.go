package model

import (
	"time"
)

type User struct {
	ID           uint      `json:"id"`
	Login        string    `json:"login"`
	Email        string    `json:"email"`
	Password     string    `json:"password"`
	RegisteredAt time.Time `json:"registered_at"`
}

func CreateUser(login, email, password string) User {
	return User{
		Login:        login,
		Email:        email,
		Password:     password,
		RegisteredAt: time.Now(),
	}
}
