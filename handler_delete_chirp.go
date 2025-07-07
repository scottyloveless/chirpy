package main

import (
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/scottyloveless/chirpy/internal/auth"
)

func (cfg *apiConfig) handlerDeleteChirp(w http.ResponseWriter, r *http.Request) {
	chirpIDString := r.PathValue("chirpID")
	if chirpIDString == "" {
		respondWithError(w, http.StatusBadRequest, "no chirpID provided", fmt.Errorf(""))
		return
	}

	bearerToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "invalid jwt token", err)
		return
	}

	validUUID, err := auth.ValidateJWT(bearerToken, cfg.secret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "problem validating jwt", err)
		return
	}

	parsedId, err := uuid.Parse(chirpIDString)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid chirp ID", err)
		return
	}

	dbChirp, err := cfg.db.GetChirp(r.Context(), parsedId)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "error fetching chirp", err)
		return
	}

	if dbChirp.UserID != validUUID {
		respondWithError(w, http.StatusForbidden, "user not authorized to delete this chirp", err)
		return
	}

	err = cfg.db.DeleteChirpByID(r.Context(), dbChirp.ID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "error deleting chirp by id", err)
		return
	}

	w.WriteHeader(204)
}
