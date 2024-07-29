package main

import (
	"database/sql"
	"encoding/json"
	"log"
)

func processMessages(message map[string]interface{}) {
	messageType := message["type"].(string)
	switch messageType {
	case "unsupported":
		log.Printf("Message deleted by %s", message["from"])
	case "text":
		log.Printf("Text message from %s: %s", message["from"], message["text"].(map[string]interface{})["body"])

		var referralJSON string
		if referral, ok := message["referral"].([]interface{}); ok && len(referral) > 0 {
			jsonData, err := json.Marshal(referral)
			if err != nil {
				log.Printf("Error marshaling referral data to JSON: %v", err)
			} else {
				referralJSON = string(jsonData)
			}
		}

		var contextJSON string
		if context, ok := message["context"].([]interface{}); ok && len(context) > 0 {
			jsonData, err := json.Marshal(context)
			if err != nil {
				log.Printf("Error marshaling context data to JSON: %v", err)
			} else {
				contextJSON = string(jsonData)
			}
		}

		_, err := db.Exec(`
		INSERT INTO conversations (message_id, from_no, type, timestamp, message, referrals, inquiries)
		VALUES ($1, $2, $3, $4, $5, $6, $7)`,
			message["id"],
			message["from"],
			messageType,
			message["timestamp"],
			message["text"].(map[string]interface{})["body"],
			referralJSON,
			contextJSON,
		)
		if err != nil {
			log.Printf("Error inserting type text message into database: %v", err)
		}

	case "reaction":
		log.Printf("Reaction from %s: %s", message["from"], message["reaction"].(map[string]interface{})["emoji"])

		_, err := db.Exec(`
		INSERT INTO conversations (message_id, frome_no, type, timestamp, emoji, emoji_id)
		VALUES ($1, $2, $3, $4, $5, $6)`,
			message["id"],
			message["from"],
			messageType,
			message["timestamp"],
			message["reaction"].(map[string]interface{})["emoji"],
			message["reaction"].(map[string]interface{})["message_id"],
		)
		if err != nil {
			log.Printf("Error inserting type reaction message into database: %v", err)

		}
	case "image":
		log.Printf("Image message from %s: %s (Caption: %s)", message["from"], message["image"].(map[string]interface{})["id"], message["image"].(map[string]interface{})["caption"])
		_, err := db.Exec(`
		INSERT INTO conversations (message_id, from_no, type, timestamp, caption, mime_type, image_hash, image_id)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`,
			message["id"],
			message["from"],
			messageType,
			message["timestamp"],
			message["image"].(map[string]interface{})["caption"],
			message["image"].(map[string]interface{})["mime_type"],
			message["image"].(map[string]interface{})["sha256"],
			message["image"].(map[string]interface{})["id"],
		)
		if err != nil {
			log.Printf("Error inserting type media message into database: %v", err)

		}
	case "sticker":
		log.Printf("Sticker message from %s: %s", message["from"], message["sticker"].(map[string]interface{})["id"])
		_, err := db.Exec(`
		INSERT INTO conversations (message_id, from_id, type, timestamp, mime_type, image_hash, image_id)
		VALUES ($1, $2, $3, $4, $5, $6, $7)`,
			message["id"],
			message["from"],
			messageType,
			message["timestamp"],
			message["sticker"].(map[string]interface{})["mime_type"],
			message["sticker"].(map[string]interface{})["sha256"],
			message["sticker"].(map[string]interface{})["id"],
		)
		if err != nil {
			log.Printf("Error inserting type sticker message into database: %v", err)

		}
	case "unknown":
		log.Printf("Unkown message from %s: %s", message["from"], message["sticker"].(map[string]interface{})["id"])
		_, err := db.Exec(`
		INSERT INTO conversations (message_id, from_no, type, timestamp, message)
		VALUES ($1, $2, $3, $4, $5)`,
			message["id"],
			message["from"],
			messageType,
			message["timestamp"],
			message["errors"].(map[string]interface{})["details"],
		)
		if err != nil {
			log.Printf("Error inserting type sticker message into database: %v", err)

		}
	case "button":
		log.Printf("button message from %s: %s", message["from"], message["sticker"].(map[string]interface{})["id"])
		_, err := db.Exec(`
	INSERT INTO conversations (message_id, from_no, type, timestamp, message)
	VALUES ($1, $2, $3, $4, $5)`,
			message["id"],
			message["from"],
			messageType,
			message["timestamp"],
			message["button"].(map[string]interface{})["text"],
		)
		if err != nil {
			log.Printf("Error inserting type button message into database: %v", err)

		}
	case "interactive":
		log.Printf("List reply message from %s: %s", message["from"], message["sticker"].(map[string]interface{})["id"])

		var title string
		interactive := message["interactive"].(map[string]interface{})

		if listReply, ok := interactive["list_reply"].(map[string]interface{}); ok {
			title = listReply["title"].(string)
		} else if buttonReply, ok := interactive["button_reply"].(map[string]interface{}); ok {
			title = buttonReply["title"].(string)
		}

		_, err := db.Exec(`
		INSERT INTO conversations (message_id, from_no, type, timestamp, message)
		VALUES ($1, $2, $3, $4, $5)`,
			message["id"],
			message["from"],
			messageType,
			message["timestamp"],
			title,
		)
		if err != nil {
			log.Printf("Error inserting type button message into database: %v", err)
		}
	case "system":
		system, _ := json.Marshal(message["system"])

		_, err := db.Exec(`
		INSERT INTO conversations (message_id, from_no, type, timestamp, message)
		VALUES ($1, $2, $3, $4, $5)`,
			message["id"],
			message["from"],
			messageType,
			string(system),
		)
		if err != nil {
			log.Printf("Error inserting type button message into database: %v", err)
		}

	default:
		log.Printf("Unknown message type from %s: %s", message["from"], messageType)
	}
	if location, ok := message["location"].([]interface{}); ok && len(location) > 0 {
		if locationMap, ok := location[0].(map[string]interface{}); ok {
			log.Printf("Locations message from %s: %s", message["from"], message["sticker"].(map[string]interface{})["id"])
			_, err := db.Exec(`
			INSERT INTO conversations (message_id, from, type, timestamp, latitude, longitude, name, address)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`,
				message["id"],
				message["from"],
				"location",
				message["timestamp"],
				locationMap["latitude"],
				locationMap["longitude"],
				locationMap["name"],
				locationMap["address"],
			)
			if err != nil {
				log.Printf("Error inserting type locations message into database: %v", err)
			}
		} else {
			log.Printf("Invalid location format")
		}
	} else {
		log.Printf("Location data is missing or invalid")
	}

	if contacts, ok := message["contacts"].([]interface{}); ok && len(contacts) > 0 {
		log.Printf("Contacts message from %s: %s", message["from"], contacts)
		contact := contacts[0].(map[string]interface{})
		// Extract and convert fields to JSON strings
		addresses, _ := json.Marshal(contact["addresses"])
		birthday := contact["birthday"].(string)
		emails, _ := json.Marshal(contact["emails"])
		name, _ := json.Marshal(contact["name"])
		org, _ := json.Marshal(contact["org"])
		phones, _ := json.Marshal(contact["phones"])
		urls, _ := json.Marshal(contact["urls"])
		_, err := db.Exec(`
		INSERT INTO conversations (message_id, from_no, type, timestamp, addresses, birthday, emails, name, org, phones, urls)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)`,
			message["id"],
			message["from"],
			"contacts",
			message["timestamp"],
			string(addresses),
			birthday,
			string(emails),
			string(name),
			string(org),
			string(phones),
			string(urls),
		)
		if err != nil {
			log.Printf("Error inserting type contacts message into database: %v", err)

		}
	}

}
func processStatuses(status map[string]interface{}) {
	statusType := status["status"].(string)
	conversation := status["conversation"].(map[string]interface{})
	origin := conversation["origin"].(map[string]interface{})

	switch statusType {
	case "sent":
		log.Printf("Message %s status: %s", status["recipient_id"], statusType)

		// Check if the message ID already exists in the database
		var existingID string
		err := db.QueryRow("SELECT message_id FROM messages WHERE message_id = $1", status["id"]).Scan(&existingID)
		if err != nil && err != sql.ErrNoRows {
			log.Printf("Error querying database: %v", err)
			return
		}

		// If the record exists, update it
		if existingID != "" {
			_, err := db.Exec(`
				UPDATE messages
				SET status = $1, timestamp = $2
				WHERE message_id = $3`,
				"delivered",
				status["timestamp"],
				status["id"],
			)
			if err != nil {
				log.Printf("Error updating status in database: %v", err)
			}
		} else {
			// If the record does not exist, insert it
			_, err := db.Exec(`
				INSERT INTO messages (message_id, status, timestamp, recipient_id, conversation_id, conversation_type)
				VALUES ($1, $2, $3, $4, $5, $6)`,
				status["id"],
				statusType,
				status["timestamp"],
				status["recipient_id"],
				conversation["id"],
				origin["type"],
			)
			if err != nil {
				log.Printf("Error inserting status into database: %v", err)
			}
		}
	case "read":
		log.Printf("Message read by %s: %s", status["recipient_id"], status["id"])

		// Update the status to failed in the database
		_, err := db.Exec(`
			UPDATE messages
			SET status = $1, timestamp = $2,
			WHERE message_id = $4`,
			statusType,
			status["timestamp"],
			status["id"],
		)
		if err != nil {
			log.Printf("Error updating status to failed in database: %v", err)
		}
	case "failed":
		// Extract the errors array and details
		var errorMessage string
		if errors, ok := status["errors"].([]interface{}); ok && len(errors) > 0 {
			firstError := errors[0].(map[string]interface{})
			errorMessage = firstError["message"].(string)
		}

		log.Printf("Message failed to %s: %s, error: %s", status["recipient_id"], status["id"], errorMessage)

		// Update the status to failed in the database
		_, err := db.Exec(`
			UPDATE messages
			SET status = $1, timestamp = $2, error = $3
			WHERE message_id = $4`,
			statusType,
			status["timestamp"],
			errorMessage,
			status["id"],
		)
		if err != nil {
			log.Printf("Error updating status to failed in database: %v", err)
		}

	default:
		log.Printf("Unknown status type: %s", statusType)
	}
}
