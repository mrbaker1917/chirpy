package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/mrbaker1917/chirpy/internal/auth"
)

func (apiCfg *apiConfig) handlerLogin(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	type reqBody struct {
		Password string `json:"password"`
		Email    string `json:"email"`
	}

	decoder := json.NewDecoder(r.Body)
	reqBdy := reqBody{}
	err := decoder.Decode(&reqBdy)
	if err != nil {
		log.Printf("Error decoding user input: %s", err)
		respondWithError(w, 500, "Error decoding request body")
		return
	}

	if len(reqBdy.Email) < 5 || len(reqBdy.Password) < 5 {
		respondWithError(w, 401, "Incorrect email or password")
		return
	}

	user, err := apiCfg.db.GetUserByEmail(ctx, reqBdy.Email)
	if err != nil {
		respondWithError(w, 401, "Incorrect email or password")
		return
	}

	valid, err := auth.CheckPasswordHash(reqBdy.Password, user.HashedPassword)
	if err != nil {
		respondWithError(w, 401, "Incorrect email or password")
		return
	}

	if !valid {
		respondWithError(w, 401, "Incorrect email or password")
		return
	}
	respondWithJSON(w, 200, User{
		ID:        user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Email:     user.Email,
	})
}
