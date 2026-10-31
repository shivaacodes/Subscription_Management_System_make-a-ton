package auth

import (
	"context"
	"errors"
	"fmt"
	"net/mail"
	"strings"
	"time"

	"github.com/subflow/backend/internal/users"
)

var (
	ErrInvalidEmail = errors.New("invalid email address format")
)

type Service interface {
	RegisterCustomer(ctx context.Context, email, password string, companyName *string, userAgent, ipAddress *string) (*users.User, string, error)
	Login(ctx context.Context, email, password string, userAgent, ipAddress *string) (*users.User, string, error)
	ValidateSessionToken(ctx context.Context, rawToken string) (*users.User, error)
	Logout(ctx context.Context, rawToken string) error
}

type authService struct {
	userRepo    users.Repository
	sessionRepo SessionRepository
	sessionTTL  time.Duration
}

func NewService(userRepo users.Repository, sessionRepo SessionRepository) Service {
	return &authService{
		userRepo:    userRepo,
		sessionRepo: sessionRepo,
		sessionTTL:  7 * 24 * time.Hour, // 7 days session
	}
}

func (s *authService) RegisterCustomer(
	ctx context.Context,
	email, password string,
	companyName *string,
	userAgent, ipAddress *string,
) (*users.User, string, error) {
	email = strings.TrimSpace(email)
	if _, err := mail.ParseAddress(email); err != nil {
		return nil, "", ErrInvalidEmail
	}

	if err := ValidatePassword(password); err != nil {
		return nil, "", err
	}

	hash, err := HashPassword(password)
	if err != nil {
		return nil, "", fmt.Errorf("hashing error: %w", err)
	}

	user, _, err := s.userRepo.CreateUserWithProfile(ctx, email, hash, users.RoleCustomer, companyName)
	if err != nil {
		return nil, "", err
	}

	rawToken, tokenHash, err := GenerateSecureToken()
	if err != nil {
		return nil, "", err
	}

	expiresAt := time.Now().UTC().Add(s.sessionTTL)
	_, err = s.sessionRepo.CreateSession(ctx, user.ID, tokenHash, userAgent, ipAddress, expiresAt)
	if err != nil {
		return nil, "", fmt.Errorf("failed to create initial session: %w", err)
	}

	return user, rawToken, nil
}

func (s *authService) Login(
	ctx context.Context,
	email, password string,
	userAgent, ipAddress *string,
) (*users.User, string, error) {
	email = strings.TrimSpace(email)
	user, err := s.userRepo.GetByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, users.ErrUserNotFound) {
			return nil, "", ErrInvalidCredentials
		}
		return nil, "", err
	}

	if !CheckPasswordHash(password, user.PasswordHash) {
		return nil, "", ErrInvalidCredentials
	}

	rawToken, tokenHash, err := GenerateSecureToken()
	if err != nil {
		return nil, "", err
	}

	expiresAt := time.Now().UTC().Add(s.sessionTTL)
	_, err = s.sessionRepo.CreateSession(ctx, user.ID, tokenHash, userAgent, ipAddress, expiresAt)
	if err != nil {
		return nil, "", fmt.Errorf("failed to store session: %w", err)
	}

	return user, rawToken, nil
}

func (s *authService) ValidateSessionToken(ctx context.Context, rawToken string) (*users.User, error) {
	if rawToken == "" {
		return nil, ErrSessionNotFound
	}
	tokenHash := HashToken(rawToken)
	session, err := s.sessionRepo.GetByTokenHash(ctx, tokenHash)
	if err != nil {
		return nil, err
	}

	user, err := s.userRepo.GetByID(ctx, session.UserID)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *authService) Logout(ctx context.Context, rawToken string) error {
	if rawToken == "" {
		return nil
	}
	tokenHash := HashToken(rawToken)
	return s.sessionRepo.DeleteByTokenHash(ctx, tokenHash)
}
