package main

import (
	// "bytes"
	"database/sql"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv" // For loading environment variables
	"log"
	"net/http"
	"os"
)

var (
	myToken string
	port    string
	logFile *os.File
	db      *sql.DB
)

// Load environment variables from .env file
func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
	}
	// token = os.Getenv("TOKEN")
	myToken = os.Getenv("MYTOKEN")
	port = os.Getenv("PORT")

	// Open or create the log file
	logFile, err = os.OpenFile("webhook.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatalf("Failed to open log file: %v", err)
	}
	// Ensure the log file is closed when the program exits
	log.SetOutput(logFile)

	// Get the database connection string from environment variables
	connStr := fmt.Sprintf("user=%s dbname=%s sslmode=%s password=%s host=%s port=%s",
		os.Getenv("DB_USER"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_SSLMODE"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"))

	// Initialize the database connection
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("Error connecting to the database: %v", err)
	}

	// Verify the database connection
	err = db.Ping()
	if err != nil {
		log.Fatalf("Error pinging the database: %v", err)
	}

	log.Println("Successfully connected to the database")
}

// Verification handler
func verifyWebhook(w http.ResponseWriter, r *http.Request) {
	mode := r.URL.Query().Get("hub.mode")
	challenge := r.URL.Query().Get("hub.challenge")
	verifyToken := r.URL.Query().Get("hub.verify_token")

	// Log detailed information about the incoming request
	log.Printf("Received verification request: Method=%s, URL=%s, Parameters: mode=%s, challenge=%s, verify_token=%s", r.Method, r.URL.String(), mode, challenge, verifyToken)

	if mode != "" && verifyToken != "" {
		if mode == "subscribe" && verifyToken == myToken {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(challenge))
			log.Printf("Webhook verified with challenge: %s", challenge)
		} else {
			w.WriteHeader(http.StatusForbidden)
			log.Printf("Failed verification: invalid mode or token")
		}
	} else {
		w.WriteHeader(http.StatusBadRequest)
		log.Printf("Failed verification: missing mode or token")
	}
}

// Simple test handler
func helloHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("hello this is webhook setup"))
}

func main() {
	router := mux.NewRouter()

	router.HandleFunc("/webhook", verifyWebhook).Methods("GET")
	router.HandleFunc("/webhook", handleWebhookCallback).Methods("POST")
	router.HandleFunc("/", helloHandler).Methods("GET")

	fmt.Printf("Server is listening on port %s\n", port)
	log.Printf("Server started on port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, router))
}
