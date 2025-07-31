package handlers

import "net/http"

func (cfg *APIConfig) resetMetrics(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	if cfg.Platform != "dev" {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	err := cfg.Queries.DeleteAllUsers(r.Context())
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
