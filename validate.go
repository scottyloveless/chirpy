package main

import (
	"encoding/json"
	"log"
	"net/http"
)

type input struct {
	Body string `json:"body"`
}

type validResp struct {
	Valid bool `json:"valid"`
}

type invalidResp struct {
	Error string `json:"error"`
}

func (cfg *apiConfig) handlerValidate(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	params := input{}
	wentWrong := invalidResp{Error: "Something went wrong"}
	wrongDat, _ := json.Marshal(wentWrong)

	err := decoder.Decode(&params)
	if err != nil {
		log.Printf("Error decoding parameters: %s", err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(500)
		w.Write(wrongDat)
		return
	}

	if len(params.Body) > 140 {
		log.Printf("Chirp is too long: %v", len(params.Body))

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

	respBody := validResp{Valid: true}

	dat, err := json.Marshal(respBody)
	if err != nil {
		log.Printf("Error marshalling JSON after valid chirp: %v", err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(500)
		w.Write(wrongDat)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	w.Write(dat)
}
