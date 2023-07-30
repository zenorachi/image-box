package postgres

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
)

type DBConfig struct {
	Host     string
	Port     int
	Username string
	Name     string
	SSLMode  string
	Password string
}

func NewDB(cfg DBConfig) (*sql.DB, error) {
	source := fmt.Sprintf("host=%s port=%d user=%s dbname=%s sslmode=%s password=%s", cfg.Host, cfg.Port, cfg.Username, cfg.Name, cfg.SSLMode, cfg.Password)
	db, err := sql.Open("postgres", source)

	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}
