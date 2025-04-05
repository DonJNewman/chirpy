package main

import (
	"log"
	"net/http"
)

func main() {
	mux := http.NewServeMux()
	s := http.Server{
		Addr:    ":8080",
		Handler: mux,
	}
	err := s.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}
