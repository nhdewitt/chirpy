package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/nhdewitt/chirpy/internal/auth"
	"github.com/nhdewitt/chirpy/internal/database"
	"github.com/nhdewitt/chirpy/internal/models"
)

func (cfg *APIConfig) login(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var resp models.User

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

	var userDetails database.QueryUserRow

	userDetails, err = cfg.Queries.QueryUser(r.Context(), params.Email)
	if err != nil {
		w.WriteHeader(500)
		log.Printf("Error retrieving User ID: %s", err)
		resp.Error = "Error retrieving User ID"
		json.NewEncoder(w).Encode(resp)
		return
	}

	hashedPwd := userDetails.HashedPassword
	uid := userDetails.ID

	err = auth.CheckPasswordHash(params.Password, hashedPwd)
	if err != nil {
		w.WriteHeader(401)
		log.Printf("Incorrect email or password: %s", err)
		resp.Error = "Incorrect email or password"
		json.NewEncoder(w).Encode(resp)
		return
	}

	resp.ID = uid
	resp.CreatedAt = userDetails.CreatedAt
	resp.UpdatedAt = userDetails.UpdatedAt
	resp.Email = params.Email

	json.NewEncoder(w).Encode(resp)
}
