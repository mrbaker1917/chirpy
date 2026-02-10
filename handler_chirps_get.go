package main

import (
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/mrbaker1917/chirpy/internal/database"
)

func (apiCfg *apiConfig) handlerGetChirps(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	chirps := []database.Chirp{}
	var err error
	authorID := r.URL.Query().Get("author_id")
	if authorID != "" {
		uAuthorID, err := uuid.Parse(authorID)
		if err != nil {
			log.Printf("We encountered an error parsing authorID: %s", err)
			respondWithError(w, 501, "we encountered an error parsing author_id")
			return
		}
		chirps, err = apiCfg.db.GetChirpsByAuthor(ctx, uAuthorID)
		if err != nil {
			log.Printf("We encountered grabbing chirps by this authorID: %s", err)
			respondWithError(w, 501, "We could not find any chirps for this author_id")
			return
		}

	} else {

		chirps, err = apiCfg.db.GetChirps(ctx)
		if err != nil {
			log.Printf("we encountered an error: %s", err)
			respondWithError(w, 501, "we encountered an error getting chirps")
			return
		}

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
		return
	}
	if chirpId == "" {
		log.Println("No ChirpID found. We need a ChirpID to get the Chirp.")
		return
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
