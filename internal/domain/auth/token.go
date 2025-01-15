// New file for token-related domain logic
package auth

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt"
)

type Claims struct {
	UserID    string `json:"user_id"`
	Email     string `json:"email"`
	ExpiresAt int64  `json:"exp"`
}

// Valid implements jwt.Claims interface
func (c Claims) Valid() error {
	if time.Now().Unix() > c.ExpiresAt {
		return errors.New("token has expired")
	}
	return nil
}

type TokenService interface {
	GenerateToken(userID, email string) (string, error)
	ValidateToken(tokenString string) (*Claims, error)
}

type JWTService struct {
	secretKey []byte
	duration  time.Duration
}

func NewJWTService(secretKey string, duration time.Duration) *JWTService {
	return &JWTService{
		secretKey: []byte(secretKey),
		duration:  duration,
	}
}

func (s *JWTService) GenerateToken(userID, email string) (string, error) {
	claims := Claims{
		UserID:    userID,
		Email:     email,
		ExpiresAt: time.Now().Add(s.duration).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(s.secretKey)
}

type UserResponse struct {
	ID    string `json:"id"`
	Email string `json:"email"`
	Name  string `json:"name"`
}

type LoginResponse struct {
	Token string       `json:"token"`
	User  UserResponse `json:"user"`
}
