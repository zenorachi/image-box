package repository

import (
	"database/sql"
	"github.com/gin-gonic/gin"
	"github.com/zenorachi/image-box/models"
)

type Users struct {
	db *sql.DB
}

func NewUsers(db *sql.DB) *Users {
	return &Users{db: db}
}

func (u *Users) Create(ctx *gin.Context, user models.User) error {
	_, err := u.db.ExecContext(ctx, "INSERT INTO users (login, email, password, registered_at) "+
		"VALUES ($1, $2, $3, $4)",
		user.Login, user.Email, user.Password, user.RegisteredAt)

	return err
}

func (u *Users) GetByCredentials(ctx *gin.Context, login, password string) (models.User, error) {
	var user models.User
	err := u.db.QueryRowContext(ctx, "SELECT id, login, email, password, registered_at FROM users "+
		"WHERE login = $1 AND password = $2", login, password).
		Scan(&user.ID, &user.Login, &user.Email, &user.Password, &user.RegisteredAt)

	return user, err
}
