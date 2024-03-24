package main

import (
	"fmt"
	"net/http"

	"github.com/JuanFelipeJaramillo20/email_indexer/utils"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
)

func main() {
	router := chi.NewRouter()

	cors := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		AllowCredentials: true,
		MaxAge:           300,
	})

	router.Use(cors.Handler)

	router.Get("/search", utils.SearchHandler)
	router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("hello world"))
	})
	initText := "Server started on :8080"
	fmt.Println(initText)

	listenAndServeText := fmt.Sprintf(":%d", 8080)
	http.ListenAndServe(listenAndServeText, router)
}
