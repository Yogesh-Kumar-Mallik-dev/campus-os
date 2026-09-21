package http

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	campusv1 "github.com/Yogesh-Kumar-Mallik-dev/campus-os/backend/pkg/proto/campus/v1"
	"google.golang.org/grpc"
)

// MockAuthServiceClient implements campusv1.AuthServiceClient for unit testing.
type MockAuthServiceClient struct {
	ValidateClaimTokenFunc       func(ctx context.Context, in *campusv1.ValidateClaimTokenRequest, opts ...grpc.CallOption) (*campusv1.ValidateClaimTokenResponse, error)
	VerifySIMAndSendOTPFunc       func(ctx context.Context, in *campusv1.VerifySIMAndSendOTPRequest, opts ...grpc.CallOption) (*campusv1.VerifySIMAndSendOTPResponse, error)
	VerifyOTPAndClaimAccountFunc func(ctx context.Context, in *campusv1.VerifyOTPAndClaimAccountRequest, opts ...grpc.CallOption) (*campusv1.VerifyOTPAndClaimAccountResponse, error)
	LoginFunc                    func(ctx context.Context, in *campusv1.LoginRequest, opts ...grpc.CallOption) (*campusv1.LoginResponse, error)
	RefreshTokenFunc             func(ctx context.Context, in *campusv1.RefreshTokenRequest, opts ...grpc.CallOption) (*campusv1.RefreshTokenResponse, error)
}

func (m *MockAuthServiceClient) BootstrapSuperAdmin(ctx context.Context, in *campusv1.BootstrapSuperAdminRequest, opts ...grpc.CallOption) (*campusv1.BootstrapSuperAdminResponse, error) {
	return nil, nil
}
func (m *MockAuthServiceClient) GenerateExecutiveQR(ctx context.Context, in *campusv1.GenerateExecutiveQRRequest, opts ...grpc.CallOption) (*campusv1.GenerateExecutiveQRResponse, error) {
	return nil, nil
}
func (m *MockAuthServiceClient) GenerateBulkStudentQR(ctx context.Context, in *campusv1.GenerateBulkStudentQRRequest, opts ...grpc.CallOption) (*campusv1.GenerateBulkStudentQRResponse, error) {
	return nil, nil
}
func (m *MockAuthServiceClient) ValidateClaimToken(ctx context.Context, in *campusv1.ValidateClaimTokenRequest, opts ...grpc.CallOption) (*campusv1.ValidateClaimTokenResponse, error) {
	if m.ValidateClaimTokenFunc != nil {
		return m.ValidateClaimTokenFunc(ctx, in, opts...)
	}
	return &campusv1.ValidateClaimTokenResponse{IsValid: true, Username: "yogesh.cse.2024.l"}, nil
}
func (m *MockAuthServiceClient) VerifySIMAndSendOTP(ctx context.Context, in *campusv1.VerifySIMAndSendOTPRequest, opts ...grpc.CallOption) (*campusv1.VerifySIMAndSendOTPResponse, error) {
	if m.VerifySIMAndSendOTPFunc != nil {
		return m.VerifySIMAndSendOTPFunc(ctx, in, opts...)
	}
	return &campusv1.VerifySIMAndSendOTPResponse{SimMatched: true, OtpChallengeId: "otp_test_123", ResendAvailableInSeconds: 45}, nil
}
func (m *MockAuthServiceClient) VerifyOTPAndClaimAccount(ctx context.Context, in *campusv1.VerifyOTPAndClaimAccountRequest, opts ...grpc.CallOption) (*campusv1.VerifyOTPAndClaimAccountResponse, error) {
	if m.VerifyOTPAndClaimAccountFunc != nil {
		return m.VerifyOTPAndClaimAccountFunc(ctx, in, opts...)
	}
	return &campusv1.VerifyOTPAndClaimAccountResponse{Success: true, AccessToken: "acc_123", RefreshToken: "ref_123"}, nil
}
func (m *MockAuthServiceClient) Login(ctx context.Context, in *campusv1.LoginRequest, opts ...grpc.CallOption) (*campusv1.LoginResponse, error) {
	if m.LoginFunc != nil {
		return m.LoginFunc(ctx, in, opts...)
	}
	return &campusv1.LoginResponse{AccessToken: "acc_123", RefreshToken: "ref_123", Username: in.Identifier}, nil
}
func (m *MockAuthServiceClient) RefreshToken(ctx context.Context, in *campusv1.RefreshTokenRequest, opts ...grpc.CallOption) (*campusv1.RefreshTokenResponse, error) {
	if m.RefreshTokenFunc != nil {
		return m.RefreshTokenFunc(ctx, in, opts...)
	}
	return &campusv1.RefreshTokenResponse{AccessToken: "acc_new", RefreshToken: "ref_new"}, nil
}
func (m *MockAuthServiceClient) RevokeSession(ctx context.Context, in *campusv1.RevokeSessionRequest, opts ...grpc.CallOption) (*campusv1.RevokeSessionResponse, error) {
	return &campusv1.RevokeSessionResponse{Success: true}, nil
}

