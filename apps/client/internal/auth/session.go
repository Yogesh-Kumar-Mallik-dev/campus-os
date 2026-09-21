package auth

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/zalando/go-keyring"
)

const (
	// KeyringService is the service name registered with the OS Secret Service / Keyring.
	KeyringService = "internal.campus-os.client"
	// KeyringUser is the key account identifier for the active desktop session.
	KeyringUser = "active_session"
)

// AuthSession represents an active authenticated user session.
type AuthSession struct {
	AccessToken  string   `json:"access_token"`
	RefreshToken string   `json:"refresh_token"`
	UserID       string   `json:"user_id"`
	Username     string   `json:"username"`
	FullName     string   `json:"full_name"`
	Email        string   `json:"email"`
	RoleCode     string   `json:"role_code"`
	RoleCodes    []string `json:"role_codes,omitempty"`
}

// SessionStore defines the contract for persisting, reading, and wiping sessions.
type SessionStore interface {
	Save(session *AuthSession) error
	Load() (*AuthSession, error)
	Clear() error
}

// KeyringSessionStore implements SessionStore using native OS Keyring with protected file fallback.
type KeyringSessionStore struct {
	service     string
	user        string
	fallbackDir string
}

// BLOCK_AUTH_KEYRING_STORE_001
// Purpose: Constructs a KeyringSessionStore with automatic path resolution.
func NewKeyringSessionStore(fallbackDir string) *KeyringSessionStore {
	if fallbackDir == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			home = os.TempDir()
		}
		fallbackDir = filepath.Join(home, ".config", "campus-os")
	}
	return &KeyringSessionStore{
		service:     KeyringService,
		user:        KeyringUser,
		fallbackDir: fallbackDir,
	}
}

// Save stores the session into OS Keyring or user-restricted fallback file (0600).
func (s *KeyringSessionStore) Save(sess *AuthSession) error {
	if sess == nil {
		return errors.New("cannot save nil session")
	}
	data, err := json.Marshal(sess)
	if err != nil {
		return fmt.Errorf("failed to encode session: %w", err)
	}

	// Attempt native OS Keyring first
	err = keyring.Set(s.service, s.user, string(data))
	if err == nil {
		return nil
	}

	// Graceful fallback to user-restricted local file (0600)
	return s.saveFallback(data)
}

// Load retrieves the session from OS Keyring or fallback file.
func (s *KeyringSessionStore) Load() (*AuthSession, error) {
	val, err := keyring.Get(s.service, s.user)
	if err == nil && val != "" {
		var sess AuthSession
		if err := json.Unmarshal([]byte(val), &sess); err == nil {
			return &sess, nil
		}
	}

	return s.loadFallback()
}

// Clear wipes credentials from OS Keyring and fallback file.
func (s *KeyringSessionStore) Clear() error {
	_ = keyring.Delete(s.service, s.user)
	_ = os.Remove(filepath.Join(s.fallbackDir, "session.json"))
	return nil
}

func (s *KeyringSessionStore) saveFallback(data []byte) error {
	if err := os.MkdirAll(s.fallbackDir, 0700); err != nil {
		return err
	}
	path := filepath.Join(s.fallbackDir, "session.json")
	return os.WriteFile(path, data, 0600)
}

func (s *KeyringSessionStore) loadFallback() (*AuthSession, error) {
	path := filepath.Join(s.fallbackDir, "session.json")
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var sess AuthSession
	if err := json.Unmarshal(data, &sess); err != nil {
		return nil, err
	}
	return &sess, nil
}
