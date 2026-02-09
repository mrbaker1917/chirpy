package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/google/uuid"
)

func (apiCfg *apiConfig) handlerUpgradeUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	type reqBody struct {
		Event string `json:"event"`
		Data  struct {
			UserID string `json:"user_id"`
		} `json:"data"`
	}

	decoder := json.NewDecoder(r.Body)
	reqBdy := reqBody{}
	err := decoder.Decode(&reqBdy)
	if err != nil {
		log.Printf("Error decoding user input: %s", err)
		respondWithError(w, 500, "Error decoding request body")
		return
	}

	if reqBdy.Event != "user.upgraded" {
		respondWithJSON(w, 204, "wrong event!")
		return
	}

	userID, err := uuid.Parse(reqBdy.Data.UserID)
	if err != nil {
		log.Fatal("Failed to parse userID to uuid.")
		return
	}

	err = apiCfg.db.UpgradeUser(ctx, userID)
	if err != nil {
		respondWithError(w, 404, "User could not be found.")
		return
	}

	respondWithJSON(w, 204, struct{}{})
}
