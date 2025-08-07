package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/nhdewitt/chirpy/internal/auth"
	"github.com/nhdewitt/chirpy/internal/models"
)

func (cfg *APIConfig) upgradeUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	decoder := json.NewDecoder(r.Body)

	apiKey, err := auth.GetAPIKey(r.Header)
	if err != nil {
		w.WriteHeader(401)
		log.Printf("Error getting API key")
		return
	}
	if apiKey != os.Getenv("POLKA_KEY") {
		w.WriteHeader(401)
		log.Printf("Invalid API key")
		return
	}

	var params models.PolkaWebhook
	err = decoder.Decode(&params)
	if err != nil {
		w.WriteHeader(500)
		log.Printf("Error decoding request: %s", err)
		return
	}

	if params.Event != "user.upgraded" {
		w.WriteHeader(204)
		return
	}

	err = cfg.Queries.UpgradeUserToChirpyRed(r.Context(), params.Data.UserID)
	if err != nil {
		w.WriteHeader(404)
		log.Printf("Error upgrading user: %s", err)
		return
	}

	w.WriteHeader(204)
}
