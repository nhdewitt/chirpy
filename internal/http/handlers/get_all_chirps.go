package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"slices"

	"github.com/google/uuid"
	"github.com/nhdewitt/chirpy/internal/database"
	"github.com/nhdewitt/chirpy/internal/models"
)

func (cfg *APIConfig) getAllChirps(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var responses []models.ChirpResponse

	author := r.URL.Query().Get("author_id")
	sort := r.URL.Query().Get("sort")

	var chirps []database.Chirp
	var err error
	if author == "" {
		chirps, err = cfg.Queries.GetAllChirps(r.Context())
		if err != nil {
			log.Printf("Error getting all chirps: %s", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
	} else {
		uid, err := uuid.Parse(author)
		if err != nil {
			log.Printf("Unable to parse user ID: %s", err)
			http.Error(w, "Bad Request", http.StatusBadRequest)
			return
		}
		chirps, err = cfg.Queries.GetChirpsFromUser(r.Context(), uid)
		if err != nil {
			log.Printf("Error getting user chirps: %s", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
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

	if sort == "desc" {
		slices.Reverse(responses)
	}

	json.NewEncoder(w).Encode(responses)
}
