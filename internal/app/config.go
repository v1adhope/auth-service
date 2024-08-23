package app

import (
	"fmt"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
)

type (
	Config struct {
		Tokens   Tokens
		Postgres Postgres
		Logger   Logger
		Server   Server
	}

	Tokens struct {
		AccessKey  string        `env-required:"true" env:"APP_TOKENS_ACCESS_KEY"`
		AceessTtl  time.Duration `env-required:"true" env:"APP_TOKENS_ACCESS_TTL"`
		RefreshKey string        `env-required:"true" env:"APP_TOKENS_REFRESH_KEY"`
		Issuer     string        `env-required:"true" env:"APP_TOKENS_ISSUER"`
	}

	Postgres struct {
		ConnStr string `env-required:"true" env:"APP_POSTGRES_CONN_STR"`
	}

	Logger struct {
		Level string `env-required:"true" env:"APP_LOGGER_LEVEL"`
	}

	Server struct {
		AllowOrigins    []string      `env-required:"true" env-separator:":" env:"APP_SERVER_ALLOW_ORIGINS"`
		AllowMethods    []string      `env-required:"true" env-separator:":" env:"APP_SERVER_ALLOW_METHODS"`
		AllowHeaders    []string      `env-required:"true" env-separator:":" env:"APP_SERVER_ALLOW_HEADERS"`
		Mode            string        `env-required:"true" env:"APP_SERVER_MODE"`
		Socket          string        `env-required:"true" env:"APP_SERVER_SOCKET"`
		ShutdownTimeout time.Duration `env-required:"true" env:"APP_SERVER_SHUTDOWN_TIMEOUT"`
		WriteTimeout    time.Duration `env-required:"true" env:"APP_SERVER_WRITE_TIMEOUT"`
		ReadTimeout     time.Duration `env-required:"true" env:"APP_SERVER_READ_TIMEOUT"`
	}
)

func MustConfig() Config {
	if err := godotenv.Load(); err != nil {
		panic(fmt.Errorf("config: can't load envs from .env: %v", err))
	}

	cfg := Config{}

	if err := cleanenv.ReadEnv(&cfg); err != nil {
		panic(fmt.Sprintf("config: can't read envs: %v", err))
	}

	return cfg
}
