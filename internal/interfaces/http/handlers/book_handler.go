package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/guisithos/save-my-read/internal/application"
	"github.com/guisithos/save-my-read/internal/domain/book"
	"github.com/guisithos/save-my-read/internal/infrastructure/googlebooks"
	"github.com/guisithos/save-my-read/internal/interfaces/http/middleware"
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

// AddBook handles adding a book to user's list
func (h *BookHandler) AddBook(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Get user ID from context
	userID, ok := r.Context().Value(middleware.UserIDKey).(string)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var req struct {
		GoogleBookID string   `json:"google_book_id"`
		Title        string   `json:"title"`
		Authors      []string `json:"authors"`
		Description  string   `json:"description"`
		Categories   []string `json:"categories"`
		ImageURL     string   `json:"image_url"`
		Status       string   `json:"status"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	status := book.Status(req.Status)
	if !status.IsValid() {
		http.Error(w, "Invalid book status", http.StatusBadRequest)
		return
	}

	newBook, err := h.bookService.AddBookToList(
		userID,
		req.GoogleBookID,
		req.Title,
		req.Authors,
		req.Description,
		req.Categories,
		req.ImageURL,
		status,
	)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"data":    newBook,
	})
}

// GetBooks handles retrieving user's books
func (h *BookHandler) GetBooks(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Get user ID from context
	userID, ok := r.Context().Value(middleware.UserIDKey).(string)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	status := r.URL.Query().Get("status")
	var books []*book.Book
	var err error

	if status != "" {
		bookStatus := book.Status(status)
		if !bookStatus.IsValid() {
			http.Error(w, "Invalid status parameter", http.StatusBadRequest)
			return
		}
		books, err = h.bookService.GetUserBooksByStatus(userID, bookStatus)
	} else {
		books, err = h.bookService.GetUserBooks(userID)
	}

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"data":    books,
	})
}

// UpdateBookStatus handles updating a book's status
func (h *BookHandler) UpdateBookStatus(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Get user ID from context (for future authorization checks)
	_, ok := r.Context().Value(middleware.UserIDKey).(string)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var req struct {
		BookID string `json:"book_id"`
		Status string `json:"status"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	status := book.Status(req.Status)
	if !status.IsValid() {
		http.Error(w, "Invalid book status", http.StatusBadRequest)
		return
	}

	if err := h.bookService.UpdateBookStatus(req.BookID, status); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
	})
}

func respondJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}
