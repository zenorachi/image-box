package repository

import (
	"errors"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/zenorachi/image-box/model"
)

func TestTokens_Get(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("error creating database connection: %v\n", err)
	}
	defer db.Close()

	repo := NewTokens(db)

	type args struct {
		token string
	}
	type mockBehaviour func(args args)

	tests := []struct {
		name          string
		mockBehaviour mockBehaviour
		args          args
		wantErr       bool
	}{
		{
			name: "OK",
			args: args{
				token: "token",
			},
			mockBehaviour: func(args args) {
				mock.ExpectBegin()

				rows := sqlmock.NewRows([]string{"id", "user_id", "token", "expires_at"}).
					AddRow(1, "2", "token", time.Now().Round(time.Second))

				expectedQuery := "SELECT id, user_id, token, expires_at FROM refresh_tokens WHERE token = $1"
				mock.ExpectQuery(regexp.QuoteMeta(expectedQuery)).WithArgs(args.token).WillReturnRows(rows)

				expectedExec := "DELETE FROM refresh_tokens WHERE token = $1"
				mock.ExpectExec(regexp.QuoteMeta(expectedExec)).WithArgs(args.token).
					WillReturnResult(sqlmock.NewResult(1, 1))

				mock.ExpectCommit()
			},
		},
		{
			name: "ERROR",
			args: args{},
			mockBehaviour: func(args args) {
				mock.ExpectBegin()

				expectedQuery := "SELECT id, user_id, token, expires_at FROM refresh_tokens WHERE token = $1"
				mock.ExpectQuery(regexp.QuoteMeta(expectedQuery)).WithArgs(args.token).WillReturnError(errors.New("test error"))

				expectedExec := "DELETE FROM refresh_tokens WHERE token = $1"
				mock.ExpectExec(regexp.QuoteMeta(expectedExec)).WithArgs(args.token).
					WillReturnResult(sqlmock.NewResult(1, 1))

				mock.ExpectRollback()
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockBehaviour(tt.args)
			token, err := repo.Get(&gin.Context{}, tt.args.token)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.args.token, token.Token)
			}
		})
	}
}

func TestTokens_Create(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("error creating database connection: %v\n", err)
	}
	defer db.Close()

	repo := NewTokens(db)

	type args struct {
		token model.RefreshToken
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
				token: model.RefreshToken{
					ID:        1,
					UserID:    2,
					Token:     "token",
					ExpiresAt: time.Now().Round(time.Second),
				},
			},
			mockBehaviour: func(args args) {
				mock.ExpectBegin()

				expectedExec := "INSERT INTO refresh_tokens (user_id, token, expires_at) VALUES ($1, $2, $3)"
				mock.ExpectExec(regexp.QuoteMeta(expectedExec)).
					WithArgs(args.token.UserID, args.token.Token, args.token.ExpiresAt).
					WillReturnResult(sqlmock.NewResult(1, 1))

				mock.ExpectCommit()
			},
		},
		{
			name: "ERROR",
			args: args{},
			mockBehaviour: func(args args) {
				mock.ExpectBegin()

				expectedExec := "INSERT INTO refresh_tokens (user_id, token, expires_at) VALUES ($1, $2, $3)"
				mock.ExpectExec(regexp.QuoteMeta(expectedExec)).
					WithArgs(args.token.UserID, args.token.Token, args.token.ExpiresAt).
					WillReturnError(errors.New("test error"))

				mock.ExpectRollback()
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockBehaviour(tt.args)
			err := repo.Create(&gin.Context{}, tt.args.token)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
