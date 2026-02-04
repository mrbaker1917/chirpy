package main

import (
	"log"
	"net/http"

	"github.com/google/uuid"
)

func (apiCfg *apiConfig) handlerGetChirps(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	chirps, err := apiCfg.db.GetChirps(ctx)
	if err != nil {
		log.Printf("we encountered an error: %s", err)
	}
	response_chirps := []ChirpResponse{}
	for _, chirp := range chirps {
		response_chirps = append(response_chirps, ChirpResponse{
			ID:        chirp.ID,
			CreatedAt: chirp.CreatedAt,
			UpdatedAt: chirp.UpdatedAt,
			Body:      chirp.Body,
			UserID:    chirp.UserID,
		})
	}
	respondWithJSON(w, 200, response_chirps)
}

func (apiCfg *apiConfig) handlerGetChirpById(w http.ResponseWriter, r *http.Request) {
	chirpId := r.PathValue("chirpID")

	uChirpId, err := uuid.Parse(chirpId)
	if err != nil {
		log.Fatalf("failed to parse UUID %q: %v", chirpId, err)
	}
	if chirpId == "" {
		log.Println("No ChirpID found. We need a ChirpID to get the Chirp.")
	}

	chirp, err := apiCfg.db.GetChirpById(r.Context(), uChirpId)
	if err != nil {
		respondWithError(w, 404, "Chirp not found.")
	}

	respondWithJSON(w, 200, ChirpResponse{
		ID:        chirp.ID,
		CreatedAt: chirp.CreatedAt,
		UpdatedAt: chirp.UpdatedAt,
		Body:      chirp.Body,
		UserID:    chirp.UserID,
	})

}
