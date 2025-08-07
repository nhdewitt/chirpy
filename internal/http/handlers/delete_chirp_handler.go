package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/nhdewitt/chirpy/internal/models"
)

func (cfg *APIConfig) deleteChirpHandler(w http.ResponseWriter, r *http.Request) {
	chirpID, err := uuid.Parse(r.PathValue("chirpID"))
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		log.Printf("Error parsing chirpID: %s", err)
		resp := models.User{
			Error: "error parsing chirpID",
		}
		json.NewEncoder(w).Encode(resp)
		return
	}

	cfg.deleteChirp(w, r, chirpID)
}
