package users

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	ErrUserNotFound    = errors.New("user not found")
	ErrDuplicateEmail  = errors.New("email address is already in use")
	ErrProfileNotFound = errors.New("customer profile not found")
)

type Repository interface {
	CreateUserWithProfile(ctx context.Context, email, passwordHash string, role Role, companyName *string) (*User, *CustomerProfile, error)
	GetByEmail(ctx context.Context, email string) (*User, error)
	GetByID(ctx context.Context, id string) (*User, error)
	GetProfileByUserID(ctx context.Context, userID string) (*CustomerProfile, error)
}

type pgRepository struct {
	pool *pgxpool.Pool
}

func NewRepository(pool *pgxpool.Pool) Repository {
	return &pgRepository{pool: pool}
}

func (r *pgRepository) CreateUserWithProfile(
	ctx context.Context,
	email, passwordHash string,
	role Role,
	companyName *string,
) (*User, *CustomerProfile, error) {
	normalizedEmail := strings.ToLower(strings.TrimSpace(email))

	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	var u User
	userQuery := `
		INSERT INTO users (email, password_hash, role, created_at, updated_at)
		VALUES ($1, $2, $3, NOW(), NOW())
		RETURNING id, email, password_hash, role, created_at, updated_at
	`
	err = tx.QueryRow(ctx, userQuery, normalizedEmail, passwordHash, role).Scan(
		&u.ID,
		&u.Email,
		&u.PasswordHash,
		&u.Role,
		&u.CreatedAt,
		&u.UpdatedAt,
	)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return nil, nil, ErrDuplicateEmail
		}
		return nil, nil, fmt.Errorf("failed to insert user: %w", err)
	}

	var cp CustomerProfile
	profileQuery := `
		INSERT INTO customer_profiles (user_id, company_name, created_at, updated_at)
		VALUES ($1, $2, NOW(), NOW())
		RETURNING id, user_id, company_name, phone, billing_address, created_at, updated_at
	`
	err = tx.QueryRow(ctx, profileQuery, u.ID, companyName).Scan(
		&cp.ID,
		&cp.UserID,
		&cp.CompanyName,
		&cp.Phone,
		&cp.BillingAddress,
		&cp.CreatedAt,
		&cp.UpdatedAt,
	)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to insert customer profile: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return &u, &cp, nil
}

func (r *pgRepository) GetByEmail(ctx context.Context, email string) (*User, error) {
	normalizedEmail := strings.ToLower(strings.TrimSpace(email))
	query := `
		SELECT id, email, password_hash, role, created_at, updated_at
		FROM users
		WHERE LOWER(email) = LOWER($1)
	`
	var u User
	err := r.pool.QueryRow(ctx, query, normalizedEmail).Scan(
		&u.ID,
		&u.Email,
		&u.PasswordHash,
		&u.Role,
		&u.CreatedAt,
		&u.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrUserNotFound
		}
		return nil, fmt.Errorf("error querying user by email: %w", err)
	}
	return &u, nil
}

func (r *pgRepository) GetByID(ctx context.Context, id string) (*User, error) {
	query := `
		SELECT id, email, password_hash, role, created_at, updated_at
		FROM users
		WHERE id = $1
	`
	var u User
	err := r.pool.QueryRow(ctx, query, id).Scan(
		&u.ID,
		&u.Email,
		&u.PasswordHash,
		&u.Role,
		&u.CreatedAt,
		&u.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrUserNotFound
		}
		return nil, fmt.Errorf("error querying user by id: %w", err)
	}
	return &u, nil
}

func (r *pgRepository) GetProfileByUserID(ctx context.Context, userID string) (*CustomerProfile, error) {
	query := `
		SELECT id, user_id, company_name, phone, billing_address, created_at, updated_at
		FROM customer_profiles
		WHERE user_id = $1
	`
	var cp CustomerProfile
	err := r.pool.QueryRow(ctx, query, userID).Scan(
		&cp.ID,
		&cp.UserID,
		&cp.CompanyName,
		&cp.Phone,
		&cp.BillingAddress,
		&cp.CreatedAt,
		&cp.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrProfileNotFound
		}
		return nil, fmt.Errorf("error querying profile: %w", err)
	}
	return &cp, nil
}
