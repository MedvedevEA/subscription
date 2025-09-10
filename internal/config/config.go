package config

import (
	"log"
	"time"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	Log Log
	Srv Srv
	Db  Db
}
type Log struct {
	Level       string   `envconfig:"LOG_LEVEL" default:"info"`
	OutPutPaths []string `envconfig:"LOG_OUTPUTPATHS" default:"stdout"`
}

type Srv struct {
	Addr         string        `envconfig:"SRV_ADDR" required:"true"`
	WriteTimeout time.Duration `envconfig:"SRV_WRITE_TIMEOUT" required:"true"`
	AppName      string        `envconfig:"SRV_APPNAME" required:"true"`
}

type Db struct {
	Host          string `envconfig:"DB_HOST" required:"true"`
	Port          int    `envconfig:"DB_PORT" required:"true"`
	Name          string `envconfig:"DB_DBNAME" required:"true"`
	User          string `envconfig:"DB_USER" required:"true"`
	Password      string `envconfig:"DB_PASSWORD" required:"true"`
	SSLMode       string `envconfig:"DB_SSL_MODE" default:"disable"`
	MigrationPath string `envconfig:"DB_MIGRATIONS_PATH" required:"true"`
}

func MustNew() *Config {
	if err := godotenv.Load("./configs/config.env"); err != nil {
		log.Printf("failed to load configuration file: %v\n", err)
	}

	cfg := new(Config)
	if err := envconfig.Process("", cfg); err != nil {
		log.Fatalf("failed to load configuration: %v\n", err)
	}
	return cfg
}
