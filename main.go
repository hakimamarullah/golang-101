package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

// Router
const (
	MOVIES    = "/movies"
	MOVIES_ID = "/movies/{id}"
)

// Header
const (
	APPLICATION_JSON = "application/json"
)

// Model
type Movie struct {
	ID       int       `json:"id"`
	Title    string    `json:"title"`
	Year     int       `json:"year"`
	Director *Director `json:"director"`
}

type Director struct {
	ID        int    `json:"id"`
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
}

type ResponseBody struct {
	Message string      `json:"message"`
	Code    int         `json:"code"`
	Total   int         `json:"total"`
	Data    interface{} `json:"data"`
}

// Storage
var movies []Movie

// Handler
func Home(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "Hello, Welcome to GoMovies API!")
}

func getMovies(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", APPLICATION_JSON)
	json.NewEncoder(w).Encode(ResponseBody{Message: "Data retrieved", Total: len(movies), Data: movies, Code: http.StatusOK})
}

func getMovieById(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", APPLICATION_JSON)
	params := mux.Vars(r)

	id, _ := strconv.Atoi(params["id"])
	for _, item := range movies {
		if item.ID == id {
			json.NewEncoder(w).Encode(item)
			return
		}
	}

	w.WriteHeader(http.StatusNotFound)
	json.NewEncoder(w).Encode(ResponseBody{Message: fmt.Sprintf("Movie with id %d not found", id), Code: http.StatusNotFound})
}

func deleteMovieById(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", APPLICATION_JSON)

	params := mux.Vars(r)
	id, _ := strconv.Atoi(params["id"])

	for index, item := range movies {
		if item.ID == id {
			movies = append(movies[:index], movies[index+1:]...)
			json.NewEncoder(w).Encode(ResponseBody{Message: fmt.Sprintf("Movie %d deleted", id), Code: http.StatusOK})
			return
		}
	}

	w.WriteHeader(http.StatusNotFound)
	json.NewEncoder(w).Encode(ResponseBody{Message: fmt.Sprintf("Movie with id %d not found", id), Code: http.StatusNotFound})
}

func updateMovieById(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", APPLICATION_JSON)
	var newMovie Movie
	_ = json.NewDecoder(r.Body).Decode(&newMovie)
	for index, item := range movies {
		if item.ID == newMovie.ID {
			movies = append(movies[:index], movies[index+1:]...)
			movies = append(movies, newMovie)
			json.NewEncoder(w).Encode(ResponseBody{Message: fmt.Sprintf("Movie %d updated", item.ID), Code: http.StatusOK})
			return
		}
	}
	w.WriteHeader(http.StatusNotFound)
	json.NewEncoder(w).Encode(ResponseBody{Message: fmt.Sprintf("Movie with id %d not found", newMovie.ID), Code: http.StatusNotFound})
}

func postMovie(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", APPLICATION_JSON)

	var movie Movie
	err := json.NewDecoder(r.Body).Decode(&movie)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ResponseBody{Message: "Invalid Data, Not Created", Code: http.StatusBadRequest})
		return
	}

	movie.ID = generateID()
	movies = append(movies, movie)

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(ResponseBody{Message: fmt.Sprintf("Movie added ID: %d", movie.ID), Code: http.StatusCreated})
}

func generateID() int {
	return rand.Intn(10000000)
}

func main() {
	r := mux.NewRouter()
	movies = append(movies, Movie{ID: 1, Title: "Koala Kumal", Year: 2016, Director: &Director{ID: 1, Firstname: "Raditya", Lastname: "Dika"}})
	movies = append(movies, Movie{ID: 2, Title: "Kambing Jantan", Year: 2009, Director: &Director{ID: 1, Firstname: "Raditya", Lastname: "Dika"}})
	movies = append(movies, Movie{ID: 3, Title: "Manusia Setengah Salmon", Year: 2013, Director: &Director{ID: 1, Firstname: "Raditya", Lastname: "Dika"}})
	r.HandleFunc("/", Home).Methods("GET")
	r.HandleFunc(MOVIES, getMovies).Methods(http.MethodGet)
	r.HandleFunc(MOVIES, postMovie).Methods(http.MethodPost)
	r.HandleFunc(MOVIES_ID, getMovieById).Methods(http.MethodGet)
	r.HandleFunc(MOVIES_ID, deleteMovieById).Methods(http.MethodDelete)
	r.HandleFunc(MOVIES, updateMovieById).Methods(http.MethodPut)

	fmt.Printf("Server up and running on port 8000\n")
	log.Fatal(http.ListenAndServe("localhost:8000", r))
}
