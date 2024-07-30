package main

import (
	// "bytes"
	"database/sql"
	"fmt"
	"github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv" // For loading environment variables
	"log"
	"net/http"
	"os"
	"whatsapp-webhook/handlers"
)

var (
	myToken string
	port    string
	logFile *os.File
	db      *sql.DB
)

// Simple test handler
func helloHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("hello this is webhook setup"))
}

func main() {
	// Load environment variables from .env file

	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
	}
	port = os.Getenv("PORT")

	// Open or create the log file
	logFile, err = os.OpenFile("webhook.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatalf("Failed to open log file: %v", err)
	}
	//database connection
	cfg := mysql.Config{
		User:   os.Getenv("DB_USER"),
		Passwd: os.Getenv("DB_PASSWORD"),
		Net:    "tcp",
		Addr:   "127.0.0.1:3306",
		DBName: "mydb",
	}
	// Get a database handle.
	// var err error
	db, err = sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		log.Fatal(err)
	}

	pingErr := db.Ping()
	if pingErr != nil {
		log.Fatal(pingErr)
	}
	fmt.Println("Connected!")
	log.Println("Successfully connected to the database")

	//routers

	router := mux.NewRouter()

	// router.HandleFunc("/webhook", handlers.verifyWebhook).Methods("GET")
	router.HandleFunc("/", helloHandler).Methods("GET")
	router.HandleFunc("/webhook", handlers.HandleWebhookCallback(db)).Methods("POST")

	fmt.Printf("Server is listening on port %s\n", port)
	log.Printf("Server started on port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, router))

}
