package config

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
)

type (
	Config struct {
		App  App  `envPrefix:"APP_"`
		Http Http `envPrefix:"HTTP_"`
		Db   Db   `envPrefix:"DB_"`
		Auth Auth `envPrefix:"AUTH_"`
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

	Auth struct {
		Token     string        `env:"TOKEN,required"`
		Cost      int           `env:"BCRYPT_COST" envDefault:"10"`
		JwtSecret string        `env:"JWT_SECRET,required"`
		JwtTtl    time.Duration `env:"JWT_TTL" envDefault:"15m"`
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
	envFile := ".env"
	envFile += os.Getenv("ENV_FILE")

	path, err := resolveEnvFile(envFile)
	if err != nil {
		return fmt.Errorf("warning: %s not loaded: %w", envFile, err)
	}

	if err := godotenv.Load(path); err != nil {
		return fmt.Errorf("warning: %s not loaded: %v\n", envFile, err)
	}

	return nil
}

// resolveEnvFile resolves env file path in a way that works under `go test`
// where working dir is the package directory. :contentReference[oaicite:1]{index=1}
func resolveEnvFile(envFile string) (string, error) {
	// absolute path -> use as is
	if filepath.IsAbs(envFile) {
		if _, err := os.Stat(envFile); err != nil {
			return "", err
		}
		return envFile, nil
	}

	// try relative to current working dir first
	if _, err := os.Stat(envFile); err == nil {
		return envFile, nil
	}

	// walk up directories and look for the file
	wd, err := os.Getwd()
	if err != nil {
		return "", err
	}

	dir := wd
	for {
		try := filepath.Join(dir, envFile)
		if _, err := os.Stat(try); err == nil {
			return try, nil
		}

		// optional: stop when we reach module root (go.mod)
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			break
		}

		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}
		dir = parent
	}

	return "", errors.New("file not found (searched upward from working directory)")
}
