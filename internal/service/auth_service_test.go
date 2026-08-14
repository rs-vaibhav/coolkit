package service

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// Create a zero-value AuthService to test password helpers.
// These methods don't require a database connection.
func newTestAuthService() *AuthService {
	return &AuthService{
		jwtSecret: "test-secret",
	}
}

func TestHashPassword(t *testing.T) {
	s := newTestAuthService()
	password := "password123"

	hash, err := s.HashPassword(password)

	assert.NoError(t, err)
	assert.NotEmpty(t, hash)
	assert.NotEqual(t, password, hash)
}

func TestCheckPassword(t *testing.T) {
	s := newTestAuthService()
	password := "password123"
	hash, err := s.HashPassword(password)
	assert.NoError(t, err)

	tests := []struct {
		name     string
		password string
		hash     string
		expected bool
	}{
		{
			name:     "correct password",
			password: password,
			hash:     hash,
			expected: true,
		},
		{
			name:     "wrong password",
			password: "wrongpassword",
			hash:     hash,
			expected: false,
		},
		{
			name:     "empty password",
			password: "",
			hash:     hash,
			expected: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := s.CheckPassword(tc.password, tc.hash)
			assert.Equal(t, tc.expected, result)
		})
	}
}

func TestCheckPasswordEmpty(t *testing.T) {
	s := newTestAuthService()
	result := s.CheckPassword("password123", "")
	assert.False(t, result)
}

func TestAuthServiceErrors(t *testing.T) {
	assert.NotNil(t, ErrEmailTaken)
	assert.NotNil(t, ErrInvalidCredentials)
	assert.NotNil(t, ErrUserNotFound)

	assert.Equal(t, "email is already taken", ErrEmailTaken.Error())
	assert.Equal(t, "invalid email or password", ErrInvalidCredentials.Error())
	assert.Equal(t, "user not found", ErrUserNotFound.Error())
}
