package service

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/zenorachi/image-box/model"
)

var UserNotFound = errors.New("user not found")

type (
	PasswordHasher interface {
		Hash(password string) (string, error)
	}

	UserRepository interface {
		Create(ctx context.Context, user model.User) error
		GetByCredentials(ctx context.Context, login, password string) (model.User, error)
	}

	TokenRepository interface {
		Create(ctx context.Context, token model.RefreshToken) error
		Get(ctx context.Context, token string) (model.RefreshToken, error)
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

func (u *Users) SignUp(ctx context.Context, input model.SignUpInput) error {
	password, err := u.hasher.Hash(input.Password)
	if err != nil {
		return err
	}

	user := model.CreateUser(input.Login, input.Email, password)
	return u.repository.Create(ctx, user)
}

func (u *Users) SignIn(ctx context.Context, input model.SignInInput) (string, string, error) {
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
