package httpv1

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func WithAllowOrigins(ao []string) func(*cors.Config) {
	return func(cfg *cors.Config) {
		cfg.AllowOrigins = ao
	}
}

func WithAllowMethods(am []string) func(*cors.Config) {
	return func(cfg *cors.Config) {
		cfg.AllowMethods = am
	}
}

func WithAllowHeaders(ah []string) func(*cors.Config) {
	return func(cfg *cors.Config) {
		cfg.AllowHeaders = ah
	}
}

func corsHandler(opts ...func(*cors.Config)) gin.HandlerFunc {
	cfg := cors.Config{
		AllowOrigins: []string{"localhost"},
		AllowMethods: []string{"POST", "HEAD", "OPTIONS"},
		AllowHeaders: []string{"Origin"},
	}

	for _, opt := range opts {
		opt(&cfg)
	}

	return cors.New(cfg)
}
