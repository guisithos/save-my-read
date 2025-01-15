package application

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/guisithos/save-my-read/internal/domain/user"
)

type AuthService struct {
	userRepo user.Repository
	jwtKey   []byte
}

func NewAuthService(userRepo user.Repository, jwtKey string) *AuthService {
	return &AuthService{
		userRepo: userRepo,
		jwtKey:   []byte(jwtKey),
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

func (s *AuthService) Login(email, password string) (string, error) {
	user, err := s.userRepo.FindByEmail(email)
	if err != nil {
		return "", errors.New("invalid credentials")
	}

	if !user.ValidatePassword(password) {
		return "", errors.New("invalid credentials")
	}

	// Generate JWT token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"exp":     time.Now().Add(time.Hour * 24).Unix(),
	})

	tokenString, err := token.SignedString(s.jwtKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
