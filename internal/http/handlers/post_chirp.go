package handlers

import (
	"encoding/json"
	"log"
	"net/http"

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
		resp.Error = "Something went wrong"
		json.NewEncoder(w).Encode(resp)
		return
	}

	if len(params.Body) > 140 {
		resp.Error = "Chirp is too long"
		w.WriteHeader(400)
		json.NewEncoder(w).Encode(resp)
		return
	}

	params.Body = cleanBody(params.Body)
	chirp, err := cfg.Queries.PostChirp(r.Context(), database.PostChirpParams{
		UserID: params.UserID,
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
	resp.UserID = chirp.UserID
	resp.CreatedAt = chirp.CreatedAt
	resp.UpdatedAt = chirp.UpdatedAt
	resp.Body = chirp.Body

	w.WriteHeader(201)
	json.NewEncoder(w).Encode(resp)
}
