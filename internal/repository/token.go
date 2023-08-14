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
	return &Tokens{db: db}
}

func (t *Tokens) Create(ctx *gin.Context, token models.RefreshToken) error {
	tx, err := t.db.BeginTx(ctx, &sql.TxOptions{
		Isolation: sql.LevelSerializable,
		ReadOnly:  false,
	})
	if err != nil {
		return err
	}

	_, err = t.db.ExecContext(ctx, "INSERT INTO refresh_tokens (user_id, token, expires_at) "+
		"VALUES ($1, $2, $3)", token.UserID, token.Token, token.ExpiresAt)
	if err != nil {
		_ = tx.Rollback()
		return err
	}

	return tx.Commit()
}

func (t *Tokens) Get(ctx *gin.Context, token string) (models.RefreshToken, error) {
	var refreshToken models.RefreshToken

	tx, err := t.db.BeginTx(ctx, &sql.TxOptions{
		Isolation: sql.LevelSerializable,
		ReadOnly:  false,
	})
	if err != nil {
		return models.RefreshToken{}, err
	}

	err = t.db.QueryRowContext(ctx, "SELECT id, user_id, token, expires_at FROM refresh_tokens "+
		"WHERE token = $1", token).
		Scan(&refreshToken.ID, &refreshToken.UserID, &refreshToken.Token, &refreshToken.ExpiresAt)
	if err != nil {
		_ = tx.Rollback()
		return models.RefreshToken{}, err
	}

	_, err = t.db.Exec("DELETE FROM refresh_tokens WHERE token = $1", token)
	if err != nil {
		_ = tx.Rollback()
		return models.RefreshToken{}, err
	}

	return refreshToken, tx.Commit()
}
