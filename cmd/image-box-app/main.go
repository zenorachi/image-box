package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/zenorachi/image-box/internal/config"
	"github.com/zenorachi/image-box/internal/repository"
	"github.com/zenorachi/image-box/internal/service"
	"github.com/zenorachi/image-box/internal/transport/rest"
	"github.com/zenorachi/image-box/pkg/database/postgres"
	"github.com/zenorachi/image-box/pkg/hash"
	"log"
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

	fmt.Println(cfg)

	db, err := postgres.NewDB(cfg.DB)
	if err != nil {
		log.Fatalln(err)
	}
	defer db.Close()

	hasher := hash.NewSHA1Hasher("testLol")

	usersRepo := repository.NewUsers(db)
	users := service.NewUsers(hasher, usersRepo, []byte("kek"))

	handler := rest.NewHandler(users)

	s := rest.NewServer(handler)
	if err := s.Run(); err != nil {
		log.Fatalln(err)
	}
}
