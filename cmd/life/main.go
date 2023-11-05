package main

import (
	"context"
	"game/internal/application"
	"game/internal/config"
	"os"
)

func main() {
	ctx := context.Background()
	os.Exit(mainWithExitCode(ctx))
}

func mainWithExitCode(ctx context.Context) int {
	cfg := config.MustLoad()
	app := application.New(cfg)

	return app.Run(ctx)
}
