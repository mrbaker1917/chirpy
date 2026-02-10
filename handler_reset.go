package main

import (
	"fmt"
	"log"
	"net/http"
)

func (apiCfg *apiConfig) handlerReset(w http.ResponseWriter, r *http.Request) {

	if apiCfg.platform != "dev" {
		respondWithError(w, 403, "forbidden")
		return
	}

	err := apiCfg.db.DeleteAll(r.Context())
	if err != nil {
		log.Printf("error deleting all users: %s", err)
		respondWithError(w, 500, "Error deleting all users")
		return
	}

	apiCfg.fileserverHits.Store(0)
	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf("Number of server hits is reset to %d", apiCfg.fileserverHits.Load())))

}
