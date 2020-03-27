package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"
	"time"
)

// Help from
// https://tutorialedge.net/golang/golang-mysql-tutorial/
// https://tutorialedge.net/golang/creating-restful-api-with-golang/
// https://flaviocopes.com/golang-sql-database/
// https://www.thepolyglotdeveloper.com/2017/04/using-sqlite-database-golang-application/

type RailJourney struct {
	Id           int       `json:"id"`
	JourneyType  string    `json:"journey_type"`
	Departing    string    `json:"departing"`
	Destination  string    `json:"destination"`
	TicketName   string    `json:"ticket_name"`
	Date         time.Time `json:"date"`
	RailcardUsed bool      `json:"railcard_used"`
	Cost         float64   `json:"cost"`
	TotalCost    float64   `json:"total_cost"`
}

func getDBConnection() *sql.DB {
	DbUsername := os.Getenv("DB_USERNAME")
	DbPassword := os.Getenv("DB_PASSWORD")
	DbName := os.Getenv("DB_NAME")
	DbUrl := os.Getenv("DB_URL")
	DbPort := os.Getenv("DB_PORT")

	// ?parseTime=true <- add this to handle mysql DATE objects into time.TIME go objects
	dbConnectionString := fmt.Sprintf("%v:%v@tcp(%v:%v)/%v?parseTime=true", DbUsername, DbPassword, DbUrl, DbPort, DbName)

	db, err := sql.Open("mysql", dbConnectionString)

	if err != nil {
		fmt.Printf("Error  Connecting to DB")
		fmt.Printf("%s", err)
	}
	return db
}

func getRailJourneys() (railJourneys []RailJourney) {
	db := getDBConnection()
	defer db.Close()

	results, err := db.Query("SELECT * FROM rail_journeys")

	if err != nil {
		fmt.Printf("Error  !!!")
		fmt.Printf("%s", err)
	}

	for results.Next() {
		var railJourney RailJourney
		// for each row, scan the result into the railJourney composite object
		err = results.Scan(
			&railJourney.Id,
			&railJourney.JourneyType,
			&railJourney.Departing,
			&railJourney.Destination,
			&railJourney.TicketName,
			&railJourney.Date,
			&railJourney.RailcardUsed,
			&railJourney.Cost,
			&railJourney.TotalCost,
		)

		if err != nil {
			log.Printf("Error Processing DB results")
			fmt.Printf("%s", err)
			return
		}
		// and then append the new railJourney to the railJourney list
		railJourneys = append(railJourneys, railJourney)
	}
	return
}

func RailJourneysHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("RailJourneysHandler called for %v", r.RequestURI)

	railJourneys := getRailJourneys()

	log.Printf("Database retured %v railJourney record(s)", len(railJourneys))

	json.NewEncoder(w).Encode(railJourneys)
}

func main() {
	// Setup env file access
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	handleRequests()
}

func handleRequests() {
	// creates a new instance of a mux router
	myRouter := mux.NewRouter().StrictSlash(true)

	myRouter.HandleFunc("/rail-journeys", RailJourneysHandler)

	Port := fmt.Sprintf(":%s", os.Getenv("PORT"))
	log.Printf("Handler listening on Port %v", Port)
	log.Fatal(http.ListenAndServe(Port, myRouter))
}
