package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/guisithos/save-my-read/internal/application"
	"github.com/guisithos/save-my-read/internal/domain/book"
	"github.com/guisithos/save-my-read/internal/infrastructure/googlebooks"
)

// BookHandler handles HTTP requests for book operations
type BookHandler struct {
	bookService  *application.BookService
	googleClient *googlebooks.Client
}

// NewBookHandler creates a new BookHandler
func NewBookHandler(bookService *application.BookService, googleClient *googlebooks.Client) *BookHandler {
	return &BookHandler{
		bookService:  bookService,
		googleClient: googleClient,
	}
}

// SearchBooks handles book search requests
func (h *BookHandler) SearchBooks(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	query := r.URL.Query().Get("q")
	if query == "" {
		http.Error(w, "Query parameter 'q' is required", http.StatusBadRequest)
		return
	}

	log.Printf("Searching for books with query: %s", query)
	books, err := h.googleClient.SearchBooks(query)
	if err != nil {
		log.Printf("Search error: %v", err)
		http.Error(w, "Failed to search books", http.StatusInternalServerError)
		return
	}

	log.Printf("Found %d books", len(books.Items))
	respondJSON(w, http.StatusOK, books)
}

// AddToList handles adding a book to user's list
func (h *BookHandler) AddToList(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		UserID      string      `json:"user_id"`
		GoogleID    string      `json:"google_id"`
		Title       string      `json:"title"`
		Authors     []string    `json:"authors"`
		Description string      `json:"description"`
		Categories  []string    `json:"categories"`
		ImageURL    string      `json:"image_url"`
		Status      book.Status `json:"status"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	newBook, err := h.bookService.AddBookToList(
		req.UserID,
		req.GoogleID,
		req.Title,
		req.Authors,
		req.Description,
		req.Categories,
		req.ImageURL,
		req.Status,
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	respondJSON(w, http.StatusCreated, newBook)
}

// GetUserBooks handles retrieving user's books
func (h *BookHandler) GetUserBooks(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	userID := r.URL.Query().Get("user_id")
	if userID == "" {
		http.Error(w, "user_id is required", http.StatusBadRequest)
		return
	}

	status := r.URL.Query().Get("status")
	var books []*book.Book
	var err error

	if status != "" {
		books, err = h.bookService.GetUserBooksByStatus(userID, book.Status(status))
	} else {
		books, err = h.bookService.GetUserBooks(userID)
	}

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	respondJSON(w, http.StatusOK, books)
}

// UpdateBookStatus handles updating a book's status
func (h *BookHandler) UpdateBookStatus(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		BookID string      `json:"book_id"`
		Status book.Status `json:"status"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := h.bookService.UpdateBookStatus(req.BookID, req.Status); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func respondJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}
