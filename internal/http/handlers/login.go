package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

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

	params.ExpiresInSeconds = 3600

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

	jwt, err := auth.MakeJWT(uid, cfg.Secret, time.Duration(params.ExpiresInSeconds)*time.Second)
	if err != nil {
		w.WriteHeader(500)
		log.Printf("Error creating JWT: %s", err)
		resp.Error = "Error creating JWT"
		json.NewEncoder(w).Encode(resp)
		return
	}

	refreshToken, err := auth.MakeRefreshToken()
	if err != nil {
		w.WriteHeader(500)
		log.Printf("Error creating refresh token: %s", err)
		resp.Error = "Error creating refresh token"
		json.NewEncoder(w).Encode(resp)
		return
	}

	err = cfg.Queries.CreateRefreshToken(r.Context(), database.CreateRefreshTokenParams{
		Token:     refreshToken,
		UserID:    uid,
		ExpiresAt: time.Now().UTC().AddDate(0, 0, 60),
	})
	if err != nil {
		w.WriteHeader(500)
		log.Printf("DB error: %s", err)
		resp.Error = "Refresh token error"
		json.NewEncoder(w).Encode(resp)
		return
	}

	resp.ID = uid
	resp.CreatedAt = userDetails.CreatedAt
	resp.UpdatedAt = userDetails.UpdatedAt
	resp.Email = params.Email
	resp.Token = jwt
	resp.RefreshToken = refreshToken
	resp.IsChirpyRed = userDetails.IsChirpyRed

	json.NewEncoder(w).Encode(resp)
}
