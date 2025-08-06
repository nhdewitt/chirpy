package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/nhdewitt/chirpy/internal/auth"
	"github.com/nhdewitt/chirpy/internal/database"
	"github.com/nhdewitt/chirpy/internal/models"
)

func (cfg *APIConfig) postChirp(w http.ResponseWriter, r *http.Request) {
	var resp models.ChirpResponse
	w.Header().Set("Content-Type", "application/json")

	decoder := json.NewDecoder(r.Body)
	var params models.ChirpPost
	err := decoder.Decode(&params)
	if err != nil {
		w.WriteHeader(500)
		log.Printf("Error processing new chirp: %s", err)
		resp.Error = "Internal server error"
		json.NewEncoder(w).Encode(resp)
		return
	}

	bearerToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		w.WriteHeader(401)
		log.Printf("Unauthorized request: %s", err)
		resp.Error = "Unauthorized"
		json.NewEncoder(w).Encode(resp)
		return
	}

	userID, err := auth.ValidateJWT(bearerToken, cfg.Secret)
	if err != nil {
		w.WriteHeader(401)
		log.Printf("Unauthorized request: %s", err)
		resp.Error = "Unauthorized"
		json.NewEncoder(w).Encode(resp)
		return
	}

	if len(params.Body) > 140 {
		resp.Error = "Chirp is too long"
		w.WriteHeader(400)
		json.NewEncoder(w).Encode(resp)
		return
	}

	chirp, err := cfg.Queries.PostChirp(r.Context(), database.PostChirpParams{
		UserID: userID,
		Body:   params.Body,
	})
	if err != nil {
		w.WriteHeader(500)
		log.Printf("Error posting chirp: %s", err)
		resp.Error = "Error posting chirp"
		json.NewEncoder(w).Encode(resp)
		return
	}

	resp.ID = chirp.ID
	resp.UserID = userID
	resp.CreatedAt = chirp.CreatedAt
	resp.UpdatedAt = chirp.UpdatedAt
	resp.Body = chirp.Body

	w.WriteHeader(201)
	json.NewEncoder(w).Encode(resp)
}
