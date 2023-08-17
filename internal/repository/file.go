package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/zenorachi/image-box/internal/service"
	"github.com/zenorachi/image-box/model"
)

type Files struct {
	db *sql.DB
}

func NewFiles(db *sql.DB) *Files {
	return &Files{db: db}
}

func (f *Files) Create(ctx context.Context, file model.File) error {
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
		_ = tx.Rollback()
	}

	return tx.Commit()
}

func (f *Files) Get(ctx context.Context, userID uint) ([]model.File, error) {
	var files []model.File

	tx, err := f.db.BeginTx(ctx, &sql.TxOptions{
		Isolation: sql.LevelSerializable,
		ReadOnly:  true,
	})
	if err != nil {
		return []model.File{}, err
	}

	rows, err := tx.QueryContext(ctx, "SELECT id, user_id, name, url, size, uploaded_at "+
		"FROM files "+
		"WHERE user_id = $1", userID)
	if err != nil {
		_ = tx.Rollback()
		if errors.Is(err, sql.ErrNoRows) {
			return nil, service.FilesNotFound
		}
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var file model.File
		err = rows.Scan(&file.ID, &file.UserID, &file.Name, &file.URL, &file.Size, &file.UploadedAt)
		if err != nil {
			return nil, err
		}
		files = append(files, file)
	}

	return files, tx.Commit()
}
