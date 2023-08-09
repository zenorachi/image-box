package app

import (
	"context"
	"database/sql"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/sirupsen/logrus"
	"github.com/zenorachi/image-box/internal/config"
	"github.com/zenorachi/image-box/internal/repository"
	"github.com/zenorachi/image-box/internal/service"
	"github.com/zenorachi/image-box/internal/transport/rest"
	"github.com/zenorachi/image-box/pkg/database/postgres"
	"github.com/zenorachi/image-box/pkg/hash"
	"github.com/zenorachi/image-box/pkg/storage"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

// @title           ImageBox API
// @version         1.0
// @description     zenorachi's ImageBox server.

// @host      localhost:8080
// @BasePath  /

// @securityDefinitions.apikey UsersAuth
// @in header
// @name Authorization

func init() {
	// Gin Setup
	gin.SetMode(gin.ReleaseMode)

	// Logrus setup
	logrus.SetFormatter(&logrus.TextFormatter{
		ForceColors:   true,
		FullTimestamp: true,
	})
}

func Run(configDir, configFile, envFile string) error {
	// Setup GIN's mode
	gin.SetMode(gin.ReleaseMode)

	// CONFIG
	cfg, err := initConfig(configDir, configFile, envFile)
	if err != nil {
		return err
	}

	// DB
	db, err := postgres.NewDB(cfg.DB)
	if err != nil {
		return err
	}
	defer db.Close()

	// MINIO CLIENT
	minioClient, err := initMinioClient(cfg)
	if err != nil {
		return err
	}

	// CREATE BUCKET
	if isExist, _ := minioClient.BucketExists(context.Background(), cfg.Minio.Bucket); !isExist {
		err = minioClient.MakeBucket(context.Background(), cfg.Minio.Bucket, minio.MakeBucketOptions{
			Region: "eu-central-1",
		})
		if err != nil {
			return err
		}
	}

	// FILE'S STORAGE
	provider, err := initStorage(minioClient, cfg)
	if err != nil {
		return err
	}

	// HANDLER
	handler := initHandler(cfg, provider, db)
	server := rest.NewServer(handler, cfg.Server.Host, cfg.Server.Port)
	go func() {
		if err := server.Run(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logrus.Fatalln("unable to start server:", err)
		}
	}()
	logrus.Infoln("server is running")

	// GRACEFUL SHUTDOWN
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop
	logrus.Infoln("shutdown gracefully...")
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		return err
	}
	logrus.Info("server gracefully stopped")
	return nil
}

func initConfig(configDir, configFile, envFile string) (*config.Config, error) {
	err := config.InitENV(envFile)
	if err != nil {
		return nil, err
	}

	cfg, err := config.New(configDir, configFile)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}

func initMinioClient(cfg *config.Config) (*minio.Client, error) {
	minioRootUser := os.Getenv("MINIO_ROOT_USER")
	minioRootPassword := os.Getenv("MINIO_ROOT_PASSWORD")

	minioClient, err := minio.New(cfg.Minio.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(minioRootUser, minioRootPassword, ""),
		Secure: false,
	})
	if err != nil {
		return nil, err
	}

	return minioClient, nil
}

func initStorage(client *minio.Client, cfg *config.Config) (*storage.FileStorage, error) {
	provider := storage.NewProvider(client, cfg.Minio.Bucket, cfg.Minio.Endpoint)
	policy := `{"Version":"2012-10-17","Statement":[{"Effect":"Allow","Principal":{"AWS":["*"]},"Action":["s3:GetObject"],"Resource":["arn:aws:s3:::` + cfg.Minio.Bucket + `/*"]}]}`

	err := client.SetBucketPolicy(context.Background(), cfg.Minio.Bucket, policy)
	if err != nil {
		return nil, err
	}

	return provider, nil
}

func initHandler(cfg *config.Config, provider *storage.FileStorage, db *sql.DB) rest.Handler {
	usersRepo := repository.NewUsers(db)
	tokenRepo := repository.NewTokens(db)
	filesRepo := repository.NewFiles(db)

	hasher := hash.NewSHA1Hasher(cfg.Hash.Salt)
	users := service.NewUsers(hasher, usersRepo, tokenRepo, []byte(cfg.Hash.Secret), cfg.Auth.TokenTTL, cfg.Auth.RefreshTTL)
	files := service.NewFiles(filesRepo, provider)

	return rest.NewHandler(users, files)
}
