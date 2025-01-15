package application

import (
	"errors"
	"fmt"

	"github.com/guisithos/save-my-read/internal/domain/auth"
	"github.com/guisithos/save-my-read/internal/domain/user"
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

func (s *AuthService) Register(email, password, name string, genres []string) (*user.User, error) {
	// Check if user already exists
	existing, _ := s.userRepo.FindByEmail(email)
	if existing != nil {
		return nil, errors.New("email already registered")
	}

	// Create new user
	newUser, err := user.NewUser(email, password, name, genres)
	if err != nil {
		return nil, err
	}

	// Save user
	if err := s.userRepo.Save(newUser); err != nil {
		return nil, err
	}

	return newUser, nil
}

func (s *AuthService) Login(email, password string) (*auth.LoginResponse, error) {
	user, err := s.userRepo.FindByEmail(email)
	if err != nil {
		return nil, errors.New("invalid credentials")
	}

	if !user.ValidatePassword(password) {
		return nil, errors.New("invalid credentials")
	}

	token, err := s.tokenService.GenerateToken(user.ID, user.Email)
	if err != nil {
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
