package server

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"game/http/server/handler"
	"game/internal/config"
	"game/internal/service"
	"io"
	"net/http"
	"strconv"
	"time"

	"go.uber.org/zap"
)

func new(ctx context.Context, logger *zap.Logger, lifeService *service.LifeService, cfg *config.Config) (http.Handler, error) {
	muxHanler := handler.New(ctx, logger, cfg, lifeService)

	muxHanler = handler.Decorate(muxHanler, contextMiddleware(logger), loggingMiddleware(logger))

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

func contextMiddleware(logger *zap.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fill := r.URL.Query().Get("fill")
			if fill == "" {
				next.ServeHTTP(w, r)
				return
			}

			fillInt, err := strconv.Atoi(fill)
			if err != nil {
				logger.Sugar().Errorf("Error with convert fill to int: %v\n", err)
				http.Error(w, fmt.Sprintf("Error with convert fill to int: %v\n", err), http.StatusBadRequest)
			}

			request := handler.SetRequest{Fill: fillInt}
			requestBytes, err := json.Marshal(request)
			if err != nil {
				http.Error(w, fmt.Sprintf("Error with encode request body: %v\n", err), http.StatusInternalServerError)
				logger.Sugar().Errorf("Error with encode request body: %v\n", err)
				return
			}

			r.Body = io.NopCloser(bytes.NewReader(requestBytes))
			r.ContentLength = int64(len(requestBytes))
			next.ServeHTTP(w, r)
		})
	}
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
