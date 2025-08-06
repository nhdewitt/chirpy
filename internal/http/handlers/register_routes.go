package handlers

import (
	"database/sql"
	"net/http"
	"os"
	"sync/atomic"

	"github.com/nhdewitt/chirpy/internal/database"
)

type APIConfig struct {
	fileserverHits atomic.Int32
	Queries        *database.Queries
	Platform       string
	Secret         string
}

func NewAPIConfig(db *sql.DB) *APIConfig {
	return &APIConfig{
		fileserverHits: atomic.Int32{},
		Queries:        database.New(db),
		Platform:       os.Getenv("PLATFORM"),
		Secret:         os.Getenv("TOKEN_STRING"),
	}
}

func (cfg *APIConfig) RegisterRoutes(mux *http.ServeMux, filepathRoot string) {
	fileServer := http.StripPrefix("/app", http.FileServer(http.Dir(filepathRoot)))
	mux.Handle("/app/", cfg.middlewareMetricsInc(fileServer))

	mux.HandleFunc("GET /api/healthz", handler)

	mux.HandleFunc("GET /admin/metrics", cfg.metricsHandler)
	mux.HandleFunc("POST /admin/reset", cfg.resetMetrics)

	mux.HandleFunc("/api/chirps", cfg.handleChirps)
	mux.HandleFunc("GET /api/chirps/{chirpID}", cfg.getChirpHandler)
	mux.HandleFunc("POST /api/users", cfg.createUser)
	mux.HandleFunc("POST /api/login", cfg.login)
	mux.HandleFunc("POST /api/refresh", cfg.refresh)
	mux.HandleFunc("POST /api/revoke", cfg.revoke)
}
