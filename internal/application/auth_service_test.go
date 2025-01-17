package application

import (
	"testing"

	"github.com/guisithos/save-my-read/internal/domain/auth"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockUserRepository is a mock implementation of auth.UserRepository
type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) Create(user *auth.User) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *MockUserRepository) FindByEmail(email string) (*auth.User, error) {
	args := m.Called(email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*auth.User), args.Error(1)
}

func (m *MockUserRepository) FindByID(id string) (*auth.User, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*auth.User), args.Error(1)
}

func (m *MockUserRepository) Update(user *auth.User) error {
	args := m.Called(user)
	return args.Error(0)
}

type MockTokenService struct{}

func (m *MockTokenService) GenerateToken(userID string, email string) (string, error) {
	return "test-token", nil
}

func (m *MockTokenService) ValidateToken(tokenString string) (*auth.Claims, error) {
	return &auth.Claims{UserID: "test-id", Email: "test@example.com"}, nil
}

func TestAuthService_Register(t *testing.T) {
	mockRepo := new(MockUserRepository)
	mockTokenService := new(MockTokenService)
	service := NewAuthService(mockRepo, mockTokenService)

	tests := []struct {
		name     string
		email    string
		password string
		userName string
		genres   []string
		mockFn   func()
		wantErr  bool
		errType  error
	}{
		{
			name:     "Successful registration",
			email:    "test@example.com",
			password: "password123",
			userName: "Test User",
			genres:   []string{"fiction"},
			mockFn: func() {
				mockRepo.On("FindByEmail", "test@example.com").Return(nil, nil)
				mockRepo.On("Create", mock.AnythingOfType("*auth.User")).Return(nil)
			},
			wantErr: false,
		},
		{
			name:     "Email already exists",
			email:    "existing@example.com",
			password: "password123",
			userName: "Test User",
			genres:   []string{"fiction"},
			mockFn: func() {
				mockRepo.On("FindByEmail", "existing@example.com").Return(&auth.User{}, nil)
			},
			wantErr: true,
			errType: auth.ErrEmailAlreadyExists,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo.ExpectedCalls = nil
			mockRepo.Calls = nil
			tt.mockFn()

			user, err := service.Register(tt.email, tt.password, tt.userName, tt.genres)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Equal(t, tt.errType, err)
				assert.Nil(t, user)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, user)
				assert.Equal(t, tt.email, user.Email)
				assert.Equal(t, tt.userName, user.Name)
				assert.Equal(t, tt.genres, user.Genres)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestAuthService_Login(t *testing.T) {
	mockRepo := new(MockUserRepository)
	mockTokenService := new(MockTokenService)
	service := NewAuthService(mockRepo, mockTokenService)

	hashedPassword := "$2a$10$abcdefghijklmnopqrstuvwxyz" // Example hashed password

	tests := []struct {
		name     string
		email    string
		password string
		mockFn   func()
		wantErr  bool
		errType  error
	}{
		{
			name:     "Successful login",
			email:    "test@example.com",
			password: "password123",
			mockFn: func() {
				mockRepo.On("FindByEmail", "test@example.com").Return(&auth.User{
					Email:    "test@example.com",
					Password: hashedPassword,
				}, nil)
			},
			wantErr: false,
		},
		{
			name:     "User not found",
			email:    "nonexistent@example.com",
			password: "password123",
			mockFn: func() {
				mockRepo.On("FindByEmail", "nonexistent@example.com").Return(nil, nil)
			},
			wantErr: true,
			errType: auth.ErrInvalidCredentials,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo.ExpectedCalls = nil
			mockRepo.Calls = nil
			tt.mockFn()

			token, err := service.Login(tt.email, tt.password)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Equal(t, tt.errType, err)
				assert.Empty(t, token)
			} else {
				assert.NoError(t, err)
				assert.NotEmpty(t, token)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}
