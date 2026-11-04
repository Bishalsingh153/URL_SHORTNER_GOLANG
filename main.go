package main

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"math/big"
	"net/http"
	"net/url"
)

type ShortenRequest struct {
	URL string `json:"url"`
}

func main() {

	fmt.Println("Server starting on port 8080")

	urls := make(map[string]string)

	//routes
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {

		fmt.Fprintln(w, "Welcome to my URL shortner")

	})

	http.HandleFunc("/shorten", func(w http.ResponseWriter, r *http.Request) {

		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var requestData ShortenRequest

		err := json.NewDecoder(r.Body).Decode(&requestData)
		if err != nil {
			http.Error(w, "invalid JSON", http.StatusBadRequest)
			return
		}

		originalUrl := requestData.URL

		if originalUrl == "" {
			http.Error(w, "URL is required", http.StatusBadRequest)
			return
		}

		parsedURL, err := url.ParseRequestURI(originalUrl)

		if err != nil || parsedURL.Host == "" || (parsedURL.Scheme != "http" && parsedURL.Scheme != "https") {
			http.Error(w, "Invalid URL", http.StatusBadRequest)
			return
		}

		shortCode := generateShortCode()

		for {
			if _, exists := urls[shortCode]; !exists {
				break
			}

			shortCode = generateShortCode()
		}

		urls[shortCode] = originalUrl

		fmt.Println("Received URL:", originalUrl)

		w.Header().Set("Content-Type", "application/json")

		response := map[string]string{
			"short_url": "http://localhost:8080/" + shortCode,
		}
		json.NewEncoder(w).Encode(response)

	})

	http.HandleFunc("/{shortCode}", func(w http.ResponseWriter, r *http.Request) {
		shortCode := r.PathValue("shortCode")

		originalUrl, exists := urls[shortCode]

		if !exists {
			http.NotFound(w, r)
			return
		}
		http.Redirect(w, r, originalUrl, http.StatusFound)
	})

	http.ListenAndServe(":8080", nil)
}

func generateShortCode() string {
	characters := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

	code := make([]byte, 6)

	for i := range code {

		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(characters))))

		if err != nil {
			panic(err)
		}

		code[i] = characters[num.Int64()]
	}

	return string(code)
}
