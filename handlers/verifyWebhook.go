package handlers

import (
	"github.com/joho/godotenv" // For loading environment variables
	"log"
	"net/http"
	"os"
)

var myToken string

// Verification handler
func verifyWebhook(w http.ResponseWriter, r *http.Request) {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
	}
	// token = os.Getenv("TOKEN")
	myToken = os.Getenv("MYTOKEN")
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
