package services

import (
	"database/sql"
	"fmt"
	"log"
)

// InsertToDb is a function for inserting data into the database
func InsertToDb(db *sql.DB, sender, recipient, messageID, msgType, message, timestamp, conversationID string) error {
	// Insert into conversations table
	_, err := db.Exec(
		`INSERT INTO conversations (sender, recipient, message_id)
		VALUES (?, ?, ?)`,
		sender,
		recipient,
		messageID,
	)
	if err != nil {
		return fmt.Errorf("error inserting into conversations: %v", err)
	}

	// Insert into messages table
	_, err = db.Exec(
		`INSERT INTO messages (message_id, type, message, timestamp)
		VALUES (?, ?, ?, ?)`,
		messageID,
		msgType,
		message,
		timestamp,
	)
	if err != nil {
		return fmt.Errorf("error inserting into messages: %v", err)
	}

	// Insert into correlator table
	_, err = db.Exec(
		`INSERT INTO correlator (message_id, conversation_id)
		VALUES (?, ?)`,
		messageID,
		conversationID,
	)
	if err != nil {
		return fmt.Errorf("error inserting into correlator: %v", err)
	}

	return nil
}

// UpdateToDb is a function for updating data in the database
func UpdateToDb(db *sql.DB, status, statusDescription, timestamp, messageID string) error {
	log.Printf(":::::::update to db::::::::")

	_, err := db.Exec(
		`UPDATE messages
		SET status = ?, status_description = ?, timestamp = ?
		WHERE message_id = ?`,
		status,
		statusDescription,
		timestamp,
		messageID,
	)
	if err != nil {
		return fmt.Errorf("error updating messages: %v", err)
	}

	return nil
}
func SelectFromDb(db *sql.DB, messageID string) (string, error) {
	log.Printf(":::::::SelectFromDb::::::::")

	var retrievedMessageID string

	err := db.QueryRow(
		`SELECT message_id FROM messages
		WHERE message_id = ?`,
		messageID,
	).Scan(&retrievedMessageID)
	if err != nil {

		return "", fmt.Errorf("error querying message_id %s: %v", messageID, err)
	}

	return retrievedMessageID, nil
}
