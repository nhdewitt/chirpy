package handlers

import "net/http"

func (cfg *APIConfig) handleChirps(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	switch r.Method {
	case http.MethodGet:
		cfg.getAllChirps(w, r)
	case http.MethodPost:
		cfg.postChirp(w, r)
	default:
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
	}
}
