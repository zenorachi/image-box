package models

import (
	"github.com/go-playground/validator/v10"
	"time"
)

var validate *validator.Validate

func init() {
	validate = validator.New()
}

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

type SignUpInput struct {
	Login    string `json:"login" validate:"required,gte=6"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,gte=6"`
}

func (i SignUpInput) Validate() error {
	return validate.Struct(i)
}

type SignInInput struct {
	Login    string `json:"login" validate:"required,gte=6"`
	Password string `json:"password" validate:"required,gte=6"`
}

func (i SignInInput) Validate() error {
	return validate.Struct(i)
}
