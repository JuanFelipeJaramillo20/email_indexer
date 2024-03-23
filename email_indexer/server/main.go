package main

import (
	"fmt"
	"net/http"

	"github.com/JuanFelipeJaramillo20/email_indexer/utils"
	"github.com/go-chi/chi/v5"
)

func main() {
	router := chi.NewRouter()

	router.Get("/search", utils.SearchHandler)
	router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("hello world"))
	})
	initText := "Server started on :8080"
	fmt.Println(initText)

	listenAndServeText := fmt.Sprintf(":%d", 8080)
	http.ListenAndServe(listenAndServeText, router)
}
