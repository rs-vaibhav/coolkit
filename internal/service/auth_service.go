package service

import (
	"context"
	"errors"
	"regexp"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	"github.com/coolkit-org/coolkit/internal/model"
	"github.com/coolkit-org/coolkit/internal/repository"
)

var (
	ErrEmailTaken         = errors.New("email is already taken")
	ErrInvalidCredentials = errors.New("invalid email or password")
	ErrUserNotFound       = errors.New("user not found")
	ErrWeakPassword       = errors.New("password must be at least 8 characters long and contain letters and numbers")
	ErrInvalidEmail       = errors.New("invalid email format")
)

// validateEmail checks if the email format is valid
func validateEmail(email string) bool {
	pattern := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	re := regexp.MustCompile(pattern)
	return re.MatchString(email)
}

// validatePassword checks if the password meets security requirements
func validatePassword(password string) bool {
	if len(password) < 8 {
		return false
	}
	
	hasLetter := false
	hasNumber := false
	for _, c := range password {
		if (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') {
			hasLetter = true
		}
		if c >= '0' && c <= '9' {
			hasNumber = true
		}
	}
	return hasLetter && hasNumber
}

type AuthService struct {
	userRepo  *repository.UserRepository
	jwtSecret string
}

func NewAuthService(userRepo *repository.UserRepository, jwtSecret string) *AuthService {
	return &AuthService{
		userRepo:  userRepo,
		jwtSecret: jwtSecret,
	}
}

func (s *AuthService) Register(ctx context.Context, name, email, password string) (string, *model.User, error) {
	// Validate email format
	if !validateEmail(email) {
		return "", nil, ErrInvalidEmail
	}

	// Validate password strength
	if !validatePassword(password) {
		return "", nil, ErrWeakPassword
	}

	existingUser, _ := s.userRepo.FindByEmail(ctx, email)
	if existingUser != nil {
		return "", nil, ErrEmailTaken
	}

	hash, err := s.HashPassword(password)
	if err != nil {
		return "", nil, err
	}

	user := &model.User{
		ID:       uuid.New(),
		Name:     name,
		Email:    email,
		PasswordHash: hash,
	}

	err = s.userRepo.Create(ctx, user)
	if err != nil {
		return "", nil, err
	}

	token, err := s.generateToken(user)
	if err != nil {
		return "", nil, err
	}

	return token, user, nil
}

func (s *AuthService) Login(ctx context.Context, email, password string) (string, *model.User, error) {
	user, err := s.userRepo.FindByEmail(ctx, email)
	if err != nil {
		return "", nil, ErrInvalidCredentials
	}

	if !s.CheckPassword(password, user.PasswordHash) {
		return "", nil, ErrInvalidCredentials
	}

	token, err := s.generateToken(user)
	if err != nil {
		return "", nil, err
	}

	return token, user, nil
}

func (s *AuthService) GetCurrentUser(ctx context.Context, userID uuid.UUID) (*model.User, error) {
	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return nil, ErrUserNotFound
	}
	return user, nil
}

func (s *AuthService) generateToken(user *model.User) (string, error) {
	claims := jwt.MapClaims{
		"user_id": user.ID.String(),
		"email":   user.Email,
		"exp":     time.Now().Add(24 * time.Hour).Unix(),
		"iat":     time.Now().Unix(),
		"iss":     "coolkit",
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.jwtSecret))
}

func (s *AuthService) HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	return string(bytes), err
}

func (s *AuthService) CheckPassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
