package service

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/zenorachi/image-box/models"
	"math/rand"
	"strconv"
	"time"
)

func (u *Users) ParseToken(ctx *gin.Context, token string) (uint, error) {
	tokenParser := func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signign method: %v", token.Header["alg"])
		}

		return u.secret, nil
	}

	t, err := jwt.Parse(token, tokenParser)
	if err != nil {
		return 0, err
	}

	if !t.Valid {
		return 0, errors.New("invalid token")
	}

	claims, ok := t.Claims.(jwt.MapClaims)
	if !ok {
		return 0, errors.New("invalid claims")
	}

	id, err := strconv.Atoi(claims["sub"].(string))
	if err != nil {
		return 0, errors.New("invalid subject")
	}

	return uint(id), nil
}

func (u *Users) generateTokens(ctx *gin.Context, userID uint) (string, string, error) {
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		IssuedAt:  time.Now().Unix(),
		ExpiresAt: time.Now().Add(u.ttl).Unix(),
		Subject:   strconv.Itoa(int(userID)),
	})

	accessToken, err := t.SignedString(u.secret)
	if err != nil {
		return "", "", err
	}

	refreshToken, err := newRefreshToken()
	if err != nil {
		return "", "", err
	}

	rT := models.CreateToken(userID, refreshToken, time.Now().Add(u.refreshTTL))
	if err = u.tokenRepo.Create(ctx, rT); err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}

func (u *Users) RefreshTokens(ctx *gin.Context, refreshToken string) (string, string, error) {
	token, err := u.tokenRepo.Get(ctx, refreshToken)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", "", UserNotFound
		}
		return "", "", err
	}

	if token.ExpiresAt.Unix() < time.Now().Unix() {
		return "", "", errors.New("token expired") //todo - const
	}

	return u.generateTokens(ctx, token.UserID)
}

func newRefreshToken() (string, error) {
	b := make([]byte, 32)

	src := rand.NewSource(time.Now().Unix())
	random := rand.New(src)

	if _, err := random.Read(b); err != nil {
		return "", err
	}

	return fmt.Sprintf("%x", b), nil
}
