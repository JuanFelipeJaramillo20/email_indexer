package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// Email struct to hold email information
type Email struct {
	MessageID string `json:"message_id"`
	Date      string `json:"date"`
	From      string `json:"from"`
	To        string `json:"to"`
	Subject   string `json:"subject"`
	Content   string `json:"content"`
}

func main() {
	start := time.Now()
	// Define the root directory containing the email folders
	rootDir := "tests"

	// Create a slice to hold JSON data for bulk sending
	var bulkData []json.RawMessage

	// Traverse the directory recursively
	err := filepath.Walk(rootDir, func(path string, info os.FileInfo, err error) error {
		// Check if the current item is a file
		if err != nil {
			// Handle the error gracefully
			fmt.Printf("Error accessing path %q: %v\n", path, err)
			return nil // Continue traversal
		}
		if info.IsDir() {
			// Skip directories
			return nil
		}
		// Process the file as an email
		jsonData, err := processEmailFile(path)
		if err != nil {
			fmt.Println("Error processing email file:", err)
			return nil // Continue traversal
		}
		bulkData = append(bulkData, jsonData)
		return nil
	})

	if err != nil {
		fmt.Println("Error traversing directory:", err)
	}

	// Send bulk JSON data to the database
	if err := sendBulkToDatabase(bulkData); err != nil {
		fmt.Println("Error sending bulk data to database:", err)
	}

	elapsed := time.Since(start)
	log.Printf("indexing took %s", elapsed)
}

// Process an individual email file and return its JSON data
func processEmailFile(filePath string) (json.RawMessage, error) {
	// Read the content of the email file
	content, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	// Parse the email content
	email := parseEmail(string(content))

	// Convert email to JSON
	jsonData, err := json.Marshal(email)
	if err != nil {
		return nil, err
	}

	return json.RawMessage(jsonData), nil
}

// Parse the content of an email file
func parseEmail(content string) *Email {
	// Split the content into lines
	lines := strings.Split(content, "\n")

	// Extract email headers and content
	var email Email
	for _, line := range lines {
		if strings.HasPrefix(line, "Message-ID:") {
			email.MessageID = strings.TrimSpace(strings.TrimPrefix(line, "Message-ID:"))
		} else if strings.HasPrefix(line, "Date:") {
			email.Date = strings.TrimSpace(strings.TrimPrefix(line, "Date:"))
		} else if strings.HasPrefix(line, "From:") {
			email.From = strings.TrimSpace(strings.TrimPrefix(line, "From:"))
		} else if strings.HasPrefix(line, "To:") {
			email.To = strings.TrimSpace(strings.TrimPrefix(line, "To:"))
		} else if strings.HasPrefix(line, "Subject:") {
			email.Subject = strings.TrimSpace(strings.TrimPrefix(line, "Subject:"))
		} else {
			// Assume the rest is email content
			email.Content += line + "\n"
		}
	}

	return &email
}

// Send bulk JSON data to the database endpoint
func sendBulkToDatabase(data []json.RawMessage) error {
	// Define the database endpoint
	endpoint := "http://localhost:4080/api/_bulkv2"

	username := "admin"
	password := "admin"

	// Create HTTP client
	client := &http.Client{}

	// Create request body
	body := map[string]interface{}{
		"index":   "emails",
		"records": data,
	}
	jsonBody, err := json.Marshal(body)
	if err != nil {
		return err
	}

	// Create request
	req, err := http.NewRequest("POST", endpoint, bytes.NewBuffer(jsonBody))
	if err != nil {
		return err
	}

	// Set basic authentication header
	req.SetBasicAuth(username, password)
	req.Header.Set("Content-Type", "application/json")

	// Send request
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Check response status
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	return nil
}
