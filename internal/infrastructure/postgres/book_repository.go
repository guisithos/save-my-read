package postgres

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/guisithos/save-my-read/internal/domain/book"
	_ "github.com/lib/pq"
)

// BookRepository implements the book.Repository interface using PostgreSQL
type BookRepository struct {
	db *sql.DB
}

// NewBookRepository creates a new PostgreSQL book repository
func NewBookRepository(db *sql.DB) *BookRepository {
	return &BookRepository{db: db}
}

// Save stores a new book in the database
func (r *BookRepository) Save(book *book.Book) error {
	query := `
		INSERT INTO books (
			google_id, title, authors, description, categories, 
			image_url, status, user_id, created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		RETURNING id`

	err := r.db.QueryRow(
		query,
		book.GoogleID,
		book.Title,
		book.Authors,
		book.Description,
		book.Categories,
		book.ImageURL,
		book.Status,
		book.UserID,
		book.CreatedAt,
		book.UpdatedAt,
	).Scan(&book.ID)

	if err != nil {
		return fmt.Errorf("error saving book: %w", err)
	}

	return nil
}

// FindByID retrieves a book by its ID
func (r *BookRepository) FindByID(id string) (*book.Book, error) {
	query := `
		SELECT id, google_id, title, authors, description, categories,
			   image_url, status, user_id, created_at, updated_at
		FROM books WHERE id = $1`

	b := &book.Book{}
	err := r.db.QueryRow(query, id).Scan(
		&b.ID, &b.GoogleID, &b.Title, &b.Authors, &b.Description,
		&b.Categories, &b.ImageURL, &b.Status, &b.UserID,
		&b.CreatedAt, &b.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("book not found")
	}
	if err != nil {
		return nil, fmt.Errorf("error finding book: %w", err)
	}

	return b, nil
}

// FindByUserID retrieves all books for a user
func (r *BookRepository) FindByUserID(userID string) ([]*book.Book, error) {
	query := `
		SELECT id, google_id, title, authors, description, categories,
			   image_url, status, user_id, created_at, updated_at
		FROM books WHERE user_id = $1`

	rows, err := r.db.Query(query, userID)
	if err != nil {
		return nil, fmt.Errorf("error querying books: %w", err)
	}
	defer rows.Close()

	var books []*book.Book
	for rows.Next() {
		b := &book.Book{}
		err := rows.Scan(
			&b.ID, &b.GoogleID, &b.Title, &b.Authors, &b.Description,
			&b.Categories, &b.ImageURL, &b.Status, &b.UserID,
			&b.CreatedAt, &b.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning book: %w", err)
		}
		books = append(books, b)
	}

	return books, nil
}

// FindByUserIDAndStatus retrieves books for a user with specific status
func (r *BookRepository) FindByUserIDAndStatus(userID string, status book.Status) ([]*book.Book, error) {
	query := `
		SELECT id, google_id, title, authors, description, categories,
			   image_url, status, user_id, created_at, updated_at
		FROM books WHERE user_id = $1 AND status = $2`

	rows, err := r.db.Query(query, userID, status)
	if err != nil {
		return nil, fmt.Errorf("error querying books: %w", err)
	}
	defer rows.Close()

	var books []*book.Book
	for rows.Next() {
		b := &book.Book{}
		err := rows.Scan(
			&b.ID, &b.GoogleID, &b.Title, &b.Authors, &b.Description,
			&b.Categories, &b.ImageURL, &b.Status, &b.UserID,
			&b.CreatedAt, &b.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning book: %w", err)
		}
		books = append(books, b)
	}

	return books, nil
}

// Update updates an existing book
func (r *BookRepository) Update(book *book.Book) error {
	query := `
		UPDATE books 
		SET status = $1, updated_at = $2
		WHERE id = $3`

	result, err := r.db.Exec(query, book.Status, time.Now(), book.ID)
	if err != nil {
		return fmt.Errorf("error updating book: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error checking rows affected: %w", err)
	}

	if rows == 0 {
		return fmt.Errorf("book not found")
	}

	return nil
}

// Delete removes a book from the database
func (r *BookRepository) Delete(id string) error {
	query := `DELETE FROM books WHERE id = $1`

	result, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("error deleting book: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error checking rows affected: %w", err)
	}

	if rows == 0 {
		return fmt.Errorf("book not found")
	}

	return nil
}
