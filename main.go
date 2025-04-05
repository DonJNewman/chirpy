package main

import (
	"chirpy/internal/database"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync/atomic"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type apiConfig struct {
	fileserverHits atomic.Int32
	dbQueries      *database.Queries
}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileserverHits.Add(1)
		next.ServeHTTP(w, r)

	})
}

// REVIEW THIS CODE TOMORROW
func (cfg *apiConfig) metricsHandler(w http.ResponseWriter, r *http.Request) {
	// Get the current hits count
	hits := cfg.fileserverHits.Load()

	// Write the response
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Hits: %d", hits)
}

func (cfg *apiConfig) resetHandler(w http.ResponseWriter, r *http.Request) {
	// Reset the hits counter to 0
	cfg.fileserverHits.Store(0)

	// Write a simple response
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Reset successful"))
}

func main() {
	godotenv.Load()
	dbURL := os.Getenv("DB_URL")
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		return
	}
	dbQueries := database.New(db)
	//create apiConfig instance
	apiCfg := &apiConfig{
		dbQueries: dbQueries, // Pass in the initialized database queries
	}
	// 1. Create the multiplexer
	mux := http.NewServeMux()

	// 2. Create the file server
	fileServer := http.FileServer(http.Dir("."))

	//2.5 handler
	handler := http.StripPrefix("/app", fileServer)

	// 3. Register the file server as the handler for /app, strip app for local files
	mux.Handle("/app/", apiCfg.middlewareMetricsInc(handler))

	//4. Health check route, return OK and 200, doesn't need file server because we're not serving a file
	mux.HandleFunc("GET /api/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))

	})
	//admin metrics page with html formatting
	mux.HandleFunc("GET /admin/metrics", func(w http.ResponseWriter, r *http.Request) {
		htmlTemp := `<html>
							<body>
								<h1>Welcome, Chirpy Admin</h1>
								<p>Chirpy has been visited %d times!</p>
							</body>
							</html>`
		formattedHtml := fmt.Sprintf(htmlTemp, apiCfg.fileserverHits.Load())
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(formattedHtml))

	})
	//post TO API VALIDATE_CHIRP
	mux.HandleFunc("POST /api/validate_chirp", func(w http.ResponseWriter, r *http.Request) {

		//THINK I NEED TYPE PARAMETERS STRUCT HERE
		type parameters struct {
			Body string `json:"body"`
		}
		decoder := json.NewDecoder(r.Body)
		params := parameters{}
		err := decoder.Decode(&params)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"error": "Something went wrong"})
			return
		}
		//if len of the body is over 140 respond statuscode 400 and chirp is too long
		if len(params.Body) > 140 {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"error": "Chirp is too long"})
			return
		}

		//else code 200 "valid":true
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		respbody, err := json.Marshal((map[string]bool{"valid": true}))
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Write(respbody)

	})

	//mux.HandleFunc("GET /api/metrics", apiCfg.metricsHandler)
	mux.HandleFunc("POST /admin/reset", apiCfg.resetHandler)

	// 4. Create and start the server
	s := http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	// 5. Start listening
	err = s.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}
