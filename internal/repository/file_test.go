package repository

import (
	"errors"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/zenorachi/image-box/models"
	"regexp"
	"testing"
	"time"
)

func TestFiles_Get(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("error creating database connection: %v\n", err)
	}
	defer db.Close()

	repo := NewFiles(db)

	type args struct {
		userID uint
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
				userID: 3,
			},
			mockBehaviour: func(args args) {
				mock.ExpectBegin()

				rows := sqlmock.NewRows([]string{"id", "user_id", "name", "url", "size", "uploaded_at"}).
					AddRow(1, 3, "file", "file.com", 134, time.Now().Round(time.Second))

				expectedQuery := "SELECT id, user_id, name, url, size, uploaded_at " +
					"FROM files " +
					"WHERE user_id = $1"
				mock.ExpectQuery(regexp.QuoteMeta(expectedQuery)).
					WithArgs(args.userID).WillReturnRows(rows)

				mock.ExpectCommit()
			},
		},
		{
			name: "ERROR",
			args: args{},
			mockBehaviour: func(args args) {
				mock.ExpectBegin()

				expectedQuery := "SELECT id, user_id, name, url, size, uploaded_at " +
					"FROM files " +
					"WHERE user_id = $1"
				mock.ExpectQuery(regexp.QuoteMeta(expectedQuery)).
					WithArgs(args.userID).WillReturnError(errors.New("test error"))

				mock.ExpectRollback()
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockBehaviour(tt.args)
			file, err := repo.Get(&gin.Context{}, tt.args.userID)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.args.userID, file[0].UserID)
			}
		})
	}
}

func TestFiles_Create(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("error creating database connection: %v\n", err)
	}
	defer db.Close()

	repo := NewFiles(db)

	type args struct {
		file models.File
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
				models.File{
					ID:         1,
					UserID:     3,
					Name:       "file",
					URL:        "file.com",
					Size:       143,
					UploadedAt: time.Now().Round(time.Second),
				},
			},
			mockBehaviour: func(args args) {
				mock.ExpectBegin()

				expectedExec := "INSERT INTO files (user_id, name, url, size, uploaded_at) " +
					"VALUES ($1, $2, $3, $4, $5)"
				mock.ExpectExec(regexp.QuoteMeta(expectedExec)).
					WithArgs(args.file.UserID, args.file.Name, args.file.URL, args.file.Size, args.file.UploadedAt).
					WillReturnResult(sqlmock.NewResult(1, 1))

				mock.ExpectCommit()
			},
		},
		{
			name: "ERROR",
			args: args{},
			mockBehaviour: func(args args) {
				mock.ExpectBegin()

				expectedExec := "INSERT INTO files (user_id, name, url, size, uploaded_at) " +
					"VALUES ($1, $2, $3, $4, $5)"
				mock.ExpectExec(regexp.QuoteMeta(expectedExec)).
					WithArgs(args.file.UserID, args.file.Name, args.file.URL, args.file.Size, args.file.UploadedAt).
					WillReturnError(errors.New("test error"))

				mock.ExpectRollback()
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockBehaviour(tt.args)
			err := repo.Create(&gin.Context{}, tt.args.file)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
