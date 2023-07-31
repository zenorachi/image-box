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
	_, err := u.db.Exec("INSERT INTO users (email, password, registered_at) "+
		"VALUES ($1, $2, $3)",
		user.Email, user.Password, user.RegisteredAt)

	return err
}

func (u *Users) GetByCredentials(ctx *gin.Context, email, password string) (error, models.User) {
	//TODO
	return nil, models.User{}
}
