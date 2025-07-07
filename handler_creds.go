package main

import (
	"encoding/json"
	"net/http"

	"github.com/scottyloveless/chirpy/internal/auth"
	"github.com/scottyloveless/chirpy/internal/database"
)

func (cfg *apiConfig) handlerUpdateCredentials(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Password string `json:"password"`
		Email    string `json:"email"`
	}
	type response struct {
		User
	}
	bearerToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "invalid token", err)
		return
	}

	userUUID, err := auth.ValidateJWT(bearerToken, cfg.secret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "invalid token", err)
		return
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err = decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "error decoding json", err)
		return
	}

	hash, err := auth.HashPassword(params.Password)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "error hashing password", err)
	}

	updatedDbUser, err := cfg.db.UpdateUserCredentials(r.Context(), database.UpdateUserCredentialsParams{
		HashedPassword: hash,
		Email:          params.Email,
		ID:             userUUID,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "error updating user", err)
	}

	respondWithJSON(w, http.StatusOK, response{
		User: User{
			ID:        updatedDbUser.ID,
			Email:     updatedDbUser.Email,
			CreatedAt: updatedDbUser.CreatedAt,
			UpdatedAt: updatedDbUser.UpdatedAt,
		},
	})
}
