package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/mrbaker1917/chirpy/internal/auth"
	"github.com/mrbaker1917/chirpy/internal/database"
)

// struct:

type ChirpResponse struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Body      string    `json:"body"`
	UserID    uuid.UUID `json:"user_id"`
}

// helpers:

func sanitizeChirp(s string) string {
	profane_words := []string{"kerfuffle", "sharbert", "fornax"}
	s_slice := strings.Split(s, " ")
	for i, w := range s_slice {
		for _, p := range profane_words {
			if strings.ToLower(w) == p {
				s_slice[i] = "****"
			}
		}
	}
	new_str := strings.Join(s_slice, " ")

	return new_str
}

func validateChirp(s string) (string, error) {

	if len(s) > 140 {
		return "", fmt.Errorf("Chirp is too long!")
	}

	cleaned_str := sanitizeChirp(s)
	return cleaned_str, nil
}

// handler:
func (apiCfg *apiConfig) handlerCreateChirp(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	type reqChirp struct {
		Body string `json:"body"`
	}

	decoder := json.NewDecoder(r.Body)
	chp := reqChirp{}
	err := decoder.Decode(&chp)
	if err != nil {
		log.Printf("Error decoding chirp: %s", err)
		respondWithError(w, 500, "Error decoding request body")
		return
	}

	bearerToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, 401, "Requester does not have a session token")
		return
	}

	userID, err := auth.ValidateJWT(bearerToken, apiCfg.secret)
	if err != nil {
		respondWithError(w, 401, "Session token not valid!")
		return
	}

	cleanedBody, err := validateChirp(chp.Body)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	chirp, err := apiCfg.db.CreateChirp(
		ctx,
		database.CreateChirpParams{
			Body:   cleanedBody,
			UserID: userID,
		},
	)

	if err != nil {
		log.Printf("Could not create chirp: %s", err)
		respondWithError(w, http.StatusInternalServerError, "Could not create chirp")
		return
	}

	respondWithJSON(w, http.StatusCreated, ChirpResponse{
		ID:        chirp.ID,
		CreatedAt: chirp.CreatedAt,
		UpdatedAt: chirp.UpdatedAt,
		Body:      chirp.Body,
		UserID:    chirp.UserID,
	})
}
