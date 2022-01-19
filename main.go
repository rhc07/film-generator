package main

import (
	"encoding/json"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	"github.com/joho/godotenv"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

var (
	APIKEY      string
	APIURL      string
	DBName      string
	DBUsername  string
	DBPassword  string
	DBHost      string
	DBPort      string
	DatabaseURL string
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

type movieList struct {
	List []movie `json:"results"`
}

func (list movieList) save() {
	db, dberr := gorm.Open("mysql", DatabaseURL)
	defer db.Close()
	if dberr != nil {
		log.Fatal(dberr)
	}
	db.Debug().DropTableIfExists(&movie{})
	db.AutoMigrate(&movie{})
	for _, row := range list.List {
		db.Debug().Create(&row)
	}
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
	//For sending API call
	w.Header().Set("Access-Control-Allow-Origin", "*") //This heading is necessary for cross origin data transfer
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(list)

}

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	//APIKEY in .env file
	APIKEY = os.Getenv("API_KEY")

	//APIURL for API
	APIURL = "https://api.themoviedb.org/3/movie/popular?api_key=" + APIKEY + "&page=2"

	//DBNAME variable for database name
	DBName = os.Getenv("DB_NAME")

	//DBUSERNAME variable for database username
	DBUsername = os.Getenv("DB_USERNAME")

	//DBPASSWORD variable for database password
	DBPassword = os.Getenv("DB_PASSWORD")

	//DBHOST variable for database address
	DBHost = os.Getenv("DB_HOST")

	//DBPORT variable for database default port number
	DBPort = os.Getenv("DB_PORT")

	//DATABASEURL variable
	DatabaseURL = DBUsername + ":" + DBPassword + "@tcp(" + DBHost + ":" + DBPort + ")/" + DBName + "?charset=utf8&parseTime=True&loc=Local"

	router := mux.NewRouter()
	router.HandleFunc("/", tmdbImplementation)
	fmt.Println("Listening of port 8080")
	http.ListenAndServe(":8080", router)
}
