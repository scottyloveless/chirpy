package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"time"

	"github.com/scottyloveless/chirpy/internal/auth"
)

func (cfg *apiConfig) handlerRefreshToken(w http.ResponseWriter, r *http.Request) {
	type response struct {
		Token string `json:"token"`
	}
	authToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "error getting token", err)
		return
	}

	dbToken, err := cfg.db.CheckRefreshToken(r.Context(), authToken)
	if err == sql.ErrNoRows {
		respondWithError(w, http.StatusUnauthorized, "token not found", err)
		return
	} else if err != nil {
		respondWithError(w, http.StatusInternalServerError, "error querying database for token", err)
		return
	}

	if dbToken.RevokedAt.Valid {
		respondWithError(w, http.StatusUnauthorized, "token revoked", err)
		return
	}

	if time.Now().After(dbToken.ExpiresAt) {
		respondWithError(w, http.StatusUnauthorized, "token expired", fmt.Errorf("%v", dbToken.ExpiresAt))
		return
	}

	dbUser, err := cfg.db.GetUserFromRefreshToken(r.Context(), dbToken.Token)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "not authorized", err)
		return
	}

	jwt, err := auth.MakeJWT(dbUser, cfg.secret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "not authorized", err)
		return
	}

	respondWithJSON(w, http.StatusOK, response{
		Token: jwt,
	})
}
