package main

import (
	"fmt"
	"net/http"
)

// review this ch2l1
func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler { //middleware for safely incrementing the hit counter

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileserverHits.Add(1) //access stateful fileserverHits and increment
		next.ServeHTTP(w, r)      //where does this next come from ? is this the next request

	})
}

// REVIEW THIS CODE TOMORROW
func (cfg *apiConfig) metricsHandler(w http.ResponseWriter, r *http.Request) {
	// Get the current hits count
	hits := cfg.fileserverHits.Load() //load is the safe method I was commenting about a cpl days ago, believe this is from http? check that

	// Write the response
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Hits: %d", hits)
}
