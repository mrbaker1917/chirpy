package main

import (
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/mrbaker1917/chirpy/internal/auth"
)

func (apiCfg *apiConfig) handlerDeleteChirp(w http.ResponseWriter, r *http.Request) {
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

	userID, err := auth.ValidateJWT(token, apiCfg.secret)
	if err != nil {
		log.Printf("Error validating access token %s", err)
		respondWithError(w, 401, "Error validating access token")
		return
	}

	chirpId := r.PathValue("chirpID")
	if chirpId == "" {
		log.Println("No chirpID found. We need a chirpID to get the Chirp.")
		respondWithError(w, 401, "We need a chirpID to delete it.")
		return
	}

	uChirpId, err := uuid.Parse(chirpId)
	if err != nil {
		log.Fatalf("failed to parse UUID %q: %v", chirpId, err)
		return
	}

	chirp, err := apiCfg.db.GetChirpById(r.Context(), uChirpId)
	if err != nil {
		respondWithError(w, 404, "Chirp not found.")
		return
	}

	if userID != chirp.UserID {
		respondWithError(w, 403, "Not your chirp, so cannot delete it.")
		return
	}

	err = apiCfg.db.DeleteChirpById(ctx, chirp.ID)
	if err != nil {
		respondWithError(w, 403, "Chirp could not be deleted.")
		return
	}

	w.WriteHeader(204)

}
