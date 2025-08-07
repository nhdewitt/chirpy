package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/nhdewitt/chirpy/internal/auth"
	"github.com/nhdewitt/chirpy/internal/models"
)

func (cfg *APIConfig) deleteChirp(w http.ResponseWriter, r *http.Request, id uuid.UUID) {
	var resp models.ChirpResponse

	chirp, err := cfg.Queries.GetOneChirp(r.Context(), id)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(404)
		log.Printf("Chirp not found: %s (%s)", err, id)
		resp.Error = "Chirp not found"
		json.NewEncoder(w).Encode(resp)
		return
	}

	bearerToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		w.WriteHeader(401)
		log.Printf("Error getting token: %s", err)
		return
	}

	uid, err := auth.ValidateJWT(bearerToken, cfg.Secret)
	if err != nil {
		log.Printf("Error validating JWT: %s", err)
		w.WriteHeader(401)
		return
	}
	if uid != chirp.UserID {
		log.Printf("Invalid user: %s", uid)
		w.WriteHeader(403)
		return
	}

	err = cfg.Queries.DeleteChirp(r.Context(), chirp.ID)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(500)
		log.Printf("Error deleting chirp: %s", err)
		resp.Error = "Internal Server Error"
		json.NewEncoder(w).Encode(resp)
		return
	}

	w.WriteHeader(204)
}
