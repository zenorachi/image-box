package storage

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/minio/minio-go/v7"
)

type FileStorage struct {
	client   *minio.Client
	bucket   string
	endpoint string
}

func NewProvider(client *minio.Client, bucket, endpoint string) *FileStorage {
	return &FileStorage{
		client:   client,
		bucket:   bucket,
		endpoint: endpoint,
	}
}

func (fs *FileStorage) Upload(ctx *gin.Context, input UploadInput) (string, error) {
	opts := minio.PutObjectOptions{
		ContentType:  input.ContentType,
		UserMetadata: map[string]string{"x-amz-acl": "public-read"},
	}

	if _, err := fs.client.PutObject(ctx, fs.bucket, input.Name, input.File, input.Size, opts); err != nil {
		return "", err
	}

	return fs.generateFileURL(input.Name), nil
}

func (fs *FileStorage) generateFileURL(filename string) string {
	return fmt.Sprintf("http://%s/%s/%s", fs.endpoint, fs.bucket, filename)
}
