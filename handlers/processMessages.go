package handlers

import (
	"database/sql"
	"encoding/json"
	"log"
	"whatsapp-webhook/services"
)

func processMessages(db *sql.DB, message map[string]interface{}) {
	log.Printf(":::::::proces messages::::::::")
	log.Printf("message recieved::::::: %+v", message)

	// Check if the message has a value called "contacts"
	if contactsSlice, ok := message["contacts"].([]interface{}); ok {

		contacts := contactsSlice[0].(map[string]interface{})

		// Log the contacts' name
		name := contacts["name"].(map[string]interface{})

		formattedName, _ := json.Marshal(name["formatted_name"])
		log.Printf("Contacts message from %s: %s", message["from"], string(formattedName))

		// Extract and convert fields to JSON strings
		messageContacts, _ := json.Marshal(message["contacts"])

		services.InsertToDb(db, message["from"].(string), message["from"].(string), message["id"].(string), "contacts", string(messageContacts), message["timestamp"].(string), "null")

	}

	// Check if the message has a value called "location"
	if location, ok := message["location"].(map[string]interface{}); ok {
		log.Printf(":::::::IN LOCATION TYPE:::: ")
		messageLocation, _ := json.Marshal(location)

		services.InsertToDb(db, message["from"].(string), message["from"].(string), message["id"].(string), "location", string(messageLocation), message["timestamp"].(string), "null")

	}

	messageType := message["type"].(string)
	switch messageType {
	case "unsupported":
		log.Printf("Message deleted by %s", message["from"])
	case "text":
		log.Printf("Text message from %s: %s", message["from"], message["text"].(map[string]interface{})["body"])
		services.InsertToDb(db, message["from"].(string), message["from"].(string), message["id"].(string), messageType, message["text"].(map[string]interface{})["body"].(string), message["timestamp"].(string), "null")

	case "reaction":
		log.Printf("Reaction from %s: %s", message["from"], message["reaction"].(map[string]interface{})["emoji"])
		messageReaction, _ := json.Marshal(message["reaction"])

		services.InsertToDb(db, message["from"].(string), message["from"].(string), message["id"].(string), messageType, string(messageReaction), message["timestamp"].(string), "null")

	case "image":
		log.Printf("Image message from %s: %s (Caption: %s)", message["from"], message["image"].(map[string]interface{})["id"], message["image"].(map[string]interface{})["caption"])
		image := message["image"].(map[string]interface{})
		imageID := image["id"].(string)
		// Fetch the media URL using the image ID
		imageURL, mimeType, err := services.GetMediaURL(imageID)
		if err != nil {
			log.Printf("Error fetching media URL: %v", err)
			return
		}

		// Download the media data
		mediaData, err := services.DownloadMediaData(imageURL)
		if err != nil {
			log.Printf("Error downloading media data: %v", err)
			return
		}

		// Upload the media data to GCS and get the signed URL
		signedURL, err := services.UploadMediaToGCS(mediaData, imageID, mimeType)
		if err != nil {
			log.Printf("Error uploading media to GCS: %v", err)
			return
		}

		services.InsertToDb(db, message["from"].(string), message["from"].(string), message["id"].(string), messageType, signedURL, message["timestamp"].(string), "null")

	case "sticker":
		log.Printf("Sticker message from %s: %s", message["from"], message["sticker"].(map[string]interface{})["id"])
		sticker := message["sticker"].(map[string]interface{})
		imageID := sticker["id"].(string)
		// Fetch the media URL using the image ID
		imageURL, mimeType, err := services.GetMediaURL(imageID)
		if err != nil {
			log.Printf("Error fetching media URL: %v", err)
			return
		}

		// Download the media data
		mediaData, err := services.DownloadMediaData(imageURL)
		if err != nil {
			log.Printf("Error downloading media data: %v", err)
			return
		}

		// Upload the media data to GCS and get the signed URL
		signedURL, err := services.UploadMediaToGCS(mediaData, imageID, mimeType)
		if err != nil {
			log.Printf("Error uploading media to GCS: %v", err)
			return
		}
		services.InsertToDb(db, message["from"].(string), message["from"].(string), message["id"].(string), messageType, signedURL, message["timestamp"].(string), "null")

	case "unknown":
		log.Printf("Unkown message from %s: %s", message["from"],
			message["errors"].([]interface{})[0].(map[string]interface{})["title"].(string),
		)
		messageUnkown, _ := json.Marshal(message["errors"])

		services.InsertToDb(db, message["from"].(string), message["from"].(string), message["id"].(string), messageType, string(messageUnkown), message["timestamp"].(string), "null")

	case "button":
		log.Printf("button message from %s: %s", message["from"], message["sticker"].(map[string]interface{})["id"])
		messageButton, _ := json.Marshal(message["button"])

		services.InsertToDb(db, message["from"].(string), message["from"].(string), message["id"].(string), messageType, string(messageButton), message["timestamp"].(string), "null")

	case "interactive":
		log.Printf("List reply message from %s: %s", message["from"], message["sticker"].(map[string]interface{})["id"])

		var typeMessage string
		var messageInteractive string

		interactive := message["interactive"].(map[string]interface{})

		if listReply, ok := interactive["list_reply"].(map[string]interface{}); ok {
			typeMessage = interactive["type"].(string)
			messageInteractive1, _ := json.Marshal(listReply)
			messageInteractive = string(messageInteractive1)

		} else if buttonReply, ok := interactive["button_reply"].(map[string]interface{}); ok {
			typeMessage = buttonReply["type"].(string)
			messageInteractive1, _ := json.Marshal(listReply)
			messageInteractive = string(messageInteractive1)

		}
		services.InsertToDb(db, message["from"].(string), message["from"].(string), message["id"].(string), typeMessage, messageInteractive, message["timestamp"].(string), "null")

	default:
		log.Printf("Unknown message type from %s: %s", message["from"], messageType)
	}

}
