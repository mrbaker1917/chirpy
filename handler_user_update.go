package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/mrbaker1917/chirpy/internal/auth"
	"github.com/mrbaker1917/chirpy/internal/database"
)

func (apiCfg *apiConfig) handlerUserUpdate(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, 401, "Could not find token in header")
		return
	}

	if token == "" {
		respondWithError(w, 401, "Could not find token in header")
		return
	}

	type reqBody struct {
		Password string `json:"password"`
		Email    string `json:"email"`
	}

	decoder := json.NewDecoder(r.Body)
	reqBdy := reqBody{}
	err = decoder.Decode(&reqBdy)
	if err != nil {
		log.Printf("Error decoding user input: %s", err)
		respondWithError(w, 500, "Error decoding request body")
		return
	}

	hashed_password, err := auth.HashPassword(reqBdy.Password)
	if err != nil {
		log.Printf("Error hashing password: %s", err)
		respondWithError(w, 500, "Error hasing password")
		return
	}

	userID, err := auth.ValidateJWT(token, apiCfg.secret)
	if err != nil {
		log.Printf("Error validating access token %s", err)
		respondWithError(w, 401, "Error validating access token")
		return
	}

	updated_user, err := apiCfg.db.UpdateUser(ctx, database.UpdateUserParams{
		ID:             userID,
		Email:          reqBdy.Email,
		HashedPassword: hashed_password,
	})

	if err != nil {
		log.Printf("Error updating user: %s", err)
		respondWithError(w, 500, "Error updating user")
		return
	}

	type updatedUser struct {
		ID        uuid.UUID `json:"id"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
		Email     string    `json:"email"`
		Token     string    `json:"token"`
	}

	respondWithJSON(w, 200, updatedUser{
		ID:        updated_user.ID,
		CreatedAt: updated_user.CreatedAt,
		UpdatedAt: updated_user.UpdatedAt,
		Email:     updated_user.Email,
		Token:     token,
	})

}
