package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync/atomic"

	"github.com/joho/godotenv"
	"github.com/justinyeh1995/go-web-server-template/internal/database"
	_ "github.com/lib/pq"
)

type apiConfig struct {
	fileserverHits atomic.Int32
	db             *database.Queries
}

func main() {
	godotenv.Load()

	dbURL := os.Getenv("DB_URL")
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("Failed to open db connection, with error msg %v", err)
	}
	dbQueries := database.New(db)

	const filepathRoot = "."
	const port = "8080"

	apiCfg := &apiConfig{
		fileserverHits: atomic.Int32{},
		db:             dbQueries,
	}

	mux := http.NewServeMux()
	mux.Handle("/app/", apiCfg.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(filepathRoot)))))
	mux.HandleFunc("GET /api/healthz", healthCheckHandler)
	mux.HandleFunc("GET /admin/metrics", apiCfg.handlerMetrics)
	mux.HandleFunc("POST /admin/reset", apiCfg.reset)
	mux.HandleFunc("POST /api/validate_chirp", apiCfg.handlerCheckChirpy)

	server := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	log.Printf("Serving files from %s on port: %s\n", filepathRoot, port)
	log.Fatal(server.ListenAndServe())
}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileserverHits.Add(1)
		next.ServeHTTP(w, r)
	})
}

func (cfg *apiConfig) handlerMetrics(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/html; charset=utf-8;")
	w.WriteHeader(http.StatusOK)

	msg := fmt.Sprintf(`
			<html>
				<body>
					<h1>Welcome, Chirpy Admin</h1>
					<p>Chirpy has been visited %d times!</p>
				</body>
			</html>`, cfg.fileserverHits.Load())
	w.Write([]byte(msg))
}

func (cfg *apiConfig) handlerCheckChirpy(w http.ResponseWriter, r *http.Request) {
	type param struct {
		Body string `json:"body"`
	}
	type errorResponse struct {
		Error string `json:"error"`
	}
	type regularReponse struct {
		Valid bool `json:"valid"`
	}
	decoder := json.NewDecoder(r.Body)
	params := param{}
	err := decoder.Decode(&params)

	if err != nil {
		// an error will be thrown if the JSON is invalid or has the wrong types
		// any missing fields will simply have their values in the struct set to their zero value
		log.Printf("Error decoding parameters: %s", err)
		w.WriteHeader(500)
		return
	}

	if len(params.Body) > 140 {
		w.WriteHeader(http.StatusBadRequest)
		errMessage := errorResponse{
			Error: "Chirp is too long",
		}
		dat, err := json.Marshal(errMessage)
		if err != nil {
			log.Printf("Error marshalling JSON: %s", err)
			w.WriteHeader(500)
			return
		}
		w.Write(dat)
		return
	}

	respBody := regularReponse{
		Valid: true,
	}

	dat, err := json.Marshal(respBody)
	if err != nil {
		log.Printf("Error marshalling JSON: %s", err)
		w.WriteHeader(500)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	w.Write(dat)

}
