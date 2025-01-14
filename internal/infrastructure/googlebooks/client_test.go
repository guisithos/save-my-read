package googlebooks

import (
	"os"
	"testing"
)

func TestSearchBooks(t *testing.T) {
	// Skip if no API key is set
	if os.Getenv("GOOGLE_BOOKS_API_KEY") == "" {
		t.Skip("GOOGLE_BOOKS_API_KEY not set")
	}

	client, err := NewClient()
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	tests := []struct {
		name    string
		query   string
		wantErr bool
	}{
		{
			name:    "Valid search",
			query:   "The Lord of the Rings",
			wantErr: false,
		},
		{
			name:    "Empty query",
			query:   "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			response, err := client.SearchBooks(tt.query)
			if (err != nil) != tt.wantErr {
				t.Errorf("SearchBooks() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && (response == nil || len(response.Items) == 0) {
				t.Error("SearchBooks() returned no results for valid query")
			}
		})
	}
}
