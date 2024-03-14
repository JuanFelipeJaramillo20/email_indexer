package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
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

	//rootDir := "tests"
	rootDir := "enron_mail_20110402"

	numWorkers := 10

	fileChannel := make(chan string)
	emailChannel := make(chan *Email)
	done := make(chan struct{})

	var wg sync.WaitGroup
	wg.Add(numWorkers)

	for i := 0; i < numWorkers; i++ {
		go worker(fileChannel, emailChannel, &wg)
	}

	// Traverse the directory recursively
	go func() {
		defer close(fileChannel)
		err := filepath.Walk(rootDir, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				fmt.Printf("Error accessing path %q: %v\n", path, err)
				return nil // Continue traversal
			}
			if info.IsDir() {
				return nil // Skip directories
			}
			fileChannel <- path
			return nil
		})
		if err != nil {
			fmt.Println("Error traversing directory:", err)
		}
	}()

	go func() {
		wg.Wait()
		close(emailChannel)
	}()

	go func() {
		var bulkData []*Email
		for email := range emailChannel {
			bulkData = append(bulkData, email)
		}
		sendBulkToDatabase(bulkData)
		close(done)
	}()

	// Wait for sending emails to finish
	<-done

	elapsed := time.Since(start)
	log.Printf("Indexing took %s", elapsed)
}

// Worker function to process email files
func worker(fileChannel <-chan string, emailChannel chan<- *Email, wg *sync.WaitGroup) {
	defer wg.Done()
	for filePath := range fileChannel {
		email, err := processEmailFile(filePath)
		if err != nil {
			fmt.Println("Error processing email file:", err)
			continue
		}
		emailChannel <- email
	}
}

// Process an individual email file and return its JSON data
func processEmailFile(filePath string) (*Email, error) {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	email := parseEmail(string(content))
	return email, nil
}

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
func sendBulkToDatabase(emails []*Email) {
	// Define the database endpoint
	endpoint := "http://localhost:4080/api/_bulkv2"

	username := "admin"
	password := "admin"

	// Create HTTP client
	client := &http.Client{}

	// Create a slice to hold JSON data for bulk sending
	var bulkData []json.RawMessage

	// Convert each email to JSON
	for _, email := range emails {
		jsonData, err := json.Marshal(email)
		if err != nil {
			fmt.Println("Error marshalling email data:", err)
			return
		}
		bulkData = append(bulkData, json.RawMessage(jsonData))
	}

	// Create request body
	body := map[string]interface{}{
		"index":   "emails",
		"records": bulkData,
	}
	jsonBody, err := json.Marshal(body)
	if err != nil {
		fmt.Println("Error marshalling bulk data:", err)
		return
	}

	// Create request
	req, err := http.NewRequest("POST", endpoint, bytes.NewBuffer(jsonBody))
	if err != nil {
		fmt.Println("Error creating request:", err)
		return
	}

	// Set basic authentication header
	req.SetBasicAuth(username, password)
	req.Header.Set("Content-Type", "application/json")

	// Send request
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending request:", err)
		return
	}
	defer resp.Body.Close()

	// Check response status
	if resp.StatusCode != http.StatusOK {
		fmt.Printf("Unexpected status code: %d\n", resp.StatusCode)
		return
	}

	fmt.Println("Bulk data sent successfully")
}
