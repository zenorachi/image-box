package repository

import (
	"database/sql"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/zenorachi/image-box/internal/service"
	"github.com/zenorachi/image-box/models"
)

type Files struct {
	db *sql.DB
}

func NewFiles(db *sql.DB) *Files {
	return &Files{db: db}
}

func (f *Files) Create(ctx *gin.Context, file models.File) error {
	_, err := f.db.Exec("INSERT INTO files (user_id, name, url, size, uploaded_at) "+
		"VALUES ($1, $2, $3, $4, $5)",
		file.UserID, file.Name, file.URL, file.Size, file.UploadedAt)

	return err
}

func (f *Files) Get(ctx *gin.Context, userID uint) ([]models.File, error) {
	var files []models.File
	rows, err := f.db.Query("SELECT id, user_id, name, url, size, uploaded_at "+
		"FROM files "+
		"WHERE user_id = $1", userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, service.FilesNotFound
		}
		return nil, err
	}

	for rows.Next() {
		var file models.File
		err = rows.Scan(&file.ID, &file.UserID, &file.Name, &file.URL, &file.Size, &file.UploadedAt)
		if err != nil {
			return nil, err
		}
		files = append(files, file)
	}

	return files, nil
}
