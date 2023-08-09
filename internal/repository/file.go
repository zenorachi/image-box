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
	tx, err := f.db.BeginTx(ctx, &sql.TxOptions{
		Isolation: sql.LevelSerializable,
		ReadOnly:  false,
	})
	if err != nil {
		return err
	}

	_, err = tx.ExecContext(ctx, "INSERT INTO files (user_id, name, url, size, uploaded_at) "+
		"VALUES ($1, $2, $3, $4, $5)",
		file.UserID, file.Name, file.URL, file.Size, file.UploadedAt)
	if err != nil {
		tx.Rollback()
	}

	return tx.Commit()
}

func (f *Files) Get(ctx *gin.Context, userID uint) ([]models.File, error) {
	var files []models.File

	tx, err := f.db.BeginTx(ctx, &sql.TxOptions{
		Isolation: sql.LevelSerializable,
		ReadOnly:  true,
	})
	if err != nil {
		return []models.File{}, err
	}

	rows, err := tx.QueryContext(ctx, "SELECT id, user_id, name, url, size, uploaded_at "+
		"FROM files "+
		"WHERE user_id = $1", userID)
	if err != nil {
		tx.Rollback()
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

	return files, tx.Commit()
}
