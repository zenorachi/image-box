package repository

import (
	"context"
	"database/sql"

	"github.com/zenorachi/image-box/model"
)

type Users struct {
	db *sql.DB
}

func NewUsers(db *sql.DB) *Users {
	return &Users{db: db}
}

func (u *Users) Create(ctx context.Context, user model.User) error {
	tx, err := u.db.BeginTx(ctx, &sql.TxOptions{
		Isolation: sql.LevelSerializable,
		ReadOnly:  false,
	})
	if err != nil {
		return err
	}

	_, err = tx.ExecContext(ctx, "INSERT INTO users (login, email, password, registered_at) "+
		"VALUES ($1, $2, $3, $4)",
		user.Login, user.Email, user.Password, user.RegisteredAt)
	if err != nil {
		_ = tx.Rollback()
		return err
	}

	return tx.Commit()
}

func (u *Users) GetByCredentials(ctx context.Context, login, password string) (model.User, error) {
	var user model.User

	tx, err := u.db.BeginTx(ctx, &sql.TxOptions{
		Isolation: sql.LevelSerializable,
		ReadOnly:  true,
	})
	if err != nil {
		return model.User{}, err
	}

	err = tx.QueryRowContext(ctx, "SELECT id, login, email, password, registered_at FROM users "+
		"WHERE login = $1 AND password = $2", login, password).
		Scan(&user.ID, &user.Login, &user.Email, &user.Password, &user.RegisteredAt)
	if err != nil {
		_ = tx.Rollback()
		return model.User{}, err
	}

	return user, tx.Commit()
}
