// New file for token-related domain logic
package auth

import (
	"time"

	"github.com/golang-jwt/jwt"
)

type Claims struct {
	UserID    string `json:"user_id"`
	Email     string `json:"email"`
	ExpiresAt int64  `json:"exp"`
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
