package config

import (
	"os"
	"testing"
	"time"
)

// BLOCK_CONFIG_TEST_001
// Purpose: Verifies config loading defaults and environment variable overrides.
func TestConfigLoad(t *testing.T) {
	cfg := Load()
	if cfg.HTTPPort != "8080" {
		t.Errorf("expected default port 8080, got %s", cfg.HTTPPort)
	}
	if cfg.DBSocketPath != "/var/run/campus-os/db.sock" {
		t.Errorf("expected default socket path, got %s", cfg.DBSocketPath)
	}

	os.Setenv("PORT", "9090")
	defer os.Unsetenv("PORT")

	cfgOverridden := Load()
	if cfgOverridden.HTTPPort != "9090" {
		t.Errorf("expected overridden port 9090, got %s", cfgOverridden.HTTPPort)
	}
	if cfgOverridden.ReadTimeout != 10*time.Second {
		t.Errorf("expected default read timeout 10s, got %v", cfgOverridden.ReadTimeout)
	}
}
