package auth

import (
	"testing"
	"time"
)

func TestNewUser(t *testing.T) {
	tests := []struct {
		name     string
		email    string
		userName string
		password string
		genres   []string
		wantErr  bool
		errType  error
	}{
		{
			name:     "Valid user",
			email:    "test@example.com",
			userName: "John Doe",
			password: "password123",
			genres:   []string{"fiction"},
			wantErr:  false,
		},
		{
			name:     "Invalid email",
			email:    "invalid-email",
			userName: "John Doe",
			password: "password123",
			genres:   []string{"fiction"},
			wantErr:  true,
			errType:  ErrInvalidEmail,
		},
		{
			name:     "Short password",
			email:    "test@example.com",
			userName: "John Doe",
			password: "short",
			genres:   []string{"fiction"},
			wantErr:  true,
			errType:  ErrInvalidPassword,
		},
		{
			name:     "Short name",
			email:    "test@example.com",
			userName: "J",
			password: "password123",
			genres:   []string{"fiction"},
			wantErr:  true,
			errType:  ErrInvalidName,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user, err := NewUser(tt.email, tt.userName, tt.password, tt.genres)

			if tt.wantErr {
				if err == nil {
					t.Error("Expected error but got none")
				}
				if err != tt.errType {
					t.Errorf("Expected error %v but got %v", tt.errType, err)
				}
				return
			}

			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}

			// Check user fields
			if user.Email != tt.email {
				t.Errorf("Expected email %s but got %s", tt.email, user.Email)
			}
			if user.Name != tt.userName {
				t.Errorf("Expected name %s but got %s", tt.userName, user.Name)
			}
			if user.Password != tt.password {
				t.Errorf("Expected password %s but got %s", tt.password, user.Password)
			}

			// Check timestamps
			now := time.Now()
			if user.CreatedAt.After(now) {
				t.Error("CreatedAt should not be in the future")
			}
			if user.UpdatedAt.After(now) {
				t.Error("UpdatedAt should not be in the future")
			}
		})
	}
}

func TestValidateEmail(t *testing.T) {
	tests := []struct {
		name    string
		email   string
		wantErr bool
	}{
		{"Valid email", "test@example.com", false},
		{"Invalid email - no @", "testexample.com", true},
		{"Invalid email - no domain", "test@", true},
		{"Invalid email - spaces", "test @example.com", true},
		{"Invalid email - empty", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateEmail(tt.email)
			if (err != nil) != tt.wantErr {
				t.Errorf("validateEmail() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestUpdatePassword(t *testing.T) {
	user := &User{
		Email:    "test@example.com",
		Name:     "Test User",
		Password: "oldpassword",
	}

	tests := []struct {
		name     string
		password string
		wantErr  bool
	}{
		{"Valid password", "newpassword123", false},
		{"Short password", "short", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			originalUpdatedAt := user.UpdatedAt
			err := user.UpdatePassword(tt.password)

			if (err != nil) != tt.wantErr {
				t.Errorf("UpdatePassword() error = %v, wantErr %v", err, tt.wantErr)
			}

			if !tt.wantErr {
				if user.Password != tt.password {
					t.Errorf("Password not updated, got = %v, want %v", user.Password, tt.password)
				}
				if !user.UpdatedAt.After(originalUpdatedAt) {
					t.Error("UpdatedAt timestamp not updated")
				}
			}
		})
	}
}
