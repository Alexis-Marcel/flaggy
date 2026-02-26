package models

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"time"
)

type Environment string

const (
	EnvLive    Environment = "live"
	EnvTest    Environment = "test"
	EnvStaging Environment = "staging"
)

func ValidateEnvironment(env Environment) error {
	switch env {
	case EnvLive, EnvTest, EnvStaging:
		return nil
	default:
		return fmt.Errorf("invalid environment: %q (must be live, test, or staging)", env)
	}
}

type APIKey struct {
	ID          string      `json:"id"`
	Name        string      `json:"name"`
	Environment Environment `json:"environment"`
	Prefix      string      `json:"prefix"`
	Revoked     bool        `json:"revoked"`
	CreatedAt   time.Time   `json:"created_at"`
	LastUsedAt  *time.Time  `json:"last_used_at,omitempty"`
}

// APIKeyWithRaw is returned only at creation time. RawKey is shown once.
type APIKeyWithRaw struct {
	APIKey
	RawKey string `json:"key"`
}

// GenerateAPIKey creates a new API key with a random suffix.
// The raw key is returned for the one-time display.
// Only the SHA-256 hash is stored.
func GenerateAPIKey(name string, env Environment) (*APIKeyWithRaw, string) {
	suffix := make([]byte, 32)
	if _, err := rand.Read(suffix); err != nil {
		panic(fmt.Sprintf("crypto/rand failed: %v", err))
	}
	suffixHex := hex.EncodeToString(suffix)

	rawKey := fmt.Sprintf("flg_%s_%s", env, suffixHex)
	prefix := fmt.Sprintf("flg_%s_%s", env, suffixHex[:4])

	hash := sha256.Sum256([]byte(rawKey))
	hashedKey := hex.EncodeToString(hash[:])

	id := fmt.Sprintf("key_%s", suffixHex[:16])

	key := &APIKeyWithRaw{
		APIKey: APIKey{
			ID:          id,
			Name:        name,
			Environment: env,
			Prefix:      prefix,
			CreatedAt:   time.Now().UTC(),
		},
		RawKey: rawKey,
	}

	return key, hashedKey
}

// HashKey returns the SHA-256 hex digest of a raw API key.
func HashKey(rawKey string) string {
	h := sha256.Sum256([]byte(rawKey))
	return hex.EncodeToString(h[:])
}
