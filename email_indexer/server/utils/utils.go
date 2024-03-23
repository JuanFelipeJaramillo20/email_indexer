package utils

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

type SearchRequest struct {
	Term  string `json:"term"`
	Page  string `json:"page,omitempty"`
	Order string `json:"order,omitempty"`
}

func SearchHandler(w http.ResponseWriter, r *http.Request) {
	var searchReq SearchRequest
	err := json.NewDecoder(r.Body).Decode(&searchReq)
	if err != nil {
		http.Error(w, "Failed to parse request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()
	fmt.Printf("Entrando", searchReq.Term)

	if searchReq.Term == "" {
		http.Error(w, "Search term can't be empty", http.StatusBadRequest)
		fmt.Println("Search term can't empty")
		return
	}

	if searchReq.Page == "" {
		searchReq.Page = "0"
	}
	if len(searchReq.Order) > 2 {
		searchReq.Order = ""
	}

	query := fmt.Sprintf(
		`{
        	"search_type": "matchphrase",
        	"query": {
            	"term": "%s"
        	},
			"sort_fields": ["%sdate"],
        	"from": %s,
        	"max_results": 20,
        	"_source": []
    	}`,
		searchReq.Term, searchReq.Order, searchReq.Page)
	fmt.Print("QUERY: ", query)
	requestURL := SERVER_URL + "_search"
	zincRequest, err := http.NewRequest("POST", requestURL, strings.NewReader(query))
	if err != nil {
		handleError(w, err)
		return
	}

	zincRequest.SetBasicAuth(USERNAME, PASSWORD)
	zincRequest.Header.Set("Content-Type", "application/json")
	zincRequest.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_4) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/81.0.4044.138 Safari/537.36")

	zincResponse, err := http.DefaultClient.Do(zincRequest)
	if err != nil {
		handleError(w, err)
		return
	}
	defer zincResponse.Body.Close()

	body, err := io.ReadAll(zincResponse.Body)
	if err != nil {
		handleError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(zincResponse.StatusCode)
	_, err = w.Write(body)
	if err != nil {
		fmt.Println(err)
	}
}

func handleError(w http.ResponseWriter, err error) {
	var statusCode int
	switch {
	case strings.Contains(err.Error(), "timeout"):
		statusCode = http.StatusGatewayTimeout
	case strings.Contains(err.Error(), "URL"):
		statusCode = http.StatusBadRequest
	default:
		statusCode = http.StatusInternalServerError
	}
	http.Error(w, err.Error(), statusCode)
}
