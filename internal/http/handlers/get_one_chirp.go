package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/nhdewitt/chirpy/internal/models"
)

func (cfg *APIConfig) getOneChirp(w http.ResponseWriter, r *http.Request, id uuid.UUID) {
	w.Header().Set("Content-Type", "application/json")

	chirp, err := cfg.Queries.GetOneChirp(r.Context(), id)
	if err != nil {
		var resp models.ChirpResponse
		w.WriteHeader(404)
		log.Printf("Chirp not found: %s", err)
		resp.Error = "Chirp not found"
		json.NewEncoder(w).Encode(resp)
		return
	}

	resp := models.ChirpResponse{
		ID:        chirp.ID,
		CreatedAt: chirp.CreatedAt,
		UpdatedAt: chirp.UpdatedAt,
		Body:      chirp.Body,
		UserID:    chirp.UserID,
	}

	w.WriteHeader(200)
	json.NewEncoder(w).Encode(resp)
}
