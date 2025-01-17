package handlers

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/guisithos/save-my-read/internal/application"
	"github.com/guisithos/save-my-read/internal/domain/auth"
)

var errUserNotFound = errors.New("user not found")

type mockUserRepo struct {
	users map[string]*auth.User
}

func newMockUserRepo() *mockUserRepo {
	return &mockUserRepo{users: make(map[string]*auth.User)}
}

func (m *mockUserRepo) Save(user *auth.User) error {
	m.users[user.Email] = user
	return nil
}

func (m *mockUserRepo) FindByEmail(email string) (*auth.User, error) {
	if user, exists := m.users[email]; exists {
		return user, nil
	}
	return nil, errUserNotFound
}

func (m *mockUserRepo) FindByID(id string) (*auth.User, error) {
	for _, u := range m.users {
		if u.ID == id {
			return u, nil
		}
	}
	return nil, errUserNotFound
}

func (m *mockUserRepo) Update(user *auth.User) error {
	if _, exists := m.users[user.Email]; !exists {
		return errUserNotFound
	}
	m.users[user.Email] = user
	return nil
}

func (m *mockUserRepo) Create(user *auth.User) error {
	m.users[user.ID] = user
	return nil
}

type mockTokenService struct{}

func (m *mockTokenService) GenerateToken(userID string, email string) (string, error) {
	return "test-token", nil
}

func (m *mockTokenService) ValidateToken(token string) (*auth.Claims, error) {
	return &auth.Claims{UserID: "test-id", Email: "test@example.com"}, nil
}

func TestAuthHandler_Register(t *testing.T) {
	repo := newMockUserRepo()
	tokenService := &mockTokenService{}
	authService := application.NewAuthService(repo, tokenService)
	handler := NewAuthHandler(authService)

	tests := []struct {
		name           string
		payload        map[string]interface{}
		expectedStatus int
	}{
		{
			name: "successful registration",
			payload: map[string]interface{}{
				"email":    "test@example.com",
				"password": "password123",
				"name":     "Test User",
				"genres":   []string{"fiction", "mystery"},
			},
			expectedStatus: http.StatusCreated,
		},
		{
			name: "duplicate email",
			payload: map[string]interface{}{
				"email":    "test@example.com",
				"password": "password123",
				"name":     "Test User",
				"genres":   []string{"fiction"},
			},
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			payloadBytes, _ := json.Marshal(tt.payload)
			req := httptest.NewRequest(http.MethodPost, "/register", bytes.NewReader(payloadBytes))
			w := httptest.NewRecorder()

			handler.Register(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, w.Code)
			}
		})
	}
}

func TestAuthHandler_Login(t *testing.T) {
	repo := newMockUserRepo()
	tokenService := &mockTokenService{}
	authService := application.NewAuthService(repo, tokenService)
	handler := NewAuthHandler(authService)

	// Create a test user first
	testUser, _ := auth.NewUser("test@example.com", "password123", "Test User", []string{"fiction"})
	repo.Save(testUser)

	tests := []struct {
		name           string
		payload        map[string]interface{}
		expectedStatus int
	}{
		{
			name: "successful login",
			payload: map[string]interface{}{
				"email":    "test@example.com",
				"password": "password123",
			},
			expectedStatus: http.StatusOK,
		},
		{
			name: "invalid credentials",
			payload: map[string]interface{}{
				"email":    "test@example.com",
				"password": "wrongpassword",
			},
			expectedStatus: http.StatusUnauthorized,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			payloadBytes, _ := json.Marshal(tt.payload)
			req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewReader(payloadBytes))
			w := httptest.NewRecorder()

			handler.Login(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, w.Code)
			}
		})
	}
}
