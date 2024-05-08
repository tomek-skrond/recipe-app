package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
)

func main() {

	listenPort := ":8000"

	workdir, err := os.Getwd()
	if err != nil {
		log.Println(err)
	}
	log.Println(workdir)

	db, err := NewStorage()
	if err != nil {
		log.Fatalln(err)
	}
	db.SeedRecipes()
	s, err := NewAPIServer(db, listenPort)
	if err != nil {
		log.Fatalln(err)
	}
	s.Start()
}

type apiFunc func(http.ResponseWriter, *http.Request) error

func writeJSON(w http.ResponseWriter, statusCode int, v any) error {
	w.WriteHeader(statusCode)
	w.Header().Add("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(v)
}

func writeHTTP(w http.ResponseWriter, statusCode int, v any) {
	w.WriteHeader(statusCode)
	w.Header().Add("Content-Type", "text/html")
	w.Header().Add("Content-Type", "multipart/form-data")
}

func makeHandler(f apiFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := f(w, r); err != nil {
			if e, ok := err.(apiError); !ok {
				writeJSON(w, e.statuscode, err)
				return
			}
			writeJSON(w, http.StatusInternalServerError, err)
		}
	}
}

func makeFrontendHandler(f apiFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := f(w, r); err != nil {
			// Handle API errors
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

type apiError struct {
	message    string
	statuscode int
}

func (e apiError) Error() string {
	return e.message
}
