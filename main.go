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
	Id           int            `json:"id"`
	JourneyType  string         `json:"journey_type"`
	Departing    string         `json:"departing"`
	Destination  sql.NullString `json:"destination"`
	TicketName   sql.NullString `json:"ticket_name"`
	Date         time.Time      `json:"date"`
	RailcardUsed bool           `json:"railcard_used"`
	Cost         float64        `json:"cost"`
	TotalCost    float64        `json:"total_cost"`
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
		fmt.Printf("Error  Connecting to Database")
		fmt.Printf("%s", err)
	}
	return db
}

func closeDBConnection(db *sql.DB) {
	err := db.Close()
	if err != nil {
		fmt.Printf("Error closing connection to Database")
		fmt.Printf("%s", err)
	}
}

func getRailJourneys() (railJourneys []RailJourney) {
	db := getDBConnection()
	defer closeDBConnection(db)

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
			log.Printf("Error Processing Database results")
			fmt.Printf("%s", err)
			return
		}
		// and then append the new railJourney to the railJourney list
		railJourneys = append(railJourneys, railJourney)
	}
	return
}

func saveRailJourney(railJourney RailJourney) {
	db := getDBConnection()
	defer closeDBConnection(db)

	// Inserting records into database https://stackoverflow.com/a/16058741
	stmt, err := db.Prepare("INSERT rail_journeys SET journey_type=?, departing=?, destination=?, ticket_name=?, date=?, railcard_used=?, cost=?, total_cost=?")
	if err != nil {
		log.Printf("Error creating INSERT statment")
		fmt.Printf("%s", err)
		return
	}

	_, err = stmt.Exec(railJourney.JourneyType, railJourney.Departing, railJourney.Destination, railJourney.TicketName, railJourney.Date, railJourney.RailcardUsed, railJourney.Cost, railJourney.TotalCost)
	if err != nil {
		log.Printf("Error inserting into Database")
		fmt.Printf("%s", err)
		return
	}

	log.Printf("Record inserted into Database")
}

func RailJourneysHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("RailJourneysHandler called, Method %v", r.Method)

	method := r.Method
	if method == "POST" {
		// Decode JSON https://stackoverflow.com/a/15685432

		decoder := json.NewDecoder(r.Body)
		var railJourney RailJourney
		err := decoder.Decode(&railJourney)

		if err != nil {
			log.Printf("Failed to decode JSON, %v", err)
		}

		saveRailJourney(railJourney)
	} else {
		railJourneys := getRailJourneys()

		log.Printf("Database retured %v railJourney record(s)", len(railJourneys))

		err := json.NewEncoder(w).Encode(railJourneys)
		if err != nil {
			log.Printf("Failed to encode JSON, %v", err)
		}
	}
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
