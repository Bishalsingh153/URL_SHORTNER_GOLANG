package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

// ShortenRequest represents the JSON data we send to our server.
type ShortenRequest struct {
	URL string `json:"url"`
}

// ShortenResponse represents the JSON response we receive from our server.
type ShortenResponse struct {
	ShortURL string `json:"short_url"`
}

// shortenURL sends the original URL to our API and returns the shortened URL.
func shortenURL(originalURL string) (string, error) {

	// Convert our Go struct into JSON.
	requestData := ShortenRequest{
		URL: originalURL,
	}

	jsonData, err := json.Marshal(requestData)
	if err != nil {
		return "", fmt.Errorf("error creating JSON: %w", err)
	}

	// Send a POST request containing our JSON data to the server.
	resp, err := http.Post(
		"http://localhost:8080/shorten",
		"application/json",
		bytes.NewBuffer(jsonData),
	)

	if err != nil {
		return "", fmt.Errorf("error connecting to server: %w", err)
	}

	// Make sure the response body is closed when this function finishes.
	defer resp.Body.Close()

	// Check whether the server returned a successful HTTP status.
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		errorBody, err := io.ReadAll(resp.Body)
		if err != nil {
			return "", fmt.Errorf("server returned %s", resp.Status)
		}

		return "", fmt.Errorf("server error (%s): %s",
			resp.Status, strings.TrimSpace(string(errorBody)))
	}

	// Decode the JSON response into our ShortenResponse struct.
	var result ShortenResponse

	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return "", fmt.Errorf("error decoding server response: %w", err)
	}

	// Return the shortened URL to whoever called this function.
	return result.ShortURL, nil
}

func main() {

	// Create a reader that accepts input from the terminal.
	reader := bufio.NewReader(os.Stdin)

	fmt.Println("=================================")
	fmt.Println("        GO URL SHORTENER")
	fmt.Println("=================================")

	// Keep running until the user chooses to exit.
	for {

		fmt.Println()
		fmt.Println("Enter a long URL (or type 'q' to quit):")

		// Read the URL entered by the user.
		originalURL, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Error reading input:", err)
			return
		}

		// Remove whitespace and newline characters from the input.
		originalURL = strings.TrimSpace(originalURL)

		// Allow the user to exit the application.
		if strings.EqualFold(originalURL, "q") {
			fmt.Println("Thank you for using Go URL Shortener!")
			return
		}

		// Reject empty input before making an HTTP request.
		if originalURL == "" {
			fmt.Println("Error: URL cannot be empty.")
			continue
		}

		// Send the URL to our server and receive the shortened URL.
		shortURL, err := shortenURL(originalURL)
		if err != nil {
			fmt.Println("Error:", err)
			continue
		}

		// Display the shortened URL.
		fmt.Println()
		fmt.Println("Your shortened URL:", shortURL)
		fmt.Println("---------------------------------")
	}
}
