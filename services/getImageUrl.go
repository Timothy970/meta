package services

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"cloud.google.com/go/storage"
	"google.golang.org/api/option"
)

// Access token for the WhatsApp API.
var accessToken = "123456"

// GetMediaURL fetches the media URL and MIME type for a given media ID from the WhatsApp API.
func GetMediaURL(mediaID string) (string, string, error) {
	url := fmt.Sprintf("https://graph.facebook.com/v20.0/%s", mediaID)

	// Create a new HTTP GET request
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", "", err
	}

	// Set the Authorization header with the access token
	req.Header.Set("Authorization", "Bearer "+accessToken)

	client := &http.Client{}
	// Execute the request
	resp, err := client.Do(req)
	if err != nil {
		return "", "", err
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", "", err
	}

	// Check for a successful response
	if resp.StatusCode != http.StatusOK {
		return "", "", fmt.Errorf("error fetching media URL: %s", string(body))
	}

	// Parse the JSON response
	var result map[string]interface{}
	err = json.Unmarshal(body, &result)
	if err != nil {
		return "", "", err
	}

	// Extract the media URL and MIME type from the response
	mediaURL, urlOk := result["url"].(string)
	mimeType, mimeOk := result["mime_type"].(string)
	if !urlOk || !mimeOk {
		return "", "", fmt.Errorf("media URL or MIME type not found in response")
	}

	return mediaURL, mimeType, nil
}

// DownloadMediaData downloads the media data from a given URL.
func DownloadMediaData(mediaURL string) ([]byte, error) {
	// Execute a GET request to download the media data
	resp, err := http.Get(mediaURL)
	if err != nil {
		return nil, fmt.Errorf("error downloading media: %v", err)
	}
	defer resp.Body.Close()

	// Read the media data from the response body
	mediaData, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading media data: %v", err)
	}

	return mediaData, nil
}

// UploadMediaToGCS uploads the media data to Google Cloud Storage
func UploadMediaToGCS(mediaData []byte, mediaID, mimeType string) (string, error) {
	// Google Cloud Storage bucket name
	bucketName := "m_tickets"
	// Determine the file extension from the MIME type
	extension := strings.Split(mimeType, "/")[1]
	// Generate a unique object name using the media ID and a unique ID
	objectName := fmt.Sprintf("attachments/%s_%s.%s", mediaID, generateUniqueID(), extension)
	ctx := context.Background()

	// Create a new Google Cloud Storage client
	client, err := storage.NewClient(ctx, option.WithCredentialsFile("/path/to/service_account.json"))
	if err != nil {
		return "", fmt.Errorf("failed to create storage client: %v", err)
	}
	defer client.Close()

	// Get a handle to the bucket and the object
	bucket := client.Bucket(bucketName)
	object := bucket.Object(objectName)
	// Create a new writer for the object
	wc := object.NewWriter(ctx)
	// Write the media data to the object
	if _, err := io.Copy(wc, bytes.NewReader(mediaData)); err != nil {
		return "", fmt.Errorf("failed to upload media data: %v", err)
	}
	// Close the writer
	if err := wc.Close(); err != nil {
		return "", fmt.Errorf("failed to close writer: %v", err)
	}

	// Generate a signed URL for the object using the `storage.SignedURL` function
	url, err := storage.SignedURL(bucketName, objectName, &storage.SignedURLOptions{
		Scheme:  storage.SigningSchemeV4,
		Method:  "GET",
		Expires: time.Now().Add(365 * 24 * time.Hour), // Set the URL to expire in 1 year
	})
	if err != nil {
		return "", fmt.Errorf("failed to generate signed URL: %v", err)
	}

	log.Printf("File uploaded to GCS: %s", objectName)
	return url, nil
}

// generateUniqueID generates a unique ID based on the current time in nanoseconds.
func generateUniqueID() string {
	return fmt.Sprintf("%d", time.Now().UnixNano())
}
