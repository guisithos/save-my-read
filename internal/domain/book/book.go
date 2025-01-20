package book

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

// Status represents the reading status of a book
type Status string

const (
	StatusToRead    Status = "TO_READ"
	StatusReading   Status = "READING"
	StatusCompleted Status = "COMPLETED"
	StatusDNF       Status = "DNF" // Did Not Finish
)

// IsValid checks if the status is one of the valid statuses
func (s Status) IsValid() bool {
	switch s {
	case StatusToRead, StatusReading, StatusCompleted, StatusDNF:
		return true
	default:
		return false
	}
}

// Book represents the book domain entity
type Book struct {
	ID          string    `json:"id"`
	GoogleID    string    `json:"google_id"`
	Title       string    `json:"title"`
	Authors     []string  `json:"authors"`
	Description string    `json:"description"`
	Categories  []string  `json:"categories"`
	ImageURL    string    `json:"image_url"`
	Status      Status    `json:"status"`
	UserID      string    `json:"user_id"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// NewBook creates a new book with validated fields
func NewBook(googleID, title string, authors []string, description string,
	categories []string, imageURL string, status Status, userID string) (*Book, error) {

	if googleID == "" {
		return nil, errors.New("google_id is required")
	}
	if title == "" {
		return nil, errors.New("title is required")
	}
	if len(authors) == 0 {
		return nil, errors.New("at least one author is required")
	}
	if userID == "" {
		return nil, errors.New("user_id is required")
	}
	if !status.IsValid() {
		return nil, errors.New("invalid status")
	}

	now := time.Now()
	return &Book{
		ID:          uuid.New().String(),
		GoogleID:    googleID,
		Title:       title,
		Authors:     authors,
		Description: description,
		Categories:  categories,
		ImageURL:    imageURL,
		Status:      status,
		UserID:      userID,
		CreatedAt:   now,
		UpdatedAt:   now,
	}, nil
}

// UpdateStatus changes the book's reading status
func (b *Book) UpdateStatus(status Status) error {
	if !status.IsValid() {
		return errors.New("invalid status")
	}
	b.Status = status
	b.UpdatedAt = time.Now()
	return nil
}
