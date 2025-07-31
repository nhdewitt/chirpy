package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/nhdewitt/chirpy/internal/models"
)

func (cfg *APIConfig) createUser(w http.ResponseWriter, r *http.Request) {
	var resp models.User

	w.Header().Set("Content-Type", "application/json")

	decoder := json.NewDecoder(r.Body)
	var params models.Email
	err := decoder.Decode(&params)
	if err != nil {
		w.WriteHeader(500)
		log.Printf("Error decoding parameters: %s", err)
		resp.Error = "Something went wrong"
		json.NewEncoder(w).Encode(resp)
		return
	}

	user, err := cfg.Queries.CreateUser(r.Context(), params.Email)
	if err != nil {
		w.WriteHeader(500)
		log.Printf("Error creating user: %s", err)
		resp.Error = "Something went wrong"
		json.NewEncoder(w).Encode(resp)
		return
	}

	resp.ID = user.ID
	resp.CreatedAt = user.CreatedAt
	resp.UpdatedAt = user.UpdatedAt
	resp.Email = user.Email

	w.WriteHeader(201)
	json.NewEncoder(w).Encode(resp)
}
