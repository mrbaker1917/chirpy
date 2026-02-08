package main

import (
	"net/http"

	"github.com/mrbaker1917/chirpy/internal/auth"
)

func (apiCfg *apiConfig) handlerRevoke(w http.ResponseWriter, r *http.Request) {
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

	err = apiCfg.db.RevokeRefreshToken(ctx, r_token)
	if err != nil {
		respondWithError(w, 500, "Error revoking refresh token")
		return
	}

	w.WriteHeader(http.StatusNoContent)

}
