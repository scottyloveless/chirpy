package main

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"github.com/scottyloveless/chirpy/internal/auth"
)

func (cfg *apiConfig) handlerChirpyRedUpgrade(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Event string `json:"event"`
		Data  struct {
			UserID uuid.UUID `json:"user_id"`
		} `json:"data"`
	}

	apiKey, err := auth.GetAPIKey(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "invalid api key", err)
		return
	}

	if apiKey != cfg.polka_key {
		respondWithError(w, http.StatusUnauthorized, "invalid api key", err)
		return
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err = decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "error decoding json", err)
		return
	}

	if params.Event != "user.upgraded" {
		w.WriteHeader(204)
		return
	}

	if params.Event == "user.upgraded" {
		err := cfg.db.UpgradeUserToRed(r.Context(), params.Data.UserID)
		if err != nil {
			respondWithError(w, http.StatusNotFound, "user not found", err)
			return
		}
	}

	w.WriteHeader(204)
}
