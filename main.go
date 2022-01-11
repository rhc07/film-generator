package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
)

//APIURL for API
const APIURL = "https://api.themoviedb.org/3/movie/popular?api_key=" + config.APIKEY + "&page=1"

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

type movieList struct {
	List []movie `json:"results"`
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
	list := movieList{}
	jsonerr := json.Unmarshal(movieBody, &list)
	if jsonerr != nil {
		log.Fatal(jsonerr)
	}
	list.save()
}