func TestAuthHTTPHandler_ValidateClaim(t *testing.T) {
	mockClient := &MockAuthServiceClient{
		ValidateClaimTokenFunc: func(ctx context.Context, in *campusv1.ValidateClaimTokenRequest, opts ...grpc.CallOption) (*campusv1.ValidateClaimTokenResponse, error) {
			if in.ClaimToken == "valid_tok" {
				return &campusv1.ValidateClaimTokenResponse{
					IsValid:              true,
					UserId:               "u-123",
					Username:             "yogesh.cse.2024.l",
					AcademicName:         "Yogesh",
					LegalFullName:        "Yogesh Kumar Mallik",
					AdmissionType:        "LATERAL_ENTRY",
					EntrySemesterNumber:  3,
					LateralEntrySummary:  "Lateral Entry Verified",
					MaskedPhoneNumber:    "+91 98XXX-XX123",
				}, nil
			}
			return &campusv1.ValidateClaimTokenResponse{IsValid: false}, nil
		},
	}

	handler := NewAuthHTTPHandler(mockClient)

	// Valid token test
	reqBody, _ := json.Marshal(map[string]string{"claim_token": "valid_tok"})
	req := httptest.NewRequest("POST", "/api/v1/auth/claim/validate", bytes.NewReader(reqBody))
	rec := httptest.NewRecorder()
	handler.HandleValidateClaim(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rec.Code)
	}

	var res map[string]any
	_ = json.NewDecoder(rec.Body).Decode(&res)
	if res["username"] != "yogesh.cse.2024.l" {
		t.Errorf("expected username yogesh.cse.2024.l, got %v", res["username"])
	}
	if res["admission_type"] != "LATERAL_ENTRY" {
		t.Errorf("expected admission_type LATERAL_ENTRY, got %v", res["admission_type"])
	}

	// Invalid token test
	invBody, _ := json.Marshal(map[string]string{"claim_token": "expired_tok"})
	invReq := httptest.NewRequest("POST", "/api/v1/auth/claim/validate", bytes.NewReader(invBody))
	invRec := httptest.NewRecorder()
	handler.HandleValidateClaim(invRec, invReq)

	if invRec.Code != http.StatusUnauthorized {
		t.Fatalf("expected status 401 for invalid token, got %d", invRec.Code)
	}
}

func TestAuthHTTPHandler_VerifySIM(t *testing.T) {
	mockClient := &MockAuthServiceClient{
		VerifySIMAndSendOTPFunc: func(ctx context.Context, in *campusv1.VerifySIMAndSendOTPRequest, opts ...grpc.CallOption) (*campusv1.VerifySIMAndSendOTPResponse, error) {
			if in.DeviceCarrierPhone == "+919876543210" {
				return &campusv1.VerifySIMAndSendOTPResponse{
					SimMatched:               true,
					OtpChallengeId:           "otp_chal_1",
					ResendAvailableInSeconds: 45,
				}, nil
			}
			return &campusv1.VerifySIMAndSendOTPResponse{SimMatched: false}, nil
		},
	}

	handler := NewAuthHTTPHandler(mockClient)

	// SIM Match
	matchBody, _ := json.Marshal(map[string]string{
		"claim_token":          "tok_1",
		"device_carrier_phone": "+919876543210",
	})
	req := httptest.NewRequest("POST", "/api/v1/auth/claim/verify-sim", bytes.NewReader(matchBody))
	rec := httptest.NewRecorder()
	handler.HandleVerifySIM(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200 for SIM match, got %d", rec.Code)
	}

	// SIM Mismatch
	mismatchBody, _ := json.Marshal(map[string]string{
		"claim_token":          "tok_1",
		"device_carrier_phone": "+919111122222",
	})
	mReq := httptest.NewRequest("POST", "/api/v1/auth/claim/verify-sim", bytes.NewReader(mismatchBody))
	mRec := httptest.NewRecorder()
	handler.HandleVerifySIM(mRec, mReq)

	if mRec.Code != http.StatusForbidden {
		t.Fatalf("expected 403 for SIM mismatch, got %d", mRec.Code)
	}
}

func TestAuthHTTPHandler_CompleteClaim(t *testing.T) {
	mockClient := &MockAuthServiceClient{}
	handler := NewAuthHTTPHandler(mockClient)

	// Short password (violates grace period min 6)
	shortBody, _ := json.Marshal(map[string]any{
		"claim_token":      "tok_1",
		"otp_challenge_id": "otp_1",
		"otp_code":         "123456",
		"new_password":     "12345",
	})
	req := httptest.NewRequest("POST", "/api/v1/auth/claim/complete", bytes.NewReader(shortBody))
	rec := httptest.NewRecorder()
	handler.HandleCompleteClaim(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected 400 for short password, got %d", rec.Code)
	}

	// Valid claim
	validBody, _ := json.Marshal(map[string]any{
		"claim_token":      "tok_1",
		"otp_challenge_id": "otp_1",
		"otp_code":         "123456",
		"new_password":     "securepassword123",
	})
	vReq := httptest.NewRequest("POST", "/api/v1/auth/claim/complete", bytes.NewReader(validBody))
	vRec := httptest.NewRecorder()
	handler.HandleCompleteClaim(vRec, vReq)

	if vRec.Code != http.StatusOK {
		t.Fatalf("expected 200 for valid claim, got %d", vRec.Code)
	}
}
