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
	XFrom     string `json:"x_from"`
	To        string `json:"to"`
	XTo       string `json:"x_to"`
	Subject   string `json:"subject"`
	Cc        string `json:"cc"`
	XCc       string `json:"x_cc"`
	Bcc       string `json:"bcc"`
	XBcc      string `json:"x_bcc"`
	Content   string `json:"content"`
}

type Worker struct {
	id              int
	filePathChannel <-chan string
	resultsChannel  chan<- *Email
}

func NewWorker(id int, filePathChannel <-chan string, resultsChannel chan<- *Email) *Worker {
	return &Worker{
		id:              id,
		filePathChannel: filePathChannel,
		resultsChannel:  resultsChannel,
	}
}

func (w *Worker) Start() {
	for filePath := range w.filePathChannel {
		email, err := processEmailFile(filePath)
		if err != nil {
			fmt.Printf("Error processing email file %s: %v\n", filePath, err)
			continue
		}
		w.resultsChannel <- email
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
	fmt.Printf("Numbers of workers created: %d", numWorkers)
	filePathChannel := make(chan string, numWorkers)
	results := make(chan *Email, batchSize)
	done := make(chan struct{})

	var wg sync.WaitGroup
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		worker := NewWorker(i, filePathChannel, results)
		go func(w *Worker) {
			defer wg.Done()
			w.Start()
		}(worker)
	}

	go func() {
		defer close(filePathChannel)
		err := filepath.Walk("enron_mail_20110402", func(path string, info os.FileInfo, err error) error {
			if err != nil {
				fmt.Printf("Error accessing path %q: %v\n", path, err)
				return nil
			}
			if !info.IsDir() {
				filePathChannel <- path
			}
			return nil
		})
		if err != nil {
			fmt.Println("Error traversing directory:", err)
		}
	}()

	go func() {
		var batch []*Email
		var batchNumber = 1
		for email := range results {
			batch = append(batch, email)
			if len(batch) == batchSize {
				sendBulkToDatabase(batch, batchNumber, len(batch))
				batch = nil
				batchNumber++
			}
		}

		if len(batch) > 0 {
			sendBulkToDatabase(batch, batchNumber, len(batch))
			batchNumber++
		}
		close(done)
	}()

	wg.Wait()
	close(results)
	<-done

	elapsed := time.Since(start)
	log.Printf("Indexing took %s", elapsed)
}

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

func parseEmail(content string) *Email {
	lines := strings.Split(content, "\n")
	var email Email
	for _, line := range lines {
		switch {
		case strings.HasPrefix(line, "Message-ID:"):
			email.MessageID = strings.TrimSpace(strings.TrimPrefix(line, "Message-ID:"))
		case strings.HasPrefix(line, "Date:"):
			email.Date = strings.TrimSpace(strings.TrimPrefix(line, "Date:"))
		case strings.HasPrefix(line, "From:"):
			email.From = strings.TrimSpace(strings.TrimPrefix(line, "From:"))
		case strings.HasPrefix(line, "X-From:"):
			email.XFrom = strings.TrimSpace(strings.TrimPrefix(line, "X-From:"))
		case strings.HasPrefix(line, "To:"):
			email.To = strings.TrimSpace(strings.TrimPrefix(line, "To:"))
		case strings.HasPrefix(line, "X-To:"):
			email.XTo = strings.TrimSpace(strings.TrimPrefix(line, "X-To:"))
		case strings.HasPrefix(line, "Subject:"):
			email.Subject = strings.TrimSpace(strings.TrimPrefix(line, "Subject:"))
		case strings.HasPrefix(line, "Cc:"):
			email.Cc = strings.TrimSpace(strings.TrimPrefix(line, "Cc:"))
		case strings.HasPrefix(line, "X-Cc:"):
			email.XCc = strings.TrimSpace(strings.TrimPrefix(line, "X-Cc:"))
		case strings.HasPrefix(line, "Bcc:"):
			email.Bcc = strings.TrimSpace(strings.TrimPrefix(line, "Bcc:"))
		case strings.HasPrefix(line, "X-Bcc:"):
			email.XBcc = strings.TrimSpace(strings.TrimPrefix(line, "X-Bcc:"))
		default:
			email.Content += line + "\n"
		}
	}
	return &email
}

func sendBulkToDatabase(emails []*Email, batchNumber int, batchSize int) {

	endpoint := "http://3.18.107.19:4080/api/_bulkv2"
	username := "admin"
	password := "admin"

	client := &http.Client{}

	var bulkData []json.RawMessage

	for _, email := range emails {
		jsonData, err := json.Marshal(email)
		if err != nil {
			fmt.Println("Error marshalling email data:", err)
			return
		}
		bulkData = append(bulkData, json.RawMessage(jsonData))
	}

	body := map[string]interface{}{
		"index":   "emails",
		"records": bulkData,
	}
	jsonBody, err := json.Marshal(body)
	if err != nil {
		fmt.Println("Error marshalling bulk data:", err)
		return
	}

	req, err := http.NewRequest("POST", endpoint, bytes.NewBuffer(jsonBody))
	if err != nil {
		fmt.Println("Error creating request:", err)
		return
	}

	req.SetBasicAuth(username, password)
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending request:", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Printf("Unexpected status code: %d\n", resp.StatusCode)
		return
	}

	fmt.Println("Bulk data sent successfully:", batchNumber*batchSize)
}
