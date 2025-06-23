package main

import (
	"log"
	"net/http"
)

func (cfg *apiConfig) handlerReset(w http.ResponseWriter, r *http.Request) {
	if cfg.platform != "dev" {
		log.Println("Reset user action not allowed in non-dev environment")
		w.WriteHeader(403)
		return
	}

	err := cfg.db.ResetUsers(r.Context())
	if err != nil {
		log.Printf("error resetting users: %s", err)
		w.WriteHeader(500)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("All users reset"))
}
