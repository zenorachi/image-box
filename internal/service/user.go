package service

import (
	"github.com/gin-gonic/gin"
	"github.com/zenorachi/image-box/models"
)

type PasswordHasher interface {
	Hash(password string) (string, error)
}

type UserRepository interface {
	Create(ctx *gin.Context, user models.User) error
	GetByCredentials(ctx *gin.Context, email, password string) (error, models.User)
}

type Users struct {
	hasher     PasswordHasher
	repository UserRepository
	secret     []byte
}

func NewUsers(hasher PasswordHasher, repository UserRepository, secret []byte) *Users {
	return &Users{
		hasher:     hasher,
		repository: repository,
		secret:     secret,
	}
}

func (u *Users) SignUp(ctx *gin.Context, input models.SignUpInput) error {
	password, err := u.hasher.Hash(input.Password)
	if err != nil {
		return err
	}

	user := models.CreateUser(input.Login, input.Email, password)
	return u.repository.Create(ctx, user)
}
