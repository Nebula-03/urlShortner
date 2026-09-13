package main

import (
	"fmt"
	"net/http"

	"url_shortner/internal/config"
	"url_shortner/internal/handler"
	"url_shortner/internal/repository"
	"url_shortner/internal/service"
)

func enableCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func main() {

	fmt.Println("URL Shortener starting...")

	db, err := config.ConnectDatabase()

	if err != nil {
		fmt.Println("Database connection failed:", err)
		return
	}

	defer db.Close()

	fmt.Println("Database connected successfully!")

	urlRepository := repository.NewURLRepository(db)
	urlService := service.NewURLService(urlRepository)
	urlHandler := handler.NewURLHandler(urlService)

	mux := http.NewServeMux()

	mux.HandleFunc("/shorten", urlHandler.CreateURL)
	mux.HandleFunc("/", urlHandler.RedirectURL)

	fmt.Println("Server running on http://localhost:8080")

	err = http.ListenAndServe(
		":8080",
		enableCORS(mux),
	)

	if err != nil {
		fmt.Println("Server error:", err)
	}
}
