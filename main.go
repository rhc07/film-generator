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

//APIKEY in .env file
var APIKEY = os.Getenv("API_KEY")

//APIURL for API
var APIURL = "https://api.themoviedb.org/3/movie/popular?api_key=" + APIKEY + "&page=1"

//DBNAME variable for database name
var DBNAME = os.Getenv("DB_NAME")

//DBUSERNAME variable for database username
var DBUSERNAME = os.Getenv("DB_USERNAME")

//DBPASSWORD variable for database password
var DBPASSWORD = os.Getenv("DB_PASSWORD")

//DBHOST variable for database address
var DBHOST = os.Getenv("DB_HOST")

//DBPORT variable for database default port number
var DBPORT = os.Getenv("DB_PORT")

//DATABASEURL variable
var DATABASEURL = DBUSERNAME + ":" + DBPASSWORD + "@tcp(" + DBHOST + ":" + DBPORT + ")/" + DBNAME + "?charset=utf8&parseTime=True&loc=Local"

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

//func loadEnvFile(key string) string {
//	err := godotenv.Load(".env")
//
//	if err != nil {
//		log.Fatalf("Error loading .env file")
//	}
//
//	return os.Getenv()
//}

func (list movieList) save() {
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatalf("Error loading .env file")
	}
	db, dberr := gorm.Open("mysql", DATABASEURL)
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
	//os.LookupEnv(DBNAME)
	router := mux.NewRouter()
	router.HandleFunc("/", tmdbImplementation)
	fmt.Print(DATABASEURL)
	//fmt.Print(DBNAME)
	//os.LookupEnv(DBNAME)
	fmt.Println("Listening of port 8080")
	http.ListenAndServe(":8080", router)
}
