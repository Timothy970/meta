package handlers

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
)

// Callback handler
func HandleWebhookCallback(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("::::::::::handleWebhookCallback::::::::::::::")

		var payload map[string]interface{}

		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			log.Printf("Could not decode callback: %v", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		// Log the received callback data
		log.Printf("Received callback: %+v", payload)

		// Process the payload if it is from a WhatsApp business account
		if payload["object"] == "whatsapp_business_account" {
			if entries, ok := payload["entry"].([]interface{}); ok {
				for _, entry := range entries {
					entryMap := entry.(map[string]interface{})

					// Loop through each change in the entry
					if changes, ok := entryMap["changes"].([]interface{}); ok {
						for _, change := range changes {
							changeMap := change.(map[string]interface{})

							// Extract the value field which contains the main data
							if value, ok := changeMap["value"].(map[string]interface{}); ok {

								// Check and process messages
								if messages, ok := value["messages"].([]interface{}); ok {
									for _, message := range messages {
										messageMap := message.(map[string]interface{})
										processMessages(db, messageMap)

									}
								}

								// Check and process statuses
								if statuses, ok := value["statuses"].([]interface{}); ok {
									for _, status := range statuses {
										statusMap := status.(map[string]interface{})

										processStatuses(db, statusMap)

									}
								}
							}
						}
					}
				}
			}
		}

	}
}
