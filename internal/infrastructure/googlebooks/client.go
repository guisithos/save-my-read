package googlebooks

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
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
	log.Printf("Making request to Google Books API with query: %s", query)
	url := fmt.Sprintf("%s/volumes?q=%s&key=%s", c.baseURL, url.QueryEscape(query), c.apiKey)

	resp, err := c.httpClient.Get(url)
	if err != nil {
		log.Printf("Google Books API request failed: %v", err)
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := ioutil.ReadAll(resp.Body)
		log.Printf("Google Books API error response: %s", string(body))
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var result BookResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		log.Printf("Failed to decode response: %v", err)
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	log.Printf("Successfully retrieved %d books", len(result.Items))
	return &result, nil
}
