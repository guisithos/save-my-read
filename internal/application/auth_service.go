package application

import (
	"github.com/google/uuid"
	"github.com/guisithos/save-my-read/internal/domain/auth"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	userRepo     auth.UserRepository
	tokenService auth.TokenService
}

func NewAuthService(userRepo auth.UserRepository, tokenService auth.TokenService) *AuthService {
	return &AuthService{
		userRepo:     userRepo,
		tokenService: tokenService,
	}
}

func (s *AuthService) Register(email, password, name string, genres []string) (*auth.User, error) {
	// Check if user already exists
	existingUser, err := s.userRepo.FindByEmail(email)
	if err != nil {
		return nil, err
	}
	if existingUser != nil {
		return nil, auth.ErrEmailAlreadyExists
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	// Create new user
	user := &auth.User{
		ID:       uuid.New().String(),
		Email:    email,
		Password: string(hashedPassword),
		Name:     name,
		Genres:   genres,
	}

	// Save user to database
	if err := s.userRepo.Create(user); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *AuthService) Login(email, password string) (string, error) {
	user, err := s.userRepo.FindByEmail(email)
	if err != nil {
		return "", err
	}
	if user == nil {
		return "", auth.ErrInvalidCredentials
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return "", auth.ErrInvalidCredentials
	}

	return s.tokenService.GenerateToken(user.ID, user.Email)
}

func (s *AuthService) GetUserByID(id string) (*auth.User, error) {
	return s.userRepo.FindByID(id)
}
