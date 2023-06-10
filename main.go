package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

// Router
const (
	MOVIES    = "/movies"
	MOVIES_ID = "/movies/{id}"
)

// Header
const ()

func Home(w http.ResponseWriter, r *http.Request) {
	data := mux.Vars(r)
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Hello %v\n", data["name"])
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/", Home).Methods("GET")
	r.HandleFunc(MOVIES, getMovies).Methods(http.MethodGet)
	r.Handle(MOVIES, postMovie).Methods(http.MethodPost)
	r.HandleFunc(MOVIES_ID, getMovieById).Methods(http.MethodGet)
	r.HandleFunc(MOVIES_ID, deleteMovieById).Methods(http.MethodDelete)
	r.HandleFunc(MOVIES_ID, updateMovieById).Methods(http.MethodPut)

	log.Fatal("Server up and running on port", http.ListenAndServe(":8000", r))
}
