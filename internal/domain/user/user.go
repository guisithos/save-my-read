package user

import (
	"errors"
	"time"

	"golang.org/x/crypto/bcrypt"
)

// Genre represents a book genre
type Genre string

// User represents the user domain entity
type User struct {
	ID        string
	Name      string
	Email     string
	Password  string
	Genres    []Genre
	CreatedAt time.Time
	UpdatedAt time.Time
}

// NewUser creates a new user with validated fields
func NewUser(name, email, password string, genres []Genre) (*User, error) {
	if name == "" {
		return nil, errors.New("name cannot be empty")
	}
	if email == "" {
		return nil, errors.New("email cannot be empty")
	}
	if password == "" {
		return nil, errors.New("password cannot be empty")
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	return &User{
		Name:      name,
		Email:     email,
		Password:  string(hashedPassword),
		Genres:    genres,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}, nil
}

// ValidatePassword checks if the provided password matches the user's password
func (u *User) ValidatePassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	return err == nil
}
