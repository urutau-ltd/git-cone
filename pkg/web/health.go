package web

import (
	"context"
	"encoding/json"
	"net/http"

	"charm.land/log/v2"
	"github.com/urutau-ltd/git-cone/pkg/db"
	"github.com/urutau-ltd/git-cone/pkg/version"
	"github.com/gorilla/mux"
)

// HealthController registers the health check routes for the web server.
// These routes are registered before any authentication middleware.
func HealthController(_ context.Context, r *mux.Router) {
	r.HandleFunc("/health", getHealth)
	r.HandleFunc("/livez", getLiveness)
	r.HandleFunc("/readyz", getReadiness)
}

type healthResponse struct {
	Version string `json:"version"`
	Status  string `json:"status"`
}

func getHealth(w http.ResponseWriter, _ *http.Request) {
	ver := version.Version
	if ver == "" {
		ver = "unknown"
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(healthResponse{Version: ver, Status: "ok"})
}

func getLiveness(w http.ResponseWriter, _ *http.Request) {
	renderStatus(http.StatusOK)(w, nil)
}

func getReadiness(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := log.FromContext(ctx)
	db := db.FromContext(ctx)

	if err := db.PingContext(ctx); err != nil {
		logger.Error("error getting db readiness", "err", err)
		renderStatus(http.StatusServiceUnavailable)(w, nil)
		return
	}

	renderStatus(http.StatusOK)(w, nil)
}
