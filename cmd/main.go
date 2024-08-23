package main

import (
	"context"
	"log"

	"github.com/v1adhope/auth-service/internal/app"
)

func main() {
	ctx := context.Background()

	cfg := app.MustConfig()

	if err := app.Run(ctx, cfg); err != nil {
		log.Fatal(err)
	}
}
