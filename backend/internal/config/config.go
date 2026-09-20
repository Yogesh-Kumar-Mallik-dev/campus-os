package config

import (
	"os"
	"strconv"
	"time"
)

// Config represents runtime configuration parameters for the Go Logical Backend.
type Config struct {
	HTTPPort       string
	DBSocketPath   string
	Environment    string
	ReadTimeout    time.Duration
	WriteTimeout   time.Duration
	DefaultTimeout time.Duration
}

// BLOCK_CONFIG_LOAD_001
// Purpose: Loads runtime configuration from environment variables with safe defaults.
// Inputs:  Environment variables (PORT, DB_SOCKET_PATH, ENV, READ_TIMEOUT_SEC, WRITE_TIMEOUT_SEC)
// Outputs: *Config
// Errors:  None (falls back to defaults)
func Load() *Config {
	return &Config{
		HTTPPort:       getEnv("PORT", "8080"),
		DBSocketPath:   getEnv("DB_SOCKET_PATH", "/var/run/campus-os/db.sock"),
		Environment:    getEnv("ENV", "development"),
		ReadTimeout:    getDurationEnv("READ_TIMEOUT_SEC", 10) * time.Second,
		WriteTimeout:   getDurationEnv("WRITE_TIMEOUT_SEC", 15) * time.Second,
		DefaultTimeout: 5 * time.Second,
	}
}

func getEnv(key, defaultVal string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return defaultVal
}

func getDurationEnv(key string, defaultSec int) time.Duration {
	val := os.Getenv(key)
	if val == "" {
		return time.Duration(defaultSec)
	}
	sec, err := strconv.Atoi(val)
	if err != nil {
		return time.Duration(defaultSec)
	}
	return time.Duration(sec)
}
