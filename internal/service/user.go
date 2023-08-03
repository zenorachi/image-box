package service

import (
	"database/sql"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/zenorachi/image-box/models"
	"strconv"
	"time"
)

var UserNotFound = errors.New("user not found")

type (
	PasswordHasher interface {
		Hash(password string) (string, error)
	}

	UserRepository interface {
		Create(ctx *gin.Context, user models.User) error
		GetByCredentials(ctx *gin.Context, login, password string) (models.User, error)
	}
)

type Users struct {
	hasher     PasswordHasher
	repository UserRepository
	secret     []byte
	ttl        time.Duration
}

func NewUsers(hasher PasswordHasher, repository UserRepository, secret []byte, tokenTTL time.Duration) *Users {
	return &Users{
		hasher:     hasher,
		repository: repository,
		secret:     secret,
		ttl:        tokenTTL,
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

func (u *Users) SignIn(ctx *gin.Context, input models.SignInInput) (string, error) {
	password, err := u.hasher.Hash(input.Password)
	if err != nil {
		return "", err
	}

	user, err := u.repository.GetByCredentials(ctx, input.Login, password)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", UserNotFound
		}
		return "", err
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		IssuedAt:  time.Now().Unix(),
		ExpiresAt: time.Now().Add(u.ttl).Unix(), //todo : config time ttl
		Subject:   strconv.Itoa(int(user.ID)),
	})

	return token.SignedString(u.secret)
}
