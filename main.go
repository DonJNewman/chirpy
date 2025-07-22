package main

import (
	"chirpy/internal/database"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync/atomic"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

// hold stateful memory data
type apiConfig struct {
	fileserverHits atomic.Int32 //atomic.int32 provides safe concurrent editing
	dbQueries      *database.Queries
}

func main() {
	godotenv.Load()                        //load dotenv variables
	dbURL := os.Getenv("DB_URL")           //grab DBurl
	db, err := sql.Open("postgres", dbURL) //opens postgres connection
	if err != nil {
		return
	}
	dbQueries := database.New(db) //start new db for queries
	//create apiConfig instance
	apiCfg := &apiConfig{
		dbQueries: dbQueries, // Pass in the initialized database queries
	}
	// 1. Create the multiplexer
	mux := http.NewServeMux()

	// 2. Create the file server
	fileServer := http.FileServer(http.Dir("."))
	//the http.Dir creates a file system to serve from current directory. Implements http.FileSystem
	//basically a converter

	//fileserver creates an HTTP HANDLER that can serve files.
	//recieves requests, parses the request (searches for files), serves the response, handle MiME types, directory listings ect

	//2.5 handler
	handler := http.StripPrefix("/app", fileServer)

	// 3. Register the file server as the handler for /app, strip app for local files
	mux.Handle("/app/", apiCfg.middlewareMetricsInc(handler))

	//4. Health check route, return OK and 200, doesn't need file server because we're not serving a file
	mux.HandleFunc("GET /api/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8") //always gotta set the header
		w.WriteHeader(http.StatusOK)                                //http status 200
		w.Write([]byte("OK"))                                       //w.Write excepts slice of bytes

	})
	//admin metrics page with html formatting
	mux.HandleFunc("GET /admin/metrics", func(w http.ResponseWriter, r *http.Request) {
		htmlTemp := `<html>
							<body>
								<h1>Welcome, Chirpy Admin</h1>
								<p>Chirpy has been visited %d times!</p>
							</body>
							</html>` //This is so cool! a mini html file :D
		formattedHtml := fmt.Sprintf(htmlTemp, apiCfg.fileserverHits.Load()) //load is the concurrent safety read from atomic.Int32
		//SprintF is taking the tempHTML and replacing the variable with the actual data
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(formattedHtml)) //expects slice of bytes

	})
	//post TO API VALIDATE_CHIRP
	mux.HandleFunc("POST /api/validate_chirp", apiCfg.handlerValidateChirp)

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
