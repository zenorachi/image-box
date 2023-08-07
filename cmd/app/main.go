package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/zenorachi/image-box/internal/config"
	"github.com/zenorachi/image-box/internal/repository"
	"github.com/zenorachi/image-box/internal/service"
	"github.com/zenorachi/image-box/internal/transport/rest"
	"github.com/zenorachi/image-box/pkg/database/postgres"
	"github.com/zenorachi/image-box/pkg/hash"
	"github.com/zenorachi/image-box/pkg/storage"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

const (
	ConfigDir  = "configs"
	ConfigFile = "main"
	ENVFile    = ".env"
)

func main() {
	gin.SetMode(gin.ReleaseMode)

	err := config.InitENV(ENVFile)
	if err != nil {
		log.Fatalln(err)
	}
	cfg, err := config.New(ConfigDir, ConfigFile)
	if err != nil {
		log.Fatalln(err)
	}

	db, err := postgres.NewDB(cfg.DB)
	if err != nil {
		log.Fatalln(err)
	}
	defer db.Close()

	minioRootUser := os.Getenv("MINIO_ROOT_USER")
	minioRootPassword := os.Getenv("MINIO_ROOT_PASSWORD")
	fmt.Println(minioRootPassword, minioRootUser)
	fmt.Println(cfg.Minio)

	minioClient, err := minio.New(cfg.Minio.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(minioRootUser, minioRootPassword, ""),
		Secure: false,
	})
	if err != nil {
		log.Fatalln(err)
	}
	if isExist, _ := minioClient.BucketExists(context.Background(), cfg.Minio.Bucket); !isExist {
		err = minioClient.MakeBucket(context.Background(), cfg.Minio.Bucket, minio.MakeBucketOptions{
			Region: "eu-central-1",
		})
		if err != nil {
			log.Fatalln(err)
		}
	}

	hasher := hash.NewSHA1Hasher(cfg.Hash.Salt)
	provider := storage.NewProvider(minioClient, cfg.Minio.Bucket, cfg.Minio.Endpoint)
	policy := `{"Version":"2012-10-17","Statement":[{"Effect":"Allow","Principal":{"AWS":["*"]},"Action":["s3:GetObject"],"Resource":["arn:aws:s3:::` + cfg.Minio.Bucket + `/*"]}]}`
	err = minioClient.SetBucketPolicy(context.Background(), cfg.Minio.Bucket, policy)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println(cfg.Hash)

	usersRepo := repository.NewUsers(db)
	tokenRepo := repository.NewTokens(db)
	filesRepo := repository.NewFiles(db)

	users := service.NewUsers(hasher, usersRepo, tokenRepo, []byte(cfg.Hash.Secret), cfg.Auth.TokenTTL, cfg.Auth.RefreshTTL)
	files := service.NewFiles(filesRepo, provider)

	handler := rest.NewHandler(users, files)

	s := rest.NewServer(handler, cfg.Server.Host, cfg.Server.Port)

	go func() {
		log.Println("Server started")
		if err := s.Run(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalln("Error starting server:", err)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	// Waiting for CTRL+C
	<-stop
	log.Println("Shutdown gracefully...")
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if err := s.Shutdown(ctx); err != nil {
		log.Printf("Error shutting down server: %s\n", err)
	}
	log.Println("Server gracefully stopped")
}
