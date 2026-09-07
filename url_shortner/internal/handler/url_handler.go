package handler

import (
	"encoding/json"
	"net/http"

	"url_shortner/internal/service"
)

type URLHandler struct {
	service *service.URLService
}

type CreateURLRequest struct {
	Alias       string `json:"alias"`
	OriginalURL string `json:"original_url"`
}

func NewURLHandler(service *service.URLService) *URLHandler {
	return &URLHandler{
		service: service,
	}
}

func (h *URLHandler) CreateURL(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var request CreateURLRequest

	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	w.Write([]byte("URL request received successfully!"))
}
