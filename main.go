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


type Movie struct{
	ID string `json:"id"`
	Isbn string `json:"isbn"`
	Title string `json:"title"`
	Director *Director `json:"director"`

}



type Director struct{
	Firstname string `json:"firstname"`
	Lastname string `json:"lastname"`

}

var movies []Movie

func getMovies(w http.ResponseWriter, r *http.Request){
	w.Header().Set("Content-Type","application/json")
	json.NewEncoder(w).Encode(movies)

}

func deleteMovie(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json") // Set the response to JSON
    
    // Get the path parameters from the URL (e.g., /movies/{id})
    params := mux.Vars(r)
    
    // Loop through the movies to find the movie with the matching ID
    for index, item := range movies {
        if item.ID == params["id"] {
            // Remove the movie from the slice
            movies = append(movies[:index], movies[index+1:]...)
            
            // Send a success response and exit the function
            w.WriteHeader(http.StatusOK)
            json.NewEncoder(w).Encode(map[string]string{
                "message": "Movie deleted successfully",
                "id": item.ID,
            })
            return // Return after deleting the movie, no need to continue the loop
        }
    }

    // If no movie with the given ID was found, return a 404 error
    http.Error(w, "Movie not found", http.StatusNotFound)
}

func getMovie(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json") // Set the response content type to JSON

    // Get the path parameters from the URL (e.g., /movies/{id})
    params := mux.Vars(r)

    // Loop through the list of movies
    for _, item := range movies {
        // Check if the current movie ID matches the parameter in the URL
        if item.ID == params["id"] {
            // If a match is found, encode the movie as JSON and send it in the response
            json.NewEncoder(w).Encode(item)
            return // Exit the function once the movie is found and returned
        }
    }

    // If no movie with the provided ID is found, return a 404 error
    http.Error(w, "Movie not found", http.StatusNotFound)
}


func createMovies(w http.ResponseWriter, r *http.Request){
	w.Header().Set("Content-Type","application/json")
	var movie Movie
	_ = json.NewDecoder(r.Body).Decode(&movie)
	movie.ID = strconv.Itoa(rand.Intn(100000000))
	movies = append(movies, movie)
	json.NewEncoder(w).Encode(movies)

}

func updateMovie(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")

    // Get the path parameters from the URL (e.g., /movies/{id})
    params := mux.Vars(r)

    // Create a variable to hold the updated movie data
    var updatedMovie Movie

    // Decode the incoming request body into the updatedMovie struct
    if err := json.NewDecoder(r.Body).Decode(&updatedMovie); err != nil {
        http.Error(w, "Invalid JSON format", http.StatusBadRequest)
        return
    }

    // Loop through the movies slice and find the movie by ID
    for index, item := range movies {
        if item.ID == params["id"] {
            // Remove the old movie from the slice
            movies = append(movies[:index], movies[index+1:]...)

            // Set the ID of the updated movie to the original movie's ID
            updatedMovie.ID = params["id"]

            // Add the updated movie to the slice
            movies = append(movies, updatedMovie)

            // Respond with the updated movie
            json.NewEncoder(w).Encode(updatedMovie)
            return
        }
    }

    // If movie with the given ID was not found
    http.Error(w, "Movie not found", http.StatusNotFound)
}

func main(){
	r := mux.NewRouter()
	movies = append(movies, Movie{ID: "1",Isbn:"438277",Title:"Movie One",Director: &Director{Firstname: "john",Lastname: "doe"} })
	movies = append(movies, Movie{ID: "2",Isbn:"45455",Title:"Movie Two",Director: &Director{Firstname: "ayush",Lastname: "jain"} })
	r.HandleFunc("/movies",getMovies).Methods("GET")
	r.HandleFunc("/movies/{id}",getMovie).Methods("GET")
	r.HandleFunc("/movies",createMovies).Methods("POST")
	r.HandleFunc("/movies/{id}",updateMovie).Methods("PUT")
	r.HandleFunc("/movies/{id}",deleteMovie).Methods("DELETE")

	fmt.Print("Starting server at port 8000\n")
	log.Fatal(http.ListenAndServe(":8000",r))
}


