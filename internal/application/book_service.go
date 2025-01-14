package application

import (
	"errors"

	"github.com/guisithos/save-my-read/internal/domain/book"
	"github.com/guisithos/save-my-read/internal/domain/user"
)

// BookService handles the application logic for books
type BookService struct {
	bookRepo book.Repository
	userRepo user.Repository
}

// NewBookService creates a new BookService
func NewBookService(bookRepo book.Repository, userRepo user.Repository) *BookService {
	return &BookService{
		bookRepo: bookRepo,
		userRepo: userRepo,
	}
}

// AddBookToList adds a book to user's reading list
func (s *BookService) AddBookToList(userID, googleBookID, title string,
	authors []string, description string, categories []string,
	imageURL string, status book.Status) (*book.Book, error) {

	// Verify user exists
	_, err := s.userRepo.FindByID(userID)
	if err != nil {
		return nil, errors.New("user not found")
	}

	// Create new book
	newBook, err := book.NewBook(
		googleBookID,
		title,
		authors,
		description,
		categories,
		imageURL,
		status,
		userID,
	)
	if err != nil {
		return nil, err
	}

	// Save book
	err = s.bookRepo.Save(newBook)
	if err != nil {
		return nil, err
	}

	return newBook, nil
}

// GetUserBooks retrieves all books for a user
func (s *BookService) GetUserBooks(userID string) ([]*book.Book, error) {
	return s.bookRepo.FindByUserID(userID)
}

// GetUserBooksByStatus retrieves books for a user filtered by status
func (s *BookService) GetUserBooksByStatus(userID string, status book.Status) ([]*book.Book, error) {
	return s.bookRepo.FindByUserIDAndStatus(userID, status)
}

// UpdateBookStatus changes the reading status of a book
func (s *BookService) UpdateBookStatus(bookID string, status book.Status) error {
	book, err := s.bookRepo.FindByID(bookID)
	if err != nil {
		return err
	}

	err = book.UpdateStatus(status)
	if err != nil {
		return err
	}

	return s.bookRepo.Update(book)
}
