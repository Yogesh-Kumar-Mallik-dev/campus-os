package api

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

// BLOCK_API_CLIENT_TEST_001
// Purpose: Verifies client health check integration with mock backend.
func TestClientCheckHealth(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/healthz" {
			t.Errorf("expected /healthz path, got %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"status":"healthy","service":"campus-os-backend","version":"0.1.0"}`))
	}))
	defer ts.Close()

	client := NewClient(ts.URL)
	resp, err := client.CheckHealth(context.Background())
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if resp.Status != "healthy" {
		t.Errorf("expected healthy, got %s", resp.Status)
	}
	if resp.Version != "0.1.0" {
		t.Errorf("expected version 0.1.0, got %s", resp.Version)
	}
}
