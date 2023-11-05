package application

import (
	"context"
	"game/http/server"
	"game/internal/config"
	"log"
	"os"
	"os/signal"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Application struct {
	Cfg config.Config
}

func New(cfg *config.Config) *Application {
	return &Application{
		Cfg: *cfg,
	}
}

func (a *Application) Run(ctx context.Context) int {
	logger := SetupLogger(a.Cfg.LogsPath)

	shutDownFunc, err := server.Run(ctx, logger, &a.Cfg)
	if err != nil {
		logger.Error(err.Error())
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	ctx, cancel := context.WithCancel(ctx)

	ctx, stop := signal.NotifyContext(ctx, os.Interrupt)
	defer stop()

	<-c
	cancel()

	shutDownFunc(ctx)

	return 0
}

func SetupLogger(logs_path string) *zap.Logger {
	config := zap.NewProductionConfig()

	config.Level = zap.NewAtomicLevelAt(zapcore.InfoLevel)
	config.OutputPaths = []string{logs_path}

	logger, err := config.Build()
	if err != nil {
		log.Fatalf("Error with config logger: %v\n", err)
	}
	return logger
}
