package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"game/internal/config"
	"game/internal/service"
	"net/http"
	"os"
	"strconv"

	"go.uber.org/zap"
)

type Decorator func(http.Handler) http.Handler

type LifeStates struct {
	service.LifeService
}

type SetRequest struct {
	Fill int `json:"fill"`
}

func New(ctx context.Context, logger *zap.Logger, cfg *config.Config, lifeService *service.LifeService) http.Handler {

	mux := http.NewServeMux()

	lifeStates := LifeStates{
		LifeService: *lifeService,
	}

	mux.HandleFunc("/nextstate", NextState(logger, cfg, &lifeStates))
	mux.HandleFunc("/setstate", SetState(logger, cfg, &lifeStates))
	return mux
}

func Decorate(next http.Handler, ds ...Decorator) http.Handler {
	decorated := next
	for d := len(ds) - 1; d >= 0; d-- {
		decorated = ds[d](decorated)
	}

	return decorated
}

func NextState(logger *zap.Logger, cfg *config.Config, ls *LifeStates) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		worldState := ls.NextState()

		str := worldState.String()
		if str == "" {
			http.Error(w, "life is empty", http.StatusInternalServerError)
			return
		}
		fmt.Fprintf(w, "Fill: %d\n", worldState.Fill)
		fmt.Fprint(w, str)

		logger.Info("NextState request processed succesfully")

		state_file, err := os.OpenFile(cfg.StatesPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
		if err != nil {
			logger.Sugar().Errorf("Cannot open state file: %v\n", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		defer state_file.Close()
		state_file.WriteString(strconv.Itoa(worldState.Fill) + "\n")
		state_file.WriteString(str)
	}
}

func SetState(logger *zap.Logger, cfg *config.Config, ls *LifeStates) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req SetRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			logger.Sugar().Errorf("Cannot decode request body: %v\n", err)
			http.Error(w, fmt.Sprintf("Error with decode request body: %v", err), http.StatusInternalServerError)
			return
		}

		world := ls.SetState(req.Fill)
		str := world.String()
		if str == "" {
			http.Error(w, "life is empty", http.StatusInternalServerError)
			return
		}
		fmt.Fprintf(w, "Fill: %d\n", req.Fill)
		fmt.Fprint(w, str)

		state_file, err := os.OpenFile(cfg.StatesPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
		if err != nil {
			logger.Sugar().Errorf("Cannot open state file: %v\n", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		defer state_file.Close()
		state_file.WriteString(strconv.Itoa(req.Fill) + "\n")
		state_file.WriteString(str)
	}
}
