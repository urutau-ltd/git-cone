package web

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/urutau-ltd/git-cone/pkg/version"
)

// HealthController registers the health check route for the web server.
func HealthController(r *mux.Router) {
	r.HandleFunc("/health", getHealth)
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
