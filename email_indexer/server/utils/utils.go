package utils

import (
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

	queryString := r.URL.Query()
	term := queryString.Get("term")
	page := queryString.Get("page")
	order := queryString.Get("order")

	if term == "" {
		http.Error(w, "Search term can't be empty", http.StatusBadRequest)
		fmt.Println("Search term can't empty")
		return
	}

	if page == "" {
		page = "0"
	}
	if len(order) > 2 {
		order = ""
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
		term, order, page)

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
