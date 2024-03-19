package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"runtime/pprof"
	"strings"
	"sync"
	"time"
)

type Email struct {
	MessageID string `json:"message_id"`
	Date      string `json:"date"`
	From      string `json:"from"`
	To        string `json:"to"`
	Subject   string `json:"subject"`
	Content   string `json:"content"`
}

type Worker struct {
	id           int
	filePathChan <-chan string
	resultChan   chan<- *Email
}

func NewWorker(id int, filePathChan <-chan string, resultChan chan<- *Email) *Worker {
	return &Worker{
		id:           id,
		filePathChan: filePathChan,
		resultChan:   resultChan,
	}
}

func (w *Worker) Start() {
	for filePath := range w.filePathChan {
		email, err := processEmailFile(filePath)
		if err != nil {
			fmt.Printf("Error processing email file %s: %v\n", filePath, err)
			continue
		}
		w.resultChan <- email
	}
}

func main() {
	cpuFile, err := os.Create("indexer.pprof")
	if err != nil {
		log.Fatalf("Error creating CPU profile file: %v\n", err)
	}
	defer cpuFile.Close()

	if err := pprof.StartCPUProfile(cpuFile); err != nil {
		log.Fatalf("Error starting CPU profile: %v\n", err)
	}
	defer pprof.StopCPUProfile()

	start := time.Now()

	const batchSize = 1000
	numWorkers := runtime.NumCPU()

	filePathChan := make(chan string, numWorkers)
	results := make(chan *Email, batchSize)
	done := make(chan struct{})

	// Start worker goroutines
	var wg sync.WaitGroup
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		worker := NewWorker(i, filePathChan, results)
		go func(w *Worker) {
			defer wg.Done()
			w.Start()
		}(worker)
	}

	// Traverse the directory and send file paths to the filePathChan
	go func() {
		defer close(filePathChan)
		err := filepath.Walk("enron_mail_20110402", func(path string, info os.FileInfo, err error) error {
			if err != nil {
				fmt.Printf("Error accessing path %q: %v\n", path, err)
				return nil
			}
			if !info.IsDir() {
				filePathChan <- path
			}
			return nil
		})
		if err != nil {
			fmt.Println("Error traversing directory:", err)
		}
	}()

	// Process results in batches and send to database
	go func() {
		defer close(results)
		var batch []*Email
		var batchNumber = 1
		for email := range results {
			batch = append(batch, email)
			if len(batch) == batchSize {
				sendBulkToDatabase(batch, batchNumber)
				batch = nil // Reset batch
				batchNumber++
			}
		}

		if len(batch) > 0 {
			sendBulkToDatabase(batch, batchNumber)
			batchNumber++
		}
		close(done)
	}()

	wg.Wait()

	<-done

	elapsed := time.Since(start)
	log.Printf("Indexing took %s", elapsed)
}

// Process an individual email file and return its parsed data
func processEmailFile(filePath string) (*Email, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var buf bytes.Buffer
	_, err = buf.ReadFrom(file)
	if err != nil {
		return nil, err
	}

	content := buf.String()
	email := parseEmail(content)
	return email, nil
}

// Parse email content and extract relevant information
func parseEmail(content string) *Email {
	lines := strings.Split(content, "\n")
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
func sendBulkToDatabase(emails []*Email, batchNumber int) {
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
		"index":   "emails_new_pattern_03",
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

	fmt.Println("Bulk data sent successfully:", batchNumber*100, " of 517.425 sent")
}
