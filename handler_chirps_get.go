package main

import (
	"net/http"
)

func (cfg *apiConfig) handlerChirpsRetrieve(w http.ResponseWriter, r *http.Request) {
	dbChirps, err := cfg.db.AllChirps(r.Context())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't retrieve chirps", err)
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
