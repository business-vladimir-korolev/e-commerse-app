package auth

import (
	"context"
	"errors"
	"net/mail"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

type Service struct {
	repo   Repository
	tokens *TokenManager
}

func NewService(repo Repository) *Service {
	return &Service{
		repo: repo,
	}
}

func toLowerEmail(email string) (string, error) {
	email = strings.ToLower(strings.TrimSpace(email))

	parsed, err := mail.ParseAddress(email)
	if err != nil || parsed.Address != email {
		return "", ErrInvalidEmail
	}

	return email, nil
}

func (s *Service) Register(ctx context.Context, email, password string) (*User, error) {
	email, err := toLowerEmail(email)
	if err != nil {
		return nil, err
	}

	if len(password) < 8 {
		return nil, ErrWeakPassword
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := &User{
		Email:        email,
		PasswordHash: string(hash),
	}

	err = s.repo.Create(ctx, user)
	if errors.Is(err, ErrEmailAlreadyExists) {
		return nil, ErrEmailAlreadyExists
	}

	return user, nil
}

func (s *Service) Login(ctx context.Context, email, password string) (string, error) {
	email, err := toLowerEmail(email)
	if err != nil {
		return "", ErrInvalidCredentials
	}

	user, err := s.repo.GetByEmail(ctx, email)
	if err != nil {
		return "", ErrInvalidCredentials
	}

	err = bcrypt.CompareHashAndPassword(
		[]byte(user.PasswordHash),
		[]byte(password),
	)
	if err != nil {
		return "", ErrInvalidCredentials
	}

	token, err := s.tokens.Generate(user.ID)
	if err != nil {
		return "", err
	}

	return token, nil
}
