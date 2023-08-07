package main

import (
	"context"
	"errors"
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

	endpoint := "127.0.0.1:9000"
	minioRootUser := "root"
	minioRootPassword := "password"

	minioClient, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(minioRootUser, minioRootPassword, ""),
		Secure: false,
	})
	if err != nil {
		log.Fatalln(err)
	}

	hasher := hash.NewSHA1Hasher("testLol")
	provider := storage.NewProvider(minioClient, "kek", endpoint)
	policy := `{"Version": "2012-10-17","Statement": [{"Action": ["s3:GetObject"],"Effect": "Allow","Principal": {"AWS": ["*"]},"Resource": ["arn:aws:s3:::kek/*"],"Sid": ""}]}`
	err = minioClient.SetBucketPolicy(context.Background(), "kek", policy)
	if err != nil {
		log.Fatalln(err)
	}

	usersRepo := repository.NewUsers(db)
	tokenRepo := repository.NewTokens(db)
	filesRepo := repository.NewFiles(db)

	users := service.NewUsers(hasher, usersRepo, tokenRepo, []byte("kekSecret"), cfg.Auth.TokenTTL, cfg.Auth.RefreshTTL)
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
