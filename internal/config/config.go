package config

import (
	"fmt"
	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
	"github.com/spf13/viper"
	"github.com/zenorachi/image-box/pkg/database/postgres"
	"time"
)

type Config struct {
	Server struct {
		Port int
	}

	Auth struct {
		TTL time.Duration
	}

	DB postgres.DBConfig
}

func InitENV(filename string) error {
	return godotenv.Load(filename)
}

func New(directory, filename string) (*Config, error) {
	cfg := new(Config)

	viper.AddConfigPath(directory)
	viper.SetConfigName(filename)

	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	if err := viper.Unmarshal(cfg); err != nil {
		return nil, err
	}
	fmt.Println(cfg.Auth)

	if err := envconfig.Process("db", &cfg.DB); err != nil {
		return nil, err
	}

	return cfg, nil
}
