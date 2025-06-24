package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/scottyloveless/chirpy/internal/database"
)

type Chirp struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Body      string    `json:"body"`
	UserID    uuid.UUID `json:"user_id"`
}

func (cfg *apiConfig) handlerAddChirp(w http.ResponseWriter, r *http.Request) {
	type newChirpRequestParams struct {
		Body   string `json:"body"`
		UserID string `json:"user_id"`
	}

	decoder := json.NewDecoder(r.Body)
	unvalidatedNewChirp := newChirpRequestParams{}

	wentWrong := invalidResp{Error: "Something went wrong"}
	wrongDat, _ := json.Marshal(wentWrong)

	err := decoder.Decode(&unvalidatedNewChirp)
	if err != nil {
		log.Printf("Error decoding parameters: %s", err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(500)
		w.Write(wrongDat)
		return
	}

	if len(unvalidatedNewChirp.Body) == 0 {
		log.Println("Chirp cannot be blank.")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(400)
		w.Write(wrongDat)
		return
	}

	if len(unvalidatedNewChirp.Body) > 140 {
		log.Printf("Chirp is too long: %v", len(unvalidatedNewChirp.Body))

		respBody := invalidResp{Error: "Chirp is too long"}

		dat, err := json.Marshal(respBody)
		if err != nil {
			log.Printf("Error marshalling JSON after chirp too long: %v", err)
			w.WriteHeader(500)
			w.Write(wrongDat)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(400)
		w.Write(dat)
		return
	}

	validUUID, err := uuid.Parse(unvalidatedNewChirp.UserID)
	if err != nil {
		log.Printf("user_id is not valid: %v", err)
		w.WriteHeader(400)
		return
	}

	santizedChirp := sanitizeProfanity(unvalidatedNewChirp.Body)

	chirpParams := database.CreateChirpParams{
		Body:   santizedChirp,
		UserID: validUUID,
	}

	newChirp, err := cfg.db.CreateChirp(r.Context(), chirpParams)
	if err != nil {
		log.Printf("error adding new chirp to database: %v", err)
		w.WriteHeader(500)
		return
	}

	newChirpJsonPrep := Chirp{
		ID:        newChirp.ID,
		CreatedAt: newChirp.CreatedAt,
		UpdatedAt: newChirp.UpdatedAt,
		Body:      newChirp.Body,
		UserID:    newChirp.UserID,
	}

	marshalledChirp, err := json.Marshal(newChirpJsonPrep)
	if err != nil {
		log.Printf("error marshalling new chirp: %v", err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(201)
	w.Write(marshalledChirp)
}
