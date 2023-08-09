package repository

import (
	"fmt"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/zenorachi/image-box/models"
	"regexp"
	"testing"
	"time"
)

func TestUsers_Get(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("error creating database connection: %v\n", err)
	}
	defer db.Close()

	repo := NewUsers(db)

	type args struct {
		login    string
		password string
	}
	type mockBehaviour func(args args)

	tests := []struct {
		name          string
		mockBehaviour mockBehaviour
		args          args
		user          models.User
		wantErr       bool
	}{
		{
			name: "OK",
			args: args{
				login:    "login",
				password: "password",
			},
			user: models.User{
				ID:           1,
				Login:        "login",
				Email:        "user-email",
				Password:     "password",
				RegisteredAt: time.Now().Round(time.Second),
			},
			mockBehaviour: func(args args) {
				mock.ExpectBegin()

				rows := sqlmock.NewRows([]string{"id", "login", "email", "password", "registered_at"}).
					AddRow(1, "login", "user-email", "password", time.Now().Round(time.Second))

				expectedQuery := "SELECT id, login, email, password, registered_at FROM users WHERE login = $1 AND password = $2"
				mock.ExpectQuery(regexp.QuoteMeta(expectedQuery)).WithArgs(args.login, args.password).WillReturnRows(rows)

				mock.ExpectCommit()
			},
		},
		{
			name: "ERROR",
			args: args{
				login:    "hello",
				password: "world",
			},
			user: models.User{},
			mockBehaviour: func(args args) {
				mock.ExpectBegin()

				rows := sqlmock.NewRows([]string{"id", "login", "email", "password", "registered_at"})
				expectedQuery := "SELECT id, login, email, password, registered_at FROM users WHERE login = $1 AND password = $2"
				mock.ExpectQuery(regexp.QuoteMeta(expectedQuery)).WithArgs(args.login, args.password).WillReturnRows(rows)

				mock.ExpectRollback()
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockBehaviour(tt.args)
			user, err := repo.GetByCredentials(&gin.Context{}, tt.args.login, tt.args.password)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.user, user)
			}
		})
	}
}

func TestUsers_Create(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("error creating database connection: %v\n", err)
	}
	defer db.Close()

	repo := NewUsers(db)

	type args struct {
		user models.User
	}
	type mockBehaviour func(args args)

	tests := []struct {
		name          string
		args          args
		mockBehaviour mockBehaviour
		wantErr       bool
	}{
		{
			name: "OK",
			args: args{
				user: models.User{
					ID:           1,
					Login:        "user",
					Email:        "email@go.dev",
					Password:     "password",
					RegisteredAt: time.Now().Round(time.Second),
				},
			},
			mockBehaviour: func(args args) {
				mock.ExpectBegin()

				expectedExec := "INSERT INTO users (login, email, password, registered_at) VALUES ($1, $2, $3, $4)"
				mock.ExpectExec(regexp.QuoteMeta(expectedExec)).
					WithArgs(args.user.Login, args.user.Email, args.user.Password, args.user.RegisteredAt).
					WillReturnResult(sqlmock.NewResult(1, 1))

				mock.ExpectCommit()
			},
		},
		{
			name: "ERROR",
			args: args{
				user: models.User{},
			},
			mockBehaviour: func(args args) {
				mock.ExpectBegin()

				expectedExec := "INSERT INTO users (login, email, password, registered_at) VALUES ($1, $2, $3, $4)"
				mock.ExpectExec(regexp.QuoteMeta(expectedExec)).
					WithArgs(args.user.Login, args.user.Email, args.user.Password, args.user.RegisteredAt).
					WillReturnError(fmt.Errorf("test error"))

				mock.ExpectRollback()
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockBehaviour(tt.args)
			err := repo.Create(&gin.Context{}, tt.args.user)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
