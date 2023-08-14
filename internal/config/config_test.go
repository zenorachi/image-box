package config

import (
	"reflect"
	"testing"
	"time"

	"github.com/joho/godotenv"
	"github.com/zenorachi/image-box/pkg/database/postgres"
)

const (
	envFile   = "./fixtures/.env"
	directory = "./fixtures"
	ymlFile   = "main"
)

func TestConfig(t *testing.T) {
	type args struct {
		directoryPath string
		ymlFile       string
	}

	initEnv := func() {
		_ = godotenv.Load(envFile)
	}

	tests := []struct {
		name    string
		args    args
		want    *Config
		wantErr bool
	}{
		{
			name: "test config",
			args: args{
				directoryPath: directory,
				ymlFile:       ymlFile,
			},
			want: &Config{
				Server: struct {
					Host string
					Port int
				}{Host: "0.0.0.0", Port: 8080},

				Auth: struct {
					TokenTTL   time.Duration
					RefreshTTL time.Duration
				}{TokenTTL: time.Minute * 15, RefreshTTL: time.Hour * 744},

				Minio: struct {
					Endpoint string
					Bucket   string
				}{Endpoint: "minio:9000", Bucket: "image-box"},

				Hash: struct {
					Salt   string
					Secret string
				}{Salt: "salt", Secret: "secret"},

				DB: postgres.DBConfig{
					Host:     "db",
					Port:     5432,
					Username: "postgres",
					Name:     "imagebox",
					SSLMode:  "disable",
					Password: "12345",
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			initEnv()

			got, err := New(tt.args.directoryPath, tt.args.ymlFile)
			if (err != nil) != tt.wantErr {
				t.Errorf("Init() error = %v, wantErr %v", err, tt.wantErr)

				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Init() got = %v, want %v", got, tt.want)
			}
		})
	}
}
