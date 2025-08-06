package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/nhdewitt/chirpy/internal/auth"
	"github.com/nhdewitt/chirpy/internal/models"
)

func (cfg *APIConfig) revoke(w http.ResponseWriter, r *http.Request) {
	var httpError models.HTTPError

	refreshToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		w.WriteHeader(500)
		log.Printf("Unable to get refresh token: %s", err)
		httpError.Error = "Internal Server Error"
		json.NewEncoder(w).Encode(httpError)
		return
	}

	err = cfg.Queries.RevokeRefreshToken(r.Context(), refreshToken)
	if err != nil {
		w.WriteHeader(500)
		log.Printf("Unable to revoke refresh token: %s", err)
		httpError.Error = "Internal Server Error"
		json.NewEncoder(w).Encode(httpError)
		return
	}

	w.WriteHeader(204)
}
