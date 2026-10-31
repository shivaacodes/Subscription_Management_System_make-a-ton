package auth

import (
	"errors"
	"testing"
)

func TestValidatePassword(t *testing.T) {
	tests := []struct {
		name     string
		password string
		wantErr  error
	}{
		{
			name:     "valid password with >8 chars, upper, lower, special",
			password: "SecurePass123!",
			wantErr:  nil,
		},
		{
			name:     "exactly 8 characters should fail (must be strictly > 8)",
			password: "Pass123!",
			wantErr:  ErrPasswordTooShort,
		},
		{
			name:     "short password 6 characters",
			password: "P123!",
			wantErr:  ErrPasswordTooShort,
		},
		{
			name:     "missing uppercase",
			password: "securepass123!",
			wantErr:  ErrPasswordNoUpper,
		},
		{
			name:     "missing lowercase",
			password: "SECUREPASS123!",
			wantErr:  ErrPasswordNoLower,
		},
		{
			name:     "missing special character",
			password: "SecurePass1234",
			wantErr:  ErrPasswordNoSpecial,
		},
		{
			name:     "valid complex password",
			password: "Brightloop#2026-Operate",
			wantErr:  nil,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := ValidatePassword(tc.password)
			if !errors.Is(err, tc.wantErr) {
				t.Errorf("ValidatePassword(%q) error = %v, wantErr = %v", tc.password, err, tc.wantErr)
			}
		})
	}
}

func TestPasswordHashingAndVerification(t *testing.T) {
	raw := "SuperSecret123!"
	hash, err := HashPassword(raw)
	if err != nil {
		t.Fatalf("HashPassword failed: %v", err)
	}

	if !CheckPasswordHash(raw, hash) {
		t.Errorf("CheckPasswordHash should return true for matching password")
	}

	if CheckPasswordHash("WrongSecret123!", hash) {
		t.Errorf("CheckPasswordHash should return false for incorrect password")
	}
}

func TestTokenGeneration(t *testing.T) {
	rawToken1, hash1, err := GenerateSecureToken()
	if err != nil {
		t.Fatalf("GenerateSecureToken failed: %v", err)
	}

	rawToken2, hash2, err := GenerateSecureToken()
	if err != nil {
		t.Fatalf("GenerateSecureToken failed: %v", err)
	}

	if rawToken1 == rawToken2 {
		t.Errorf("Generated tokens must be distinct")
	}

	if HashToken(rawToken1) != hash1 {
		t.Errorf("HashToken mismatch")
	}

	if hash1 == hash2 {
		t.Errorf("Token hashes must be distinct")
	}
}
