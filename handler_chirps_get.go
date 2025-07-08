package main

import (
	"database/sql"
	"errors"
	"net/http"

	"github.com/google/uuid"
	"github.com/scottyloveless/chirpy/internal/database"
)

func (cfg *apiConfig) handlerChirpsRetrieve(w http.ResponseWriter, r *http.Request) {
	var dbChirps []database.Chirp
	var err error
	author_string := r.URL.Query().Get("author_id")
	if author_string != "" {
		author_id, parseErr := uuid.Parse(author_string)
		if parseErr != nil {
			respondWithError(w, http.StatusBadRequest, "invalid uuid", parseErr)
			return
		}
		dbChirps, err = cfg.db.GetChirpsByAuthor(r.Context(), author_id)
	} else {
		dbChirps, err = cfg.db.GetChirps(r.Context())
	}
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't retrieve chirps", err)
		return
	}
	chirps := []Chirp{}

	for _, dbChirps := range dbChirps {
		chirps = append(chirps, Chirp{
			ID:        dbChirps.ID,
			CreatedAt: dbChirps.CreatedAt,
			UpdatedAt: dbChirps.UpdatedAt,
			Body:      dbChirps.Body,
			UserID:    dbChirps.UserID,
		})
	}

	respondWithJSON(w, http.StatusOK, chirps)
}

func (cfg *apiConfig) handlerChirpsGet(w http.ResponseWriter, r *http.Request) {
	chirpIDString := r.PathValue("chirpID")
	if chirpIDString == "" {
		respondWithError(w, http.StatusBadRequest, "no id provided", nil)
		return
	}
	parsedId, err := uuid.Parse(chirpIDString)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid chirp ID", err)
		return
	}
	dbChirp, err := cfg.db.GetChirp(r.Context(), parsedId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			respondWithError(w, http.StatusNotFound, "No chirp exists with that id", err)
			return

		}
		respondWithError(w, http.StatusInternalServerError, "Couldn't get chirp", err)
		return
	}

	respondWithJSON(w, http.StatusOK, Chirp{
		ID:        dbChirp.ID,
		CreatedAt: dbChirp.CreatedAt,
		UpdatedAt: dbChirp.UpdatedAt,
		Body:      dbChirp.Body,
		UserID:    dbChirp.UserID,
	})
}
