package repository

import (
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/zenorachi/image-box/models"
	"regexp"
	"testing"
	"time"
)

func TestGetUser(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("error creating database connection: %v", err)
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
				rows := sqlmock.NewRows([]string{"id", "login", "email", "password", "registered_at"}).
					AddRow(1, "login", "user-email", "password", time.Now().Round(time.Second))

				expectedQuery := "SELECT id, login, email, password, registered_at FROM users WHERE login = $1 AND password = $2"
				mock.ExpectQuery(regexp.QuoteMeta(expectedQuery)).WithArgs(args.login, args.password).WillReturnRows(rows)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockBehaviour(tt.args)
			user, err := repo.GetByCredentials(nil, tt.args.login, tt.args.password)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.user, user)
			}
		})
	}

}
