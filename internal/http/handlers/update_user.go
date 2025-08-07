package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/nhdewitt/chirpy/internal/auth"
	"github.com/nhdewitt/chirpy/internal/database"
	"github.com/nhdewitt/chirpy/internal/models"
)

func (cfg *APIConfig) updateUser(w http.ResponseWriter, r *http.Request) {
	var resp models.User

	w.Header().Set("Content-Type", "application/json")

	decoder := json.NewDecoder(r.Body)
	var params models.UserLogin
	err := decoder.Decode(&params)
	if err != nil {
		w.WriteHeader(500)
		log.Printf("Error decoding parameters: %s", err)
		resp.Error = "Internal Server Error"
		json.NewEncoder(w).Encode(resp)
		return
	}

	accessToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		w.WriteHeader(401)
		log.Printf("Error getting token: %s", err)
		resp.Error = "Unauthorized"
		json.NewEncoder(w).Encode(resp)
		return
	}

	uid, err := auth.ValidateJWT(accessToken, cfg.Secret)
	if err != nil {
		w.WriteHeader(401)
		log.Printf("Unable to validate JWT: %s", err)
		resp.Error = "Unauthorized"
		json.NewEncoder(w).Encode(resp)
		return
	}

	hashedPwd, err := auth.HashPassword(params.Password)
	if err != nil {
		w.WriteHeader(500)
		log.Printf("Error hashing password: %s", err)
		resp.Error = "Internal Server Error"
		json.NewEncoder(w).Encode(resp)
		return
	}

	updatePasswordRow, err := cfg.Queries.UpdatePassword(r.Context(), database.UpdatePasswordParams{
		Email:          params.Email,
		HashedPassword: hashedPwd,
		ID:             uid,
	})
	if err != nil {
		w.WriteHeader(500)
		log.Printf("Error updating password: %s", err)
		resp.Error = "Internal Server Error"
		json.NewEncoder(w).Encode(resp)
		return
	}

	resp.ID = updatePasswordRow.ID
	resp.CreatedAt = updatePasswordRow.CreatedAt
	resp.UpdatedAt = updatePasswordRow.UpdatedAt
	resp.Email = updatePasswordRow.Email

	w.WriteHeader(200)
	json.NewEncoder(w).Encode(resp)
}
