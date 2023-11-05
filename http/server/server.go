package server

import (
	"context"
	"game/http/server/handler"
	"game/internal/config"
	"game/internal/service"
	"net/http"
	"time"

	"go.uber.org/zap"
)

func new(ctx context.Context, logger *zap.Logger, lifeService *service.LifeService, cfg *config.Config) (http.Handler, error) {
	muxHanler := handler.New(ctx, logger, cfg, lifeService)

	muxHanler = handler.Decorate(muxHanler, loggingMiddleware(logger))

	return muxHanler, nil
}

func Run(ctx context.Context, logger *zap.Logger, cfg *config.Config) (func(context.Context) error, error) {
	lifeService, err := service.New(cfg.Height, cfg.Width, cfg.Fill)
	if err != nil {
		return nil, err
	}

	muxHandler, err := new(ctx, logger, lifeService, cfg)
	if err != nil {
		return nil, err
	}

	srv := &http.Server{
		Addr:         cfg.Addres,
		Handler:      muxHandler,
		ReadTimeout:  cfg.Timeout,
		WriteTimeout: cfg.Timeout,
		IdleTimeout:  cfg.IdleTimeout,
	}

	go func() {
		// Запускаем сервер
		if err := srv.ListenAndServe(); err != nil {
			logger.Error("ListenAndServe",
				zap.String("err", err.Error()))
		}
	}()
	return srv.Shutdown, nil
}

func loggingMiddleware(logger *zap.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			next.ServeHTTP(w, r)

			// Завершение логирования после выполнения запроса
			duration := time.Since(start)
			logger.Info("HTTP request",
				zap.String("method", r.Method),
				zap.String("path", r.URL.Path),
				zap.Duration("duration", duration),
			)
		})
	}
}
