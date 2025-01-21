package application

import (
	"errors"
	"fmt"

	"github.com/guisithos/save-my-read/internal/domain/auth"
	"github.com/guisithos/save-my-read/internal/domain/user"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	userRepo     user.Repository
	tokenService auth.TokenService
}

func NewAuthService(userRepo user.Repository, tokenService auth.TokenService) *AuthService {
	return &AuthService{
		userRepo:     userRepo,
		tokenService: tokenService,
	}
}

func (s *AuthService) Register(email, password, name string, genres []string) (*auth.LoginResponse, error) {
	fmt.Printf("Starting registration process for email: %s, name: %s\n", email, name)

	// Check if user already exists
	existing, err := s.userRepo.FindByEmail(email)
	if err != nil {
		fmt.Printf("Error checking existing user: %v\n", err)
	}
	if existing != nil {
		fmt.Println("User already exists with this email")
		return nil, errors.New("email already registered")
	}

	// Create new user
	fmt.Println("Creating new user...")
	newUser, err := user.NewUser(email, password, name, genres)
	if err != nil {
		fmt.Printf("Error creating new user: %v\n", err)
		return nil, fmt.Errorf("failed to create user: %w", err)
	}
	fmt.Printf("User created with ID: %s\n", newUser.ID)

	// Save user
	fmt.Println("Saving user to database...")
	if err := s.userRepo.Save(newUser); err != nil {
		fmt.Printf("Error saving user to database: %v\n", err)
		return nil, fmt.Errorf("failed to save user: %w", err)
	}
	fmt.Println("User saved successfully")

	// Generate token
	fmt.Println("Generating authentication token...")
	token, err := s.tokenService.GenerateToken(newUser.ID, newUser.Email)
	if err != nil {
		fmt.Printf("Error generating token: %v\n", err)
		return nil, fmt.Errorf("failed to generate token: %w", err)
	}
	fmt.Println("Token generated successfully")

	return &auth.LoginResponse{
		Token: token,
		User: auth.UserResponse{
			ID:    newUser.ID,
			Email: newUser.Email,
			Name:  newUser.Name,
		},
	}, nil
}

func (s *AuthService) Login(email, password string) (*auth.LoginResponse, error) {
	// Find user by email
	user, err := s.userRepo.FindByEmail(email)
	if err != nil {
		// Log the error but don't expose internal details
		fmt.Printf("Error finding user: %v\n", err)
		return nil, auth.ErrInvalidCredentials
	}

	// Validate password using bcrypt
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		fmt.Printf("Invalid password for user %s\n", email)
		return nil, auth.ErrInvalidCredentials
	}

	// Generate JWT token
	token, err := s.tokenService.GenerateToken(user.ID, user.Email)
	if err != nil {
		fmt.Printf("Error generating token: %v\n", err)
		return nil, fmt.Errorf("failed to generate token: %w", err)
	}

	return &auth.LoginResponse{
		Token: token,
		User: auth.UserResponse{
			ID:    user.ID,
			Email: user.Email,
			Name:  user.Name,
		},
	}, nil
}
