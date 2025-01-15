package postgres

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/guisithos/save-my-read/internal/domain/user"
	"github.com/lib/pq"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Save(user *user.User) error {
	query := `
		INSERT INTO users (id, email, password_hash, name, genres, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)`

	_, err := r.db.Exec(
		query,
		user.ID,
		user.Email,
		user.Password,
		user.Name,
		pq.Array(user.Genres),
		user.CreatedAt,
		user.UpdatedAt,
	)

	if err != nil {
		// Check for unique constraint violation
		if pqErr, ok := err.(*pq.Error); ok && pqErr.Code == "23505" {
			return errors.New("email already exists")
		}
		return fmt.Errorf("error saving user: %w", err)
	}

	return nil
}

func (r *UserRepository) FindByEmail(email string) (*user.User, error) {
	query := `
		SELECT id, email, password_hash, name, genres, created_at, updated_at
		FROM users
		WHERE email = $1`

	u := &user.User{}
	var genres []string

	err := r.db.QueryRow(query, email).Scan(
		&u.ID,
		&u.Email,
		&u.Password,
		&u.Name,
		pq.Array(&genres),
		&u.CreatedAt,
		&u.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, errors.New("user not found")
	}
	if err != nil {
		return nil, fmt.Errorf("error finding user: %w", err)
	}

	u.Genres = genres
	return u, nil
}

func (r *UserRepository) FindByID(id string) (*user.User, error) {
	query := `
		SELECT id, email, password_hash, name, genres, created_at, updated_at
		FROM users
		WHERE id = $1`

	u := &user.User{}
	var genres []string

	err := r.db.QueryRow(query, id).Scan(
		&u.ID,
		&u.Email,
		&u.Password,
		&u.Name,
		pq.Array(&genres),
		&u.CreatedAt,
		&u.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, errors.New("user not found")
	}
	if err != nil {
		return nil, fmt.Errorf("error finding user: %w", err)
	}

	u.Genres = genres
	return u, nil
}

func (r *UserRepository) Update(user *user.User) error {
	query := `
		UPDATE users 
		SET name = $1, email = $2, password = $3, genres = $4, updated_at = $5
		WHERE id = $6`

	result, err := r.db.Exec(query, user.Name, user.Email, user.Password, user.Genres, user.UpdatedAt, user.ID)
	if err != nil {
		return fmt.Errorf("error updating user: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error checking rows affected: %w", err)
	}

	if rows == 0 {
		return fmt.Errorf("user not found")
	}

	return nil
}
