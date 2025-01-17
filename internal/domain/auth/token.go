// New file for token-related domain logic
package auth

import (
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt"
)

type Claims struct {
	UserID string `json:"user_id"`
	Email  string `json:"email"`
	jwt.StandardClaims
}

// Valid implements jwt.Claims interface
func (c Claims) Valid() error {
	if time.Now().Unix() > c.ExpiresAt {
		return errors.New("token has expired")
	}
	return nil
}

type TokenService interface {
	GenerateToken(userID string, email string) (string, error)
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
		UserID: userID,
		Email:  email,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(s.duration).Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(s.secretKey)
}

func (s *JWTService) ValidateToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return s.secretKey, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
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

func GenerateToken(userID string, email string) (string, error) {
	// Get secret key from environment variable
	secretKey := []byte(os.Getenv("JWT_SECRET"))
	if len(secretKey) == 0 {
		secretKey = []byte("your-256-bit-secret") // Default for development
	}

	// Create claims
	claims := &Claims{
		UserID: userID,
		Email:  email,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(24 * time.Hour).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
	}

	// Create token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign token with secret key
	return token.SignedString(secretKey)
}

func ValidateToken(tokenString string) (*Claims, error) {
	// Get secret key from environment variable
	secretKey := []byte(os.Getenv("JWT_SECRET"))
	if len(secretKey) == 0 {
		secretKey = []byte("your-256-bit-secret") // Default for development
	}

	// Parse token
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return secretKey, nil
	})

	if err != nil {
		return nil, err
	}

	// Validate token and return claims
	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, jwt.ErrSignatureInvalid
}
