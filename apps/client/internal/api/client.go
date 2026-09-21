package api

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// Client handles network interaction with the Campus OS Go Logical Backend.
type Client struct {
	baseURL    string
	httpClient *http.Client
}

// HealthResponse represents the health payload from backend.
type HealthResponse struct {
	Status  string `json:"status"`
	Service string `json:"service"`
	Version string `json:"version"`
}

// BLOCK_CLIENT_API_INIT_001
// Purpose: Constructs a new backend API client with default timeout.
func NewClient(baseURL string) *Client {
	if baseURL == "" {
		baseURL = "http://localhost:8080"
	}
	return &Client{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// BLOCK_CLIENT_API_HEALTH_001
// Purpose: Pings backend /healthz endpoint to check connectivity.
func (c *Client) CheckHealth(ctx context.Context) (*HealthResponse, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.baseURL+"/healthz", nil)
	if err != nil {
		return nil, fmt.Errorf("BLOCK_CLIENT_API_HEALTH_001: failed to create request: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("BLOCK_CLIENT_API_HEALTH_001: backend unreachable: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("BLOCK_CLIENT_API_HEALTH_001: unexpected status %d", resp.StatusCode)
	}

	var res HealthResponse
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return nil, fmt.Errorf("BLOCK_CLIENT_API_HEALTH_001: failed to decode health response: %w", err)
	}
	return &res, nil
}

// ClaimValidationResponse represents profile information returned for a valid claim token.
type ClaimValidationResponse struct {
	IsValid                     bool   `json:"is_valid"`
	UserID                      string `json:"user_id"`
	MaskedPhoneNumber           string `json:"masked_phone_number"`
	Username                    string `json:"username"`
	AcademicName                string `json:"academic_name"`
	LegalFullName               string `json:"legal_full_name"`
	AdmissionType               string `json:"admission_type"`
	EntrySemesterNumber         int    `json:"entry_semester_number"`
	LateralEntrySummary         string `json:"lateral_entry_summary"`
	GracePeriodRemainingSeconds int64  `json:"grace_period_remaining_seconds"`
}

// BLOCK_CLIENT_API_CLAIM_VAL_001
// Purpose: Validates a scanned QR claim token with the backend.
func (c *Client) ValidateClaim(ctx context.Context, claimToken string) (*ClaimValidationResponse, error) {
	reqBody, _ := json.Marshal(map[string]string{"claim_token": claimToken})
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL+"/api/v1/auth/claim/validate", bytes.NewReader(reqBody))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("backend unreachable: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("invalid claim token (status %d)", resp.StatusCode)
	}

	var res ClaimValidationResponse
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return nil, err
	}
	return &res, nil
}

// SIMVerificationResponse confirms hardware SIM match and OTP dispatch.
type SIMVerificationResponse struct {
	SimMatched               bool   `json:"sim_matched"`
	OTPChallengeID           string `json:"otp_challenge_id"`
	ResendAvailableInSeconds int    `json:"resend_available_in_seconds"`
}

// BLOCK_CLIENT_API_SIM_VERIFY_001
// Purpose: Submits telephony SIM info to verify hardware binding and dispatch an OTP.
func (c *Client) VerifySIM(ctx context.Context, claimToken, simHash, phone string) (*SIMVerificationResponse, error) {
	reqBody, _ := json.Marshal(map[string]string{
		"claim_token":           claimToken,
		"device_sim_iccid_hash": simHash,
		"device_carrier_phone":  phone,
	})
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL+"/api/v1/auth/claim/verify-sim", bytes.NewReader(reqBody))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("backend unreachable: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusForbidden {
		return &SIMVerificationResponse{SimMatched: false}, nil
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("SIM verification failed (status %d)", resp.StatusCode)
	}

	var res SIMVerificationResponse
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return nil, err
	}
	return &res, nil
}

// ClaimCompletionResponse contains the issued credentials and session tokens.
type ClaimCompletionResponse struct {
	Success      bool   `json:"success"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	UserID       string `json:"user_id"`
	Username     string `json:"username"`
	Email        string `json:"email"`
	RoleCode     string `json:"role_code"`
}

// BLOCK_CLIENT_API_CLAIM_COMPLETE_001
// Purpose: Submits OTP and initial password to complete account activation.
func (c *Client) CompleteClaim(ctx context.Context, claimToken, otpID, otpCode, password string, biometrics bool) (*ClaimCompletionResponse, error) {
	reqBody, _ := json.Marshal(map[string]any{
		"claim_token":       claimToken,
		"otp_challenge_id":  otpID,
		"otp_code":          otpCode,
		"new_password":      password,
		"enable_biometrics": biometrics,
	})
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL+"/api/v1/auth/claim/complete", bytes.NewReader(reqBody))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("backend unreachable: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("account claim failed (status %d)", resp.StatusCode)
	}

	var res ClaimCompletionResponse
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return nil, err
	}
	return &res, nil
}
