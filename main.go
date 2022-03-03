package main

import (
	"encoding/json"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"time"
)

var (
	APIKEY        string
	APIURL        string
	MOVIE_ID      string
	RANDOM_NUMBER int
)

type movie struct {
	MovieID     int     `gorm:"column:movie_id;primary_key" json:"id"`
	Title       string  `json:"title"`
	ReleaseDate string  `json:"release_date"`
	Language    string  `json:"original_language"`
	Adult       bool    `json:"adult"`
	Image       string  `json:"poster_path"`
	Overview    string  `gorm:"type:varchar(1000)" json:"overview"`
	VoteAverage float32 `json:"vote_average"`
}

func tmdbImplementation(w http.ResponseWriter, r *http.Request) {
	//For receiving API call
	client := http.Client{}
	movieRequest, httperr := http.NewRequest(http.MethodGet, APIURL, nil)
	if httperr != nil {
		log.Fatal(httperr)
	}
	movieResponse, geterr := client.Do(movieRequest)
	if geterr != nil {
		log.Fatal(geterr)
	}
	movieBody, readerr := ioutil.ReadAll(movieResponse.Body)
	if readerr != nil {
		log.Fatal(readerr)
	}
	randomMovie := movie{}
	jsonerr := json.Unmarshal(movieBody, &randomMovie)
	if jsonerr != nil {
		log.Fatal(jsonerr)
	}
	//For sending API call
	w.Header().Set("Access-Control-Allow-Origin", "*") //This heading is necessary for cross origin data transfer
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(randomMovie)

}

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	//Random number generator
	rand.Seed(time.Now().UnixNano())
	min := 10
	max := 634649
	RANDOM_NUMBER = rand.Intn(max-min+1) + min
	MOVIE_ID = strconv.Itoa(RANDOM_NUMBER)

	//APIKEY in .env file
	APIKEY = os.Getenv("API_KEY")

	//APIURL for API
	APIURL = "https://api.themoviedb.org/3/movie/" + MOVIE_ID + "?api_key=" + APIKEY

	router := mux.NewRouter()
	router.HandleFunc("/", tmdbImplementation)
	fmt.Println("Listening of port 8080")
	rand.Seed(time.Now().UnixNano())
	fmt.Println(RANDOM_NUMBER)
	http.ListenAndServe(":8080", router)
}
