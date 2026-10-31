package auth

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var ErrSessionNotFound = errors.New("session not found or expired")

type Session struct {
	ID        string    `json:"id"`
	TokenHash string    `json:"-"`
	UserID    string    `json:"userId"`
	UserAgent *string   `json:"userAgent,omitempty"`
	IPAddress *string   `json:"ipAddress,omitempty"`
	ExpiresAt time.Time `json:"expiresAt"`
	CreatedAt time.Time `json:"createdAt"`
}

type SessionRepository interface {
	CreateSession(ctx context.Context, userID string, tokenHash string, userAgent, ipAddress *string, expiresAt time.Time) (*Session, error)
	GetByTokenHash(ctx context.Context, tokenHash string) (*Session, error)
	DeleteByTokenHash(ctx context.Context, tokenHash string) error
	DeleteExpired(ctx context.Context) error
}

type pgSessionRepo struct {
	pool *pgxpool.Pool
}

func NewSessionRepository(pool *pgxpool.Pool) SessionRepository {
	return &pgSessionRepo{pool: pool}
}

func GenerateSecureToken() (rawToken string, tokenHash string, err error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", "", fmt.Errorf("failed generating random token: %w", err)
	}
	rawToken = hex.EncodeToString(b)
	tokenHash = HashToken(rawToken)
	return rawToken, tokenHash, nil
}

func HashToken(rawToken string) string {
	sum := sha256.Sum256([]byte(rawToken))
	return hex.EncodeToString(sum[:])
}

func (r *pgSessionRepo) CreateSession(
	ctx context.Context,
	userID string,
	tokenHash string,
	userAgent, ipAddress *string,
	expiresAt time.Time,
) (*Session, error) {
	query := `
		INSERT INTO sessions (token_hash, user_id, user_agent, ip_address, expires_at, created_at)
		VALUES ($1, $2, $3, $4, $5, NOW())
		RETURNING id, token_hash, user_id, user_agent, ip_address, expires_at, created_at
	`
	var s Session
	err := r.pool.QueryRow(ctx, query, tokenHash, userID, userAgent, ipAddress, expiresAt).Scan(
		&s.ID,
		&s.TokenHash,
		&s.UserID,
		&s.UserAgent,
		&s.IPAddress,
		&s.ExpiresAt,
		&s.CreatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create session: %w", err)
	}
	return &s, nil
}

func (r *pgSessionRepo) GetByTokenHash(ctx context.Context, tokenHash string) (*Session, error) {
	query := `
		SELECT id, token_hash, user_id, user_agent, ip_address, expires_at, created_at
		FROM sessions
		WHERE token_hash = $1 AND expires_at > NOW()
	`
	var s Session
	err := r.pool.QueryRow(ctx, query, tokenHash).Scan(
		&s.ID,
		&s.TokenHash,
		&s.UserID,
		&s.UserAgent,
		&s.IPAddress,
		&s.ExpiresAt,
		&s.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrSessionNotFound
		}
		return nil, fmt.Errorf("failed querying session: %w", err)
	}
	return &s, nil
}

func (r *pgSessionRepo) DeleteByTokenHash(ctx context.Context, tokenHash string) error {
	query := `DELETE FROM sessions WHERE token_hash = $1`
	_, err := r.pool.Exec(ctx, query, tokenHash)
	return err
}

func (r *pgSessionRepo) DeleteExpired(ctx context.Context) error {
	query := `DELETE FROM sessions WHERE expires_at <= NOW()`
	_, err := r.pool.Exec(ctx, query)
	return err
}
