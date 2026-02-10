package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/mrbaker1917/chirpy/internal/auth"
	"github.com/mrbaker1917/chirpy/internal/database"
)

func (apiCfg *apiConfig) handlerCreateUser(w http.ResponseWriter, r *http.Request) {
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

	hashed_password, err := auth.HashPassword(reqBdy.Password)
	if err != nil {
		log.Printf("Error hasing password: %s", err)
	}

	user, err := apiCfg.db.CreateUser(ctx, database.CreateUserParams{
		Email:          reqBdy.Email,
		HashedPassword: hashed_password,
	})

	if err != nil {
		log.Printf("Could not create new user: %s", err)
		respondWithError(w, 500, "Error trying to create new user")
		return
	}
	userResp := User{
		ID:          user.ID,
		CreatedAt:   user.CreatedAt,
		UpdatedAt:   user.UpdatedAt,
		Email:       user.Email,
		IsChirpyRed: user.IsChirpyRed,
	}
	respondWithJSON(w, 201, userResp)

}
