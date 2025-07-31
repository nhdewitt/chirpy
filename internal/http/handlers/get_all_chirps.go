package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/nhdewitt/chirpy/internal/models"
)

func (cfg *APIConfig) getAllChirps(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var responses []models.ChirpResponse

	chirps, err := cfg.Queries.GetAllChirps(r.Context())
	if err != nil {
		var resp models.ChirpResponse
		w.WriteHeader(500)
		log.Printf("Error getting all chirps: %s", err)
		resp.Error = "Error getting all chirps"
		json.NewEncoder(w).Encode(resp)
		return
	}

	for _, chirp := range chirps {
		resp := models.ChirpResponse{
			ID:        chirp.ID,
			CreatedAt: chirp.CreatedAt,
			UpdatedAt: chirp.UpdatedAt,
			Body:      chirp.Body,
			UserID:    chirp.UserID,
		}
		responses = append(responses, resp)
	}

	json.NewEncoder(w).Encode(responses)
}
