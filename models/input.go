package models

import "github.com/go-playground/validator/v10"

var validate *validator.Validate

func init() {
	validate = validator.New()
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
