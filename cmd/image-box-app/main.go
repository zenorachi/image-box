package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/zenorachi/image-box/internal/config"
	"github.com/zenorachi/image-box/internal/repository"
	"github.com/zenorachi/image-box/internal/service"
	"github.com/zenorachi/image-box/internal/transport/rest"
	"github.com/zenorachi/image-box/pkg/database/postgres"
	"github.com/zenorachi/image-box/pkg/hash"
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

	hasher := hash.NewSHA1Hasher("testLol")

	usersRepo := repository.NewUsers(db)
	users := service.NewUsers(hasher, usersRepo, []byte("kekSecret"), cfg.Auth.TTL)

	handler := rest.NewHandler(users)

	s := rest.NewServer(handler, cfg.Server.Host, cfg.Server.Port)
	fmt.Println(cfg)

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
