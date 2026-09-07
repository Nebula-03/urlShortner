package main

import (
	"fmt"
	"log"
	"net/http"

	"url_shortner/internal/config"
	"url_shortner/internal/handler"
	"url_shortner/internal/repository"
	"url_shortner/internal/service"
)

func main() {
	fmt.Println("URL Shortner starting...")

	db, err := config.ConnectDatabase()
	if err != nil {
		log.Fatal("Database connection failed:", err)
	}

	defer db.Close()

	fmt.Println("Database connected successfully!")

	urlRepository := repository.NewURLRepository(db)

	urlService := service.NewURLService(urlRepository)

	urlHandler := handler.NewURLHandler(urlService)

	http.HandleFunc("/shorten", urlHandler.CreateURL)

	fmt.Println("Server running on http://localhost:8080")

	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("Server failed:", err)
	}
}
