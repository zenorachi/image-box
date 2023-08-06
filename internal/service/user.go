package service

import (
	"database/sql"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/zenorachi/image-box/models"
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

	TokenRepository interface {
		Create(ctx *gin.Context, token models.RefreshToken) error
		Get(ctx *gin.Context, token string) (models.RefreshToken, error)
	}
)

type Users struct {
	hasher     PasswordHasher
	repository UserRepository
	tokenRepo  TokenRepository
	secret     []byte
	ttl        time.Duration
	refreshTTL time.Duration
}

func NewUsers(hasher PasswordHasher, repository UserRepository, tokenRepo TokenRepository, secret []byte, tokenTTL time.Duration, refreshTTL time.Duration) *Users {
	return &Users{
		hasher:     hasher,
		repository: repository,
		tokenRepo:  tokenRepo,
		secret:     secret,
		ttl:        tokenTTL,
		refreshTTL: refreshTTL,
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

func (u *Users) SignIn(ctx *gin.Context, input models.SignInInput) (string, string, error) {
	password, err := u.hasher.Hash(input.Password)
	if err != nil {
		return "", "", err
	}

	user, err := u.repository.GetByCredentials(ctx, input.Login, password)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", "", UserNotFound
		}
		return "", "", err
	}

	return u.generateTokens(ctx, user.ID)
}
