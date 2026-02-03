package main

import (
	"fmt"
	"net/http"
)

func (apiCfg *apiConfig) handlerGetChirps(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	chirps, err := apiCfg.db.GetChirps(ctx)
	if err != nil {
		fmt.Errorf("we encountered an error: %w", err)
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
