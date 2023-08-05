package repository

import (
	"database/sql"
	"github.com/gin-gonic/gin"
	"github.com/zenorachi/image-box/models"
)

type Tokens struct {
	db *sql.DB
}

func NewTokens(db *sql.DB) *Tokens {
	return &Tokens{
		db: db,
	}
}

func (t *Tokens) Create(ctx *gin.Context, token models.RefreshToken) error {
	_, err := t.db.Exec("INSERT INTO refresh_tokens (user_id, token, expires_at) "+
		"VALUES ($1, $2, $3)", token.UserID, token.Token, token.ExpiresAt)

	return err
}

func (t *Tokens) Get(ctx *gin.Context, token string) (models.RefreshToken, error) {
	var refreshToken models.RefreshToken
	err := t.db.QueryRow("SELECT id, user_id, token, expires_at FROM refresh_tokens "+
		"WHERE token = $1", token).
		Scan(&refreshToken.ID, &refreshToken.UserID, &refreshToken.Token, &refreshToken.ExpiresAt)

	if err != nil {
		return models.RefreshToken{}, err
	}

	_, err = t.db.Exec("DELETE FROM refresh_tokens WHERE user_id = $1", refreshToken.UserID)

	return refreshToken, err
}
