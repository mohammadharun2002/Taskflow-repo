package auth

import (
	"context"
	"errors"
	"fmt"
	"net/mail"
	"strconv"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type Service struct {
	repo       Repository
	jwtSecrect []byte
	tokenTTL   time.Duration
}

func NewService(repo Repository, jwtSecret string) *Service {
	return &Service{
		repo:       repo,
		jwtSecrect: []byte(jwtSecret),
		tokenTTL:   24 * time.Hour,
	}
}

func (s *Service) Register(
	ctx context.Context,
	request RegisterRequest,
) (User, error) {
	name := strings.TrimSpace(request.Name)
	if name == "" {
		return User{}, ErrNameRequired
	}

	email := strings.ToLower(strings.TrimSpace(request.Email))
	if email == "" {
		return User{}, ErrEmailRequired
	}

	parsedEmail, err := mail.ParseAddress(email)
	if err != nil || parsedEmail.Address != email {
		return User{}, ErrInvalidEmail
	}

	if len(request.Password) < 8 {
		return User{}, ErrPasswordTooShort
	}

	passwordHash, err := bcrypt.GenerateFromPassword(
		[]byte(request.Password),
		bcrypt.DefaultCost,
	)
	if err != nil {
		return User{}, err
	}

	now := time.Now()

	user := User{
		Name:         name,
		Email:        email,
		PasswordHash: string(passwordHash),
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	return s.repo.Create(ctx, user)
}

func (s *Service) Login(ctx context.Context, request LoginRequest) (AuthResponse, error) {
	email := strings.ToLower(strings.TrimSpace(request.Email))
	user, err := s.repo.FindByEmail(ctx, email)
	if errors.Is(err, ErrUserNotFound) {
		return AuthResponse{}, ErrInvalidCredentials
	}
	if err != nil {
		return AuthResponse{}, err
	}

	err = bcrypt.CompareHashAndPassword(
		[]byte(user.PasswordHash),
		[]byte(request.Password),
	)
	if err != nil {
		return AuthResponse{}, ErrInvalidCredentials
	}

	token, err := s.generateToken(user)
	if err != nil {
		return AuthResponse{}, err
	}
	return AuthResponse{
		Token: token,
		User:  user,
	}, nil
}

func (s *Service) generateToken(user User) (string, error) {
	now := time.Now()

	claims := jwt.RegisteredClaims{
		Subject:   strconv.FormatInt(user.ID, 10),
		Issuer:    "taskflow",
		IssuedAt:  jwt.NewNumericDate(now),
		ExpiresAt: jwt.NewNumericDate(now.Add(s.tokenTTL)),
	}

	token := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		claims,
	)

	signedToken, err := token.SignedString(s.jwtSecrect)
	if err != nil {
		return "", fmt.Errorf("sign JWT: %w", err)
	}

	return signedToken, nil
}
