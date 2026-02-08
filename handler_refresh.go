package main

import (
	"database/sql"
	"errors"
	"log"
	"net/http"

	"github.com/mrbaker1917/chirpy/internal/auth"
)

func (apiCfg *apiConfig) handlerRefresh(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	r_token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, 401, "Could not find token in header")
		return
	}

	if r_token == "" {
		respondWithError(w, 401, "Could not find token in header")
		return
	}

	user, err := apiCfg.db.GetUserFromRefreshToken(ctx, r_token)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			respondWithError(w, 401, "Unauthorized")
			return
		}
		respondWithError(w, 500, "Error looking up refresh token")
		return
	}

	token, err := auth.MakeJWT(user.ID, apiCfg.secret)
	if err != nil {
		log.Printf("Error acquiring JWT: %s", err)
		respondWithError(w, 500, "Error acquiring JWT")
		return
	}

	type resp struct {
		Token string `json:"token"`
	}

	respondWithJSON(w, 200, resp{Token: token})
}
