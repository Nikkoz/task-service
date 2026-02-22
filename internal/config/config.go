package config

import (
	"fmt"
	"time"

	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
)

type (
	Config struct {
		App  App  `envPrefix:"APP_"`
		Http Http `envPrefix:"HTTP_"`
		Db   Db   `envPrefix:"DB_"`
		Log  Log  `envPrefix:"LOG_"`
	}

	App struct {
		Name        string      `env:"NAME,required"`
		Version     string      `env:"VERSION,required"`
		Environment Environment `env:"ENV" envDefault:"local"`
	}

	Http struct {
		Host string `env:"HOST" envDefault:"localhost"`
		Port uint16 `env:"PORT" envDefault:"8080"`
	}

	Db struct {
		Host     string `env:"HOST" envDefault:"localhost"`
		Port     uint16 `env:"PORT" envDefault:"5432"`
		Name     string `env:"NAME,required"`
		User     string `env:"USER,required"`
		Password string `env:"PASSWORD,required"`
		SslMode  bool   `env:"USE_SSL" envDefault:"false"`

		MaxConns        int32         `env:"DB_MAX_CONNS" envDefault:"10"`
		MinConns        int32         `env:"DB_MIN_CONNS" envDefault:"0"`
		MaxConnLifetime time.Duration `env:"DB_MAX_CONN_LIFETIME" envDefault:"1h"`
		MaxConnIdleTime time.Duration `env:"DB_MAX_CONN_IDLE_TIME" envDefault:"30m"`
	}

	Log struct {
		Level LogLevel `env:"LEVEL" envDefault:"debug"`
	}
)

func Load() (Config, error) {
	if err := envInit(); err != nil {
		return Config{}, err
	}

	var cfg Config
	if err := env.Parse(&cfg); err != nil {
		return Config{}, err
	}

	return cfg, nil
}

func envInit() error {
	if err := godotenv.Load(); err != nil {
		return fmt.Errorf("warning: .env not loaded: %v\n", err)
	}

	return nil
}
