package models

import "time"

type File struct {
	ID         uint
	UserID     uint
	Name       string
	URL        string
	Size       int64
	UploadedAt time.Time
}

func CreateFile(userID uint, name, url string, size int64, uploadedAt time.Time) File {
	return File{
		UserID:     userID,
		Name:       name,
		URL:        url,
		Size:       size,
		UploadedAt: uploadedAt,
	}
}
