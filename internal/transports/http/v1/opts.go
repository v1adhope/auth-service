package httpv1

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type Option func(*Config)

type Config struct {
	Cors cors.Config
	Mode string
}

func WithAllowOrigins(ao []string) Option {
	return func(cfg *Config) {
		cfg.Cors.AllowOrigins = ao
	}
}

func WithAllowMethods(am []string) Option {
	return func(cfg *Config) {
		cfg.Cors.AllowMethods = am
	}
}

func WithAllowHeaders(ah []string) Option {
	return func(cfg *Config) {
		cfg.Cors.AllowHeaders = ah
	}
}

func WithMode(m string) Option {
	return func(cfg *Config) {
		cfg.Mode = m
	}
}

func config(opts ...Option) Config {
	cfg := Config{
		Cors: cors.Config{
			AllowOrigins: []string{"*"},
			AllowMethods: []string{"POST", "HEAD", "OPTIONS"},
			AllowHeaders: []string{"Origin"},
		},
		Mode: gin.DebugMode,
	}

	for _, opt := range opts {
		opt(&cfg)
	}

	return cfg
}
