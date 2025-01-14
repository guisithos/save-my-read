package book

import (
	"errors"
	"time"
)

// Status represents the reading status of a book
type Status string

const (
	StatusRead     Status = "READ"
	StatusWishlist Status = "WISHLIST"
)

// Book represents the book domain entity
type Book struct {
	ID          string
	GoogleID    string // ID from Google Books API
	Title       string
	Authors     []string
	Description string
	Categories  []string
	ImageURL    string
	Status      Status
	UserID      string // Reference to the user who added this book
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// NewBook creates a new book with validated fields
func NewBook(googleID, title string, authors []string, description string,
	categories []string, imageURL string, status Status, userID string) (*Book, error) {

	if googleID == "" {
		return nil, errors.New("google book ID cannot be empty")
	}
	if title == "" {
		return nil, errors.New("title cannot be empty")
	}
	if len(authors) == 0 {
		return nil, errors.New("book must have at least one author")
	}
	if userID == "" {
		return nil, errors.New("user ID cannot be empty")
	}
	if !isValidStatus(status) {
		return nil, errors.New("invalid book status")
	}

	return &Book{
		GoogleID:    googleID,
		Title:       title,
		Authors:     authors,
		Description: description,
		Categories:  categories,
		ImageURL:    imageURL,
		Status:      status,
		UserID:      userID,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}, nil
}

func isValidStatus(s Status) bool {
	return s == StatusRead || s == StatusWishlist
}

// UpdateStatus changes the book's reading status
func (b *Book) UpdateStatus(status Status) error {
	if !isValidStatus(status) {
		return errors.New("invalid book status")
	}
	b.Status = status
	b.UpdatedAt = time.Now()
	return nil
}
