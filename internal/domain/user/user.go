package user

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

// Genre represents a book genre
type Genre string

// User represents the user domain entity
type User struct {
	ID        string
	Email     string
	Password  string
	Name      string
	Genres    []string
	CreatedAt time.Time
	UpdatedAt time.Time
}

// NewUser creates a new user with validated fields
func NewUser(email, password, name string, genres []string) (*User, error) {
	if email == "" {
		return nil, errors.New("email cannot be empty")
	}
	if password == "" {
		return nil, errors.New("password cannot be empty")
	}
	if name == "" {
		return nil, errors.New("name cannot be empty")
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	now := time.Now()
	return &User{
		ID:        uuid.New().String(),
		Email:     email,
		Password:  string(hashedPassword),
		Name:      name,
		Genres:    genres,
		CreatedAt: now,
		UpdatedAt: now,
	}, nil
}

// ValidatePassword checks if the provided password matches the user's password
func (u *User) ValidatePassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	return err == nil
}
