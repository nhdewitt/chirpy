package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/nhdewitt/chirpy/internal/auth"
	"github.com/nhdewitt/chirpy/internal/database"
	"github.com/nhdewitt/chirpy/internal/models"
)

func (cfg *APIConfig) refresh(w http.ResponseWriter, r *http.Request) {
	var resp models.HTTPError

	w.Header().Set("Content-Type", "application/json")

	bearerToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		w.WriteHeader(500)
		log.Printf("Error getting token: %s", err)
		resp.Error = "Error getting token"
		json.NewEncoder(w).Encode(resp)
		return
	}

	refreshTokenRow, err := cfg.Queries.GetRefreshToken(r.Context(), bearerToken)
	if err != nil || time.Now().UTC().After(refreshTokenRow.ExpiresAt) || refreshTokenRow.RevokedAt.Valid {
		w.WriteHeader(401)
		log.Printf("Token expired or doesn't exist: %s", err)
		resp.Error = "Token expired or doesn't exist"
		json.NewEncoder(w).Encode(resp)
		return
	}

	newToken, err := auth.MakeRefreshToken()
	if err != nil {
		w.WriteHeader(500)
		log.Printf("Error creating refresh token: %s", err)
		resp.Error = "Internal server error"
		json.NewEncoder(w).Encode(resp)
		return
	}

	err = cfg.Queries.GetUserFromRefreshToken(r.Context(), database.GetUserFromRefreshTokenParams{
		Token:     newToken,
		ExpiresAt: time.Now().UTC().AddDate(0, 0, 60),
		Token_2:   refreshTokenRow.Token,
	})
	if err != nil {
		w.WriteHeader(500)
		log.Printf("Error setting new refresh token: %s", err)
		resp.Error = "Internal server error"
		json.NewEncoder(w).Encode(resp)
		return
	}

	var token models.Token

	token.Token, err = auth.MakeJWT(refreshTokenRow.UserID, cfg.Secret, time.Duration(3600*time.Second))
	if err != nil {
		w.WriteHeader(500)
		log.Printf("Error creating JWT: %s", err)
		resp.Error = "Internal server error"
		json.NewEncoder(w).Encode(resp)
		return
	}

	w.WriteHeader(200)
	json.NewEncoder(w).Encode(token)
}
