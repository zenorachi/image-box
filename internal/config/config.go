package config

import (
	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
	"github.com/spf13/viper"
	"github.com/zenorachi/image-box/pkg/database/postgres"
)

type Config struct {
	Server struct {
		Port int
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

	if err := envconfig.Process("db", &cfg.DB); err != nil {
		return nil, err
	}

	return cfg, nil
}
