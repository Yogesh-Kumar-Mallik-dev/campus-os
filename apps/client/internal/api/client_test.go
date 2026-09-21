package api

import (
	"context"
	"errors"
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

// BLOCK_API_CLIENT_TEST_002
// Purpose: Verifies ValidateClaim for both success and RFC 7807 rejection.
func TestClientValidateClaim(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v1/auth/claim/validate" {
			t.Errorf("unexpected path %s", r.URL.Path)
		}
		if r.Header.Get("Content-Type") != "application/json" {
			t.Errorf("expected application/json, got %s", r.Header.Get("Content-Type"))
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{
			"is_valid": true,
			"user_id": "u-123",
			"masked_phone_number": "+91 98XXX-XX210",
			"username": "yogesh.cse.2024.l",
			"academic_name": "Yogesh",
			"legal_full_name": "Yogesh Kumar Mallik",
			"admission_type": "LATERAL_ENTRY",
			"entry_semester_number": 3,
			"lateral_entry_summary": "Direct entry to Sem 3",
			"grace_period_remaining_seconds": 259200
		}`))
	}))
	defer ts.Close()

	client := NewClient(ts.URL)
	resp, err := client.ValidateClaim(context.Background(), "claim_test_123")
	if err != nil {
		t.Fatalf("expected success, got %v", err)
	}
	if !resp.IsValid || resp.Username != "yogesh.cse.2024.l" || resp.EntrySemesterNumber != 3 {
		t.Errorf("unexpected claim response: %+v", resp)
	}

	// Test RFC 7807 problem details error
	errTs := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/problem+json")
		w.WriteHeader(http.StatusUnauthorized)
		_, _ = w.Write([]byte(`{"title":"Invalid Claim Token","detail":"The token is expired, consumed, or invalid","status":401}`))
	}))
	defer errTs.Close()

	errClient := NewClient(errTs.URL)
	_, err = errClient.ValidateClaim(context.Background(), "bad_token")
	if err == nil {
		t.Fatal("expected error for invalid token")
	}
	expected := "Invalid Claim Token: The token is expired, consumed, or invalid"
	if err.Error() != expected {
		t.Errorf("expected %q, got %q", expected, err.Error())
	}
}

// BLOCK_API_CLIENT_TEST_003
// Purpose: Verifies VerifySIM for match and 403 Forbidden SIM mismatch.
func TestClientVerifySIM(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v1/auth/claim/verify-sim" {
			t.Errorf("unexpected path %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{
			"sim_matched": true,
			"otp_challenge_id": "otp_chal_999",
			"resend_available_in_seconds": 45
		}`))
	}))
	defer ts.Close()

	client := NewClient(ts.URL)
	resp, err := client.VerifySIM(context.Background(), "token", "hash", "+919876543210")
	if err != nil {
		t.Fatalf("expected success, got %v", err)
	}
	if !resp.SimMatched || resp.OTPChallengeID != "otp_chal_999" {
		t.Errorf("unexpected response %+v", resp)
	}

	// Mismatch returns 403 Forbidden
	mismatchTs := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusForbidden)
	}))
	defer mismatchTs.Close()

	mismatchClient := NewClient(mismatchTs.URL)
	resp, err = mismatchClient.VerifySIM(context.Background(), "token", "hash", "+919123456789")
	if err != nil {
		t.Fatalf("expected nil error on 403, got %v", err)
	}
	if resp.SimMatched {
		t.Errorf("expected SimMatched to be false on 403")
	}
}

// BLOCK_API_CLIENT_TEST_004
// Purpose: Verifies CompleteClaim for activation success and RFC 7807 error handling.
func TestClientCompleteClaim(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v1/auth/claim/complete" {
			t.Errorf("unexpected path %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{
			"success": true,
			"access_token": "jwt_access_xyz",
			"refresh_token": "jwt_refresh_xyz",
			"user_id": "u-123",
			"username": "yogesh.cse.2024.l",
			"email": "yogesh@campus.edu",
			"role_code": "STUDENT"
		}`))
	}))
	defer ts.Close()

	client := NewClient(ts.URL)
	resp, err := client.CompleteClaim(context.Background(), "token", "otp-id", "123456", "secretPass", true)
	if err != nil {
		t.Fatalf("expected success, got %v", err)
	}
	if !resp.Success || resp.Username != "yogesh.cse.2024.l" {
		t.Errorf("unexpected complete response %+v", resp)
	}
}

// BLOCK_API_CLIENT_TEST_005
// Purpose: Verifies IsUnreachable helper for various network failure scenarios.
func TestIsUnreachable(t *testing.T) {
	if IsUnreachable(nil) {
		t.Errorf("expected nil error to not be unreachable")
	}
	if !IsUnreachable(errors.New("backend unreachable: dial tcp 127.0.0.1:8080: connect: connection refused")) {
		t.Errorf("expected connection refused to be unreachable")
	}
	if IsUnreachable(errors.New("invalid claim token (status 401)")) {
		t.Errorf("expected HTTP 401 to not be unreachable")
	}
}
