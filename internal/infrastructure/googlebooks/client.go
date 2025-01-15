package googlebooks

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"

	"github.com/joho/godotenv"
)

// Client handles communication with Google Books AP
type Client struct {
	apiKey     string
	baseURL    string
	httpClient *http.Client
}

// BookResponse represents the Google Books API response structure
type BookResponse struct {
	Items []struct {
		ID         string `json:"id"`
		VolumeInfo struct {
			Title       string   `json:"title"`
			Authors     []string `json:"authors"`
			Description string   `json:"description"`
			Categories  []string `json:"categories"`
			ImageLinks  struct {
				Thumbnail string `json:"thumbnail"`
			} `json:"imageLinks"`
		} `json:"volumeInfo"`
	} `json:"items"`
}

// NewClient creates a new Google Books API client
func NewClient() (*Client, error) {
	// Try to load from .env file if godotenv is available
	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: .env file not found or error loading it")
	}

	apiKey := os.Getenv("GOOGLE_BOOKS_API_KEY")
	if apiKey == "" {
		return nil, fmt.Errorf("GOOGLE_BOOKS_API_KEY environment variable is not set")
	}

	return &Client{
		apiKey:     apiKey,
		baseURL:    "https://www.googleapis.com/books/v1",
		httpClient: &http.Client{},
	}, nil
}

// SearchBooks searches for books using the Google Books API
func (c *Client) SearchBooks(query string) (*BookResponse, error) {
	// Build URL with query parameters
	endpoint := fmt.Sprintf("%s/volumes", c.baseURL)
	params := url.Values{}
	params.Add("q", query)
	params.Add("key", c.apiKey)

	// Create request
	req, err := http.NewRequest("GET", endpoint+"?"+params.Encode(), nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	// Execute request
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error executing request: %w", err)
	}
	defer resp.Body.Close()

	// Check status code
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API request failed with status: %d", resp.StatusCode)
	}

	// Parse response
	var bookResponse BookResponse
	if err := json.NewDecoder(resp.Body).Decode(&bookResponse); err != nil {
		return nil, fmt.Errorf("error decoding response: %w", err)
	}

	return &bookResponse, nil
}
