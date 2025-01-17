package auth

import (
	"errors"
	"regexp"
	"time"

	"github.com/google/uuid"
)

var (
	ErrInvalidEmail    = errors.New("invalid email format")
	ErrInvalidPassword = errors.New("password must be at least 8 characters")
	ErrInvalidName     = errors.New("name must be between 2 and 50 characters")
)

// User represents the user entity in our domain
type User struct {
	ID        string    `json:"id"`
	Email     string    `json:"email"`
	Name      string    `json:"name"`
	Password  string    `json:"-"` // "-" means this field won't be included in JSON
	Genres    []string  `json:"genres"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// NewUser creates a new user with validation
func NewUser(email, name, password string, genres []string) (*User, error) {
	if err := validateEmail(email); err != nil {
		return nil, err
	}
	if err := validatePassword(password); err != nil {
		return nil, err
	}
	if err := validateName(name); err != nil {
		return nil, err
	}

	now := time.Now().UTC()
	return &User{
		ID:        uuid.New().String(),
		Email:     email,
		Name:      name,
		Password:  password,
		Genres:    genres,
		CreatedAt: now,
		UpdatedAt: now,
	}, nil
}

// validateEmail checks if the email format is valid
func validateEmail(email string) error {
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
	if !emailRegex.MatchString(email) {
		return ErrInvalidEmail
	}
	return nil
}

// validatePassword checks if the password meets security requirements
func validatePassword(password string) error {
	if len(password) < 8 {
		return ErrInvalidPassword
	}
	// Add more password requirements as needed
	return nil
}

// validateName checks if the name is valid
func validateName(name string) error {
	if len(name) < 2 || len(name) > 50 {
		return ErrInvalidName
	}
	return nil
}

// UpdatePassword updates the user's password with validation
func (u *User) UpdatePassword(password string) error {
	if err := validatePassword(password); err != nil {
		return err
	}
	u.Password = password
	u.UpdatedAt = time.Now().UTC()
	return nil
}

// UpdateEmail updates the user's email with validation
func (u *User) UpdateEmail(email string) error {
	if err := validateEmail(email); err != nil {
		return err
	}
	u.Email = email
	u.UpdatedAt = time.Now().UTC()
	return nil
}

// UpdateName updates the user's name with validation
func (u *User) UpdateName(name string) error {
	if err := validateName(name); err != nil {
		return err
	}
	u.Name = name
	u.UpdatedAt = time.Now().UTC()
	return nil
}
