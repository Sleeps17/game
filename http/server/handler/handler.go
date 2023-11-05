package handler

import (
	"bufio"
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

	mux.HandleFunc("/nextstate", NextState(logger, &lifeStates))
	mux.HandleFunc("/setstate", SetState(logger, cfg, &lifeStates))
	mux.HandleFunc("/reset", Reset(logger, cfg, &lifeStates))
	return mux
}

func Decorate(next http.Handler, ds ...Decorator) http.Handler {
	decorated := next
	for d := len(ds) - 1; d >= 0; d-- {
		decorated = ds[d](decorated)
	}

	return decorated
}

func NextState(logger *zap.Logger, ls *LifeStates) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		worldState := ls.LifeService.NextState()

		str := worldState.String()
		if str == "" {
			http.Error(w, "life is empty", http.StatusInternalServerError)
			return
		}
		logger.Info("NextState request processed succesfully")
		fmt.Fprint(w, str)
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

		ls.CurrentWorld.Seed(req.Fill)
		state_file, err := os.OpenFile(cfg.StatesPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
		if err != nil {
			logger.Sugar().Errorf("Cannot open state file: %v\n", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		defer state_file.Close()
		state_file.WriteString(strconv.Itoa(req.Fill) + "\n")
		fmt.Fprintf(w, "Fill: %d\n", req.Fill)
		fmt.Fprint(w, ls.CurrentWorld.String())
	}
}

func Reset(logger *zap.Logger, cfg *config.Config, ls *LifeStates) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		states, err := os.OpenFile(cfg.StatesPath, os.O_RDONLY, 0644)
		if err != nil {
			logger.Error("Cannot open states_file: ", zap.Error(err))
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer states.Close()
		sc := bufio.NewScanner(states)
		var fill string
		for sc.Scan() {
			if sc.Text() != "" {
				fill = sc.Text()
			}
		}
		if fill == "" {
			logger.Error("States file is empty")
			http.Error(w, "States file is empty", http.StatusBadRequest)
			return
		}
		intFill, err := strconv.Atoi(fill)
		if err != nil {
			logger.Error("Cannot convert fill to int: ", zap.Error(err))
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		ls.CurrentWorld.Seed(intFill)

		if err := json.NewEncoder(w).Encode("New Fill: " + fill); err != nil {
			logger.Sugar().Errorf("Cannot encode response: %v\n", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		logger.Sugar().Infoln("Reset request succesfully processed")
		fmt.Fprint(w, ls.CurrentWorld.String())
	}
}
