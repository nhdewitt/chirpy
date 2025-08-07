package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/nhdewitt/chirpy/internal/auth"
	"github.com/nhdewitt/chirpy/internal/database"
	"github.com/nhdewitt/chirpy/internal/models"
)

func (cfg *APIConfig) createUser(w http.ResponseWriter, r *http.Request) {
	var resp models.User

	w.Header().Set("Content-Type", "application/json")

	decoder := json.NewDecoder(r.Body)
	var params models.UserLogin
	err := decoder.Decode(&params)
	if err != nil {
		w.WriteHeader(500)
		log.Printf("Error decoding parameters: %s", err)
		resp.Error = "Error decoding parameters"
		json.NewEncoder(w).Encode(resp)
		return
	}

	hashedPwd, err := auth.HashPassword(params.Password)
	if err != nil {
		w.WriteHeader(500)
		log.Printf("Error hashing password: %s", err)
		resp.Error = "Error hashing password"
		json.NewEncoder(w).Encode(resp)
		return
	}

	user, err := cfg.Queries.CreateUser(r.Context(), database.CreateUserParams{
		Email:          params.Email,
		HashedPassword: hashedPwd,
	})
	if err != nil {
		w.WriteHeader(500)
		log.Printf("Error creating user: %s", err)
		resp.Error = "Error creating user"
		json.NewEncoder(w).Encode(resp)
		return
	}

	resp.ID = user.ID
	resp.CreatedAt = user.CreatedAt
	resp.UpdatedAt = user.UpdatedAt
	resp.Email = user.Email
	resp.IsChirpyRed = user.IsChirpyRed

	w.WriteHeader(201)
	json.NewEncoder(w).Encode(resp)
}
