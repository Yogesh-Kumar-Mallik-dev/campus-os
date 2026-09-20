package http

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

// BLOCK_ROUTER_TEST_001
// Purpose: Verifies health check and API base routes with middleware execution.
func TestRouterEndpoints(t *testing.T) {
	router := NewRouter(RouterConfig{Version: "0.1.0"})

	// 1. Healthz test
	req := httptest.NewRequest(http.MethodGet, "/healthz", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected healthz status 200, got %d", w.Code)
	}
	if reqID := w.Header().Get("X-Request-ID"); reqID == "" {
		t.Error("expected X-Request-ID header to be present")
	}

	var healthRes map[string]any
	if err := json.NewDecoder(w.Body).Decode(&healthRes); err != nil {
		t.Fatalf("failed to decode health response: %v", err)
	}
	if healthRes["status"] != "healthy" {
		t.Errorf("expected status 'healthy', got %v", healthRes["status"])
	}

	// 2. Base API info test
	reqAPI := httptest.NewRequest(http.MethodGet, "/api/v1", nil)
	wAPI := httptest.NewRecorder()
	router.ServeHTTP(wAPI, reqAPI)

	if wAPI.Code != http.StatusOK {
		t.Errorf("expected /api/v1 status 200, got %d", wAPI.Code)
	}
}
