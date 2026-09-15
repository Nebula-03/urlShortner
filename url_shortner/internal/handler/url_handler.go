package handler

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strings"

	"github.com/jackc/pgx/v5"

	"url_shortner/internal/service"
)

var aliasPattern = regexp.MustCompile(`^[A-Za-z0-9_-]+$`)

type URLHandler struct {
	service *service.URLService
}

type CreateURLRequest struct {
	Alias       string `json:"alias"`
	OriginalURL string `json:"original_url"`
}

type BasicURLResponse struct {
	Message  string `json:"message"`
	ShortURL string `json:"short_url"`
}

type CustomAliasResponse struct {
	Message     string `json:"message"`
	CustomAlias string `json:"custom_alias"`
}

func NewURLHandler(service *service.URLService) *URLHandler {
	return &URLHandler{
		service: service,
	}
}

func (h *URLHandler) CreateURL(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {

		http.Error(
			w,
			"Method not allowed",
			http.StatusMethodNotAllowed,
		)

		return
	}

	var request CreateURLRequest

	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {

		http.Error(
			w,
			"Invalid request body",
			http.StatusBadRequest,
		)

		return
	}

	var validationErrors []string

	if strings.TrimSpace(request.OriginalURL) == "" &&
		strings.TrimSpace(request.Alias) == "" {

		http.Error(
			w,
			"Please enter the URL and alias name you want",
			http.StatusBadRequest,
		)

		return
	}

	if request.Alias != "" {

		if len(request.Alias) > 100 {
			validationErrors = append(
				validationErrors,
				"Alias must be 100 characters or less",
			)
		}

		if !aliasPattern.MatchString(request.Alias) {

			validationErrors = append(
				validationErrors,
				"Invalid alias! Use only letters, numbers, hyphens, and underscores",
			)
		}
	}

	if request.OriginalURL == "" {

		validationErrors = append(
			validationErrors,
			"Please enter the URL",
		)

	} else {

		parsedURL, err := url.ParseRequestURI(request.OriginalURL)
		if err != nil ||
			(parsedURL.Scheme != "http" &&
				parsedURL.Scheme != "https") ||
			parsedURL.Host == "" {

			validationErrors = append(
				validationErrors,
				"Please enter a valid HTTP or HTTPS URL",
			)
		}
	}

	if len(validationErrors) > 0 {

		http.Error(
			w,
			strings.Join(validationErrors, "\n"),
			http.StatusBadRequest,
		)

		return
	}

	url, err := h.service.CreateURL(
		r.Context(),
		request.Alias,
		request.OriginalURL,
	)

	if err != nil {
		log.Println("Create URL error:", err)

		if errors.Is(err, service.ErrURLUnreachable) {

			http.Error(
				w,
				"Please enter a valid and reachable URL",
				http.StatusBadRequest,
			)

			return
		}

		if errors.Is(err, service.ErrAliasAlreadyExists) {

			http.Error(
				w,
				"This custom alias is already taken. Please choose a different alias.",
				http.StatusConflict,
			)

			return
		}

		http.Error(
			w,
			"Failed to create URL",
			http.StatusInternalServerError,
		)

		return
	}

	w.Header().Set(
		"Content-Type",
		"application/json",
	)

	w.WriteHeader(http.StatusCreated)

	if request.Alias == "" {

		baseURL := os.Getenv("BASE_URL")

		if baseURL == "" {
			baseURL = "http://localhost:8080"
		}

		response := BasicURLResponse{
			Message:  "Your shortened URL is ready!",
			ShortURL: strings.TrimRight(baseURL, "/") + "/" + url.Alias,
		}

		json.NewEncoder(w).Encode(response)

		return
	}

	response := CustomAliasResponse{
		Message:     "Your custom alias is ready!",
		CustomAlias: url.Alias,
	}

	json.NewEncoder(w).Encode(response)
}

func (h *URLHandler) RedirectURL(
	w http.ResponseWriter,
	r *http.Request,
) {

	if r.Method != http.MethodGet {

		http.Error(
			w,
			"Method not allowed",
			http.StatusMethodNotAllowed,
		)

		return
	}

	alias := strings.TrimPrefix(
		r.URL.Path,
		"/",
	)

	if alias == "" {

		http.Error(
			w,
			"Alias is required",
			http.StatusBadRequest,
		)

		return
	}

	if strings.Contains(alias, "/") {

		http.Error(
			w,
			"Invalid alias",
			http.StatusBadRequest,
		)

		return
	}

	url, err := h.service.GetURLByAlias(
		r.Context(),
		alias,
	)

	if err != nil {

		log.Println("Get URL error:", err)

		if errors.Is(err, pgx.ErrNoRows) {

			http.Error(
				w,
				"URL not found",
				http.StatusNotFound,
			)

			return
		}

		http.Error(
			w,
			"Failed to retrieve URL",
			http.StatusInternalServerError,
		)

		return
	}

	http.Redirect(
		w,
		r,
		url.OriginalURL,
		http.StatusFound,
	)
}
