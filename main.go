package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"sync/atomic"
	"time"

	"github.com/google/uuid"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/mrbaker1917/chirpy/internal/database"
)

type apiConfig struct {
	fileserverHits atomic.Int32
	db             *database.Queries
	platform       string
}

type User struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email     string    `json:"email"`
}

// helpers:

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	dat, err := json.Marshal(payload)
	if err != nil {
		log.Printf("Error marshalling JSON: %s", err)
		w.WriteHeader(500)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(dat)
}

func respondWithError(w http.ResponseWriter, code int, msg string) {
	payld := make(map[string]string)
	payld["error"] = msg
	respondWithJSON(w, code, payld)
}

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

// handlers:

func handlerValidateChirp(w http.ResponseWriter, r *http.Request) {
	type chirp struct {
		Body string `json:"body"`
	}

	decoder := json.NewDecoder(r.Body)
	chp := chirp{}
	err := decoder.Decode(&chp)
	if err != nil {
		log.Printf("Error decoding chirp: %s", err)
		respondWithError(w, 500, "Error decoding request body")
		return
	}

	type repbody struct {
		CleanedBody string `json:"cleaned_body"`
	}

	if len(chp.Body) > 140 {
		respondWithError(w, 400, "Chirp is too long")
		return
	}

	cleaned_str := sanitizeChirp(chp.Body)

	respondWithJSON(w, 200, repbody{CleanedBody: cleaned_str})
}

func (apiCfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		apiCfg.fileserverHits.Add(1)
		next.ServeHTTP(w, r)
	})
}

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

func (apiCfg *apiConfig) handlerMetrics(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf(`
	<html>
  		<body>
    		<h1>Welcome, Chirpy Admin</h1>
    		<p>Chirpy has been visited %d times!</p>
  		</body>
	</html>`,
		apiCfg.fileserverHits.Load())))
}

func handlerReadiness(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(http.StatusText(http.StatusOK)))
}

func (cfg *apiConfig) handlerCreateUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	type reqEmail struct {
		Email string `json:"email"`
	}

	decoder := json.NewDecoder(r.Body)
	rEmail := reqEmail{}
	err := decoder.Decode(&rEmail)
	if err != nil {
		log.Printf("Error decoding chirp: %s", err)
		respondWithError(w, 500, "Error decoding request body")
		return
	}

	user, err := cfg.db.CreateUser(ctx, rEmail.Email)
	if err != nil {
		log.Printf("Could not create new user: %s", err)
		respondWithError(w, 500, "Error trying to create new user")
		return
	}
	userResp := User{
		ID:        user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Email:     user.Email,
	}
	respondWithJSON(w, 201, userResp)

}

func main() {
	godotenv.Load()
	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		log.Fatal("DB_URL must be set")
	}
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("database failed to open: %s", err)
	}
	dbQueries := database.New(db)
	apiCfg := apiConfig{
		fileserverHits: atomic.Int32{},
		db:             dbQueries,
		platform:       os.Getenv("PLATFORM"),
	}

	const filepathRoot = "."
	const port = "8080"
	mux := http.NewServeMux()
	mux.HandleFunc("GET /admin/metrics", apiCfg.handlerMetrics)
	mux.HandleFunc("POST /admin/reset", apiCfg.handlerReset)
	mux.Handle("/app/", apiCfg.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(filepathRoot)))))
	mux.HandleFunc("GET /api/healthz", handlerReadiness)
	mux.HandleFunc("POST /api/validate_chirp", handlerValidateChirp)
	mux.HandleFunc("POST /api/users", apiCfg.handlerCreateUser)

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	log.Printf("Serving files from %s on port: %s\n", filepathRoot, port)
	log.Fatal(srv.ListenAndServe())
}

// here is model for handlers:
// func (cfg *apiConfig) someHandler(w http.ResponseWriter, r *http.Request) {
//     ctx := r.Context() // <- comes from the request

//     user, err := cfg.db.CreateUser(ctx, email)
//     // handle err, write response, etc.
// }
