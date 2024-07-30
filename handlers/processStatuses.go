package handlers

import (
	"database/sql"
	"log"
	"whatsapp-webhook/services"
)

func processStatuses(db *sql.DB, status map[string]interface{}) {

	log.Printf(":::::::proces statuses::::::::")

	statusType := status["status"].(string)

	switch statusType {
	case "sent":
		log.Printf("Message %s status: %s", status["recipient_id"], statusType)
		timestamp := status["timestamp"].(string)
		id := status["id"].(string)
		// Check if the message ID already exists in the database
		var existingID string
		log.Printf(":::::::Select existingId::::::::")
		retrievedMessageID, _ := services.SelectFromDb(db, id)
		existingID = retrievedMessageID
		// If the record exists, update it
		if existingID != "" {

			services.UpdateToDb(db, "delivered", "delivered", timestamp, id)

		} else {
			// If the record does not exist, insert it

			services.UpdateToDb(db, "sent", "sent", timestamp, id)

		}
	case "read":
		log.Printf("Message read by %s: %s", status["recipient_id"], status["id"])
		timestamp := status["timestamp"].(string)
		id := status["id"].(string)
		services.UpdateToDb(db, "read", "read", timestamp, id)

	case "failed":
		// Extract the errors array and details
		var errorMessage string
		if errors, ok := status["errors"].([]interface{}); ok && len(errors) > 0 {
			firstError := errors[0].(map[string]interface{})
			errorMessage = firstError["message"].(string)
		}
		timestamp := status["timestamp"].(string)
		id := status["id"].(string)
		services.UpdateToDb(db, "failed", errorMessage, timestamp, id)

	default:
		log.Printf("Unknown status type: %s", statusType)
	}
}
