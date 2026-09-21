package auth

import (
	"os"
	"path/filepath"
	"testing"
)

// BLOCK_AUTH_SESSION_TEST_001
// Purpose: Verifies AuthSession persistence, encoding, and clearing semantics.
func TestKeyringSessionStore_FallbackFile(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "campus-session-test-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	store := NewKeyringSessionStore(tempDir)

	// Initially empty
	sess, err := store.loadFallback()
	if err == nil || sess != nil {
		t.Errorf("expected error on empty store, got %v", sess)
	}

	// Save test session
	testSess := &AuthSession{
		AccessToken:  "jwt_access_123",
		RefreshToken: "jwt_refresh_456",
		UserID:       "u-chairperson-1",
		Username:     "chairperson.2024",
		FullName:     "Yogesh Kumar Mallik",
		Email:        "chairperson@campus.edu",
		RoleCode:     "SUPER_ADMIN",
		RoleCodes:    []string{"SUPER_ADMIN"},
	}

	if err := store.Save(testSess); err != nil {
		t.Fatalf("failed to save session: %v", err)
	}

	// Load session back
	loaded, err := store.Load()
	if err != nil {
		t.Fatalf("failed to load session: %v", err)
	}
	if loaded.Username != "chairperson.2024" || loaded.RoleCode != "SUPER_ADMIN" {
		t.Errorf("loaded session mismatch: %+v", loaded)
	}

	// Verify file permissions (0600)
	filePath := filepath.Join(tempDir, "session.json")
	if fi, err := os.Stat(filePath); err == nil {
		perm := fi.Mode().Perm()
		if perm != 0600 {
			t.Errorf("expected 0600 file permissions, got %o", perm)
		}
	}

	// Clear session
	if err := store.Clear(); err != nil {
		t.Fatalf("failed to clear session: %v", err)
	}

	// Confirm cleared
	afterClear, _ := store.loadFallback()
	if afterClear != nil {
		t.Errorf("expected nil session after clear, got %+v", afterClear)
	}
}

// BLOCK_AUTH_SESSION_TEST_002
// Purpose: Verifies error when saving a nil session.
func TestKeyringSessionStore_NilSave(t *testing.T) {
	store := NewKeyringSessionStore("")
	if err := store.Save(nil); err == nil {
		t.Errorf("expected error when saving nil session")
	}
}
