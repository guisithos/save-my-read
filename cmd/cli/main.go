package main

import (
	"fmt"
	"log"
	"os"

	"github.com/guisithos/save-my-read/internal/infrastructure/googlebooks"
)

func main() {
	// Initialize Google Books client
	client, err := googlebooks.NewClient()
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	// Get search query from command line
	if len(os.Args) < 2 {
		log.Fatal("Please provide a search query")
	}
	query := os.Args[1]

	// Search for books
	response, err := client.SearchBooks(query)
	if err != nil {
		log.Fatalf("Search failed: %v", err)
	}

	// Print results
	fmt.Printf("Found %d books:\n\n", len(response.Items))
	for _, item := range response.Items {
		fmt.Printf("Title: %s\n", item.VolumeInfo.Title)
		fmt.Printf("Authors: %v\n", item.VolumeInfo.Authors)
		fmt.Printf("Categories: %v\n", item.VolumeInfo.Categories)
		fmt.Printf("Description: %.100s...\n", item.VolumeInfo.Description)
		fmt.Printf("Thumbnail: %s\n", item.VolumeInfo.ImageLinks.Thumbnail)
		fmt.Println("---")
	}
}
