package main

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"net/http"
)

func main() {

	fmt.Println("Server starting on port 8080")

	urls := make(map[string]string)

	//routes
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {

		fmt.Fprintln(w, "Welcome to my URL shortner")

	})

	http.HandleFunc("/shorten", func(w http.ResponseWriter, r *http.Request) {
		originalUrl := r.URL.Query().Get("url")
		if originalUrl == "" {
			http.Error(w, "URL is required", http.StatusBadRequest)
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

		fmt.Fprintln(w, "Your short url is: http://localhost:8080/"+shortCode)

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
