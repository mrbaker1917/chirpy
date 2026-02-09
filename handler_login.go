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

	token, err := auth.MakeJWT(user.ID, apiCfg.secret)
	if err != nil {
		log.Printf("Error acquiring JWT: %s", err)
		respondWithError(w, 500, "Error acquiring JWT")
		return
	}

	r_token, err := auth.MakeRefreshToken()
	if err != nil {
		log.Printf("Error creating Refresh Token: %s", err)
		respondWithError(w, 500, "Error creating Refresh Token")
		return
	}

	_, err = apiCfg.db.CreateRefreshToken(ctx, database.CreateRefreshTokenParams{
		Token:     r_token,
		UserID:    user.ID,
		ExpiresAt: time.Now().Add(60 * 24 * time.Hour),
	})

	if err != nil {
		log.Printf("Error saving refresh token: %s", err)
		respondWithError(w, 500, "Error saving refresh token")
		return
	}

	type userWithToken struct {
		ID           uuid.UUID `json:"id"`
		CreatedAt    time.Time `json:"created_at"`
		UpdatedAt    time.Time `json:"updated_at"`
		Email        string    `json:"email"`
		Token        string    `json:"token"`
		RefreshToken string    `json:"refresh_token"`
		IsChirpyRed  bool      `json:"is_chirpy_red"`
	}

	respondWithJSON(w, 200, userWithToken{
		ID:           user.ID,
		CreatedAt:    user.CreatedAt,
		UpdatedAt:    user.UpdatedAt,
		Email:        user.Email,
		Token:        token,
		RefreshToken: r_token,
		IsChirpyRed:  user.IsChirpyRed,
	})
}
