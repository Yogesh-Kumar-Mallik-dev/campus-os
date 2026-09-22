package http

import (
	"encoding/json"
	"net/http"
	"strings"

	campusv1 "github.com/Yogesh-Kumar-Mallik-dev/campus-os/backend/pkg/proto/campus/v1"
	"google.golang.org/grpc/codes"
	grpcStatus "google.golang.org/grpc/status"
)

// mapGRPCError translates gRPC transport and status errors into clean RFC 7807 problem details
// preventing internal unix domain socket paths or raw transport traces from leaking to clients.
func mapGRPCError(err error, fallbackCode, fallbackTitle string) (int, string, string, string) {
	if err == nil {
		return http.StatusOK, "", "", ""
	}
	st, ok := grpcStatus.FromError(err)
	if ok {
		switch st.Code() {
		case codes.Unavailable:
			return http.StatusServiceUnavailable, "SERVICE_UNAVAILABLE", "Database Persistence Offline", "Unable to connect to local database persistence layer. Please ensure the database daemon is active."
		case codes.DeadlineExceeded:
			return http.StatusGatewayTimeout, "PERSISTENCE_TIMEOUT", "Database Request Timeout", "The database persistence layer timed out processing your request."
		case codes.NotFound:
			return http.StatusNotFound, "NOT_FOUND", fallbackTitle, st.Message()
		case codes.AlreadyExists:
			return http.StatusConflict, "ALREADY_EXISTS", fallbackTitle, st.Message()
		case codes.PermissionDenied:
			return http.StatusForbidden, "PERMISSION_DENIED", fallbackTitle, st.Message()
		case codes.Unauthenticated:
			return http.StatusUnauthorized, "UNAUTHENTICATED", fallbackTitle, st.Message()
		case codes.InvalidArgument:
			return http.StatusBadRequest, "INVALID_ARGUMENT", fallbackTitle, st.Message()
		}
		return http.StatusInternalServerError, fallbackCode, fallbackTitle, st.Message()
	}
	errMsg := err.Error()
	if strings.Contains(errMsg, "dial unix") || strings.Contains(errMsg, "connection error") || strings.Contains(errMsg, "no such file or directory") {
		return http.StatusServiceUnavailable, "SERVICE_UNAVAILABLE", "Database Persistence Offline", "Unable to connect to local database persistence layer. Please ensure the database daemon is active."
	}
	return http.StatusInternalServerError, fallbackCode, fallbackTitle, errMsg
}

// AuthHTTPHandler maps incoming HTTP onboarding and login requests to the gRPC AuthService.
type AuthHTTPHandler struct {
	client campusv1.AuthServiceClient
}

// NewAuthHTTPHandler constructs an AuthHTTPHandler with the provided gRPC client.
func NewAuthHTTPHandler(client campusv1.AuthServiceClient) *AuthHTTPHandler {
	return &AuthHTTPHandler{client: client}
}

// BLOCK_AUTH_HTTP_VALIDATE_001
// Purpose: Validates a scanned single-use claim token and returns profile details.
func (h *AuthHTTPHandler) HandleValidateClaim(w http.ResponseWriter, r *http.Request) {
	var body struct {
		ClaimToken string `json:"claim_token"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil || body.ClaimToken == "" {
		WriteProblem(w, r, http.StatusBadRequest, "INVALID_CLAIM_REQUEST", "Invalid Claim Request", "claim_token is required", nil)
		return
	}

	if h.client == nil {
		WriteProblem(w, r, http.StatusServiceUnavailable, "SERVICE_UNAVAILABLE", "Auth Service Unavailable", "Backend persistence link is uninitialized", nil)
		return
	}

	resp, err := h.client.ValidateClaimToken(r.Context(), &campusv1.ValidateClaimTokenRequest{
		ClaimToken: body.ClaimToken,
	})
	if err != nil {
		status, code, title, detail := mapGRPCError(err, "CLAIM_VALIDATION_FAILED", "Claim Validation Error")
		WriteProblem(w, r, status, code, title, detail, nil)
		return
	}

	if !resp.IsValid {
		WriteProblem(w, r, http.StatusUnauthorized, "CLAIM_TOKEN_INVALID", "Invalid Claim Token", "The token is expired, consumed, or invalid", nil)
		return
	}

	WriteJSON(w, http.StatusOK, map[string]any{
		"is_valid":                        resp.IsValid,
		"user_id":                         resp.UserId,
		"masked_phone_number":             resp.MaskedPhoneNumber,
		"username":                        resp.Username,
		"academic_name":                   resp.AcademicName,
		"legal_full_name":                 resp.LegalFullName,
		"admission_type":                  resp.AdmissionType,
		"entry_semester_number":           resp.EntrySemesterNumber,
		"lateral_entry_summary":           resp.LateralEntrySummary,
		"grace_period_remaining_seconds": resp.GracePeriodRemainingSeconds,
	})
}

// BLOCK_AUTH_HTTP_VERIFY_SIM_001
// Purpose: Checks hardware SIM telephony match and triggers an SMS OTP.
func (h *AuthHTTPHandler) HandleVerifySIM(w http.ResponseWriter, r *http.Request) {
	var body struct {
		ClaimToken         string `json:"claim_token"`
		DeviceSIMIccidHash string `json:"device_sim_iccid_hash"`
		DeviceCarrierPhone string `json:"device_carrier_phone"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil || body.ClaimToken == "" {
		WriteProblem(w, r, http.StatusBadRequest, "INVALID_SIM_REQUEST", "Invalid SIM Request", "claim_token and phone are required", nil)
		return
	}

	if h.client == nil {
		WriteProblem(w, r, http.StatusServiceUnavailable, "SERVICE_UNAVAILABLE", "Auth Service Unavailable", "Backend persistence link is uninitialized", nil)
		return
	}

	resp, err := h.client.VerifySIMAndSendOTP(r.Context(), &campusv1.VerifySIMAndSendOTPRequest{
		ClaimToken:         body.ClaimToken,
		DeviceSimIccidHash: body.DeviceSIMIccidHash,
		DeviceCarrierPhone: body.DeviceCarrierPhone,
	})
	if err != nil {
		status, code, title, detail := mapGRPCError(err, "SIM_VERIFY_FAILED", "SIM Verification Error")
		WriteProblem(w, r, status, code, title, detail, nil)
		return
	}

	if !resp.SimMatched {
		WriteProblem(w, r, http.StatusForbidden, "ERR_SIM_ABSENT_OR_MISMATCH", "SIM Mismatch or Absent", "The physical SIM matching this student record is not present in this device", nil)
		return
	}

	WriteJSON(w, http.StatusOK, map[string]any{
		"sim_matched":                 resp.SimMatched,
		"otp_challenge_id":            resp.OtpChallengeId,
		"resend_available_in_seconds": resp.ResendAvailableInSeconds,
	})
}

// BLOCK_AUTH_HTTP_COMPLETE_001
// Purpose: Verifies the SMS OTP, sets user password, and completes account claim.
func (h *AuthHTTPHandler) HandleCompleteClaim(w http.ResponseWriter, r *http.Request) {
	var body struct {
		ClaimToken       string `json:"claim_token"`
		OTPChallengeID   string `json:"otp_challenge_id"`
		OTPCode          string `json:"otp_code"`
		NewPassword      string `json:"new_password"`
		EnableBiometrics bool   `json:"enable_biometrics"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil || body.ClaimToken == "" || body.NewPassword == "" {
		WriteProblem(w, r, http.StatusBadRequest, "INVALID_COMPLETE_REQUEST", "Invalid Request", "Missing required claim parameters", nil)
		return
	}

	if len(body.NewPassword) < 6 {
		WriteProblem(w, r, http.StatusBadRequest, "PASSWORD_TOO_SHORT", "Weak Password", "Password must be at least 6 characters during grace period", nil)
		return
	}

	if h.client == nil {
		WriteProblem(w, r, http.StatusServiceUnavailable, "SERVICE_UNAVAILABLE", "Auth Service Unavailable", "Backend persistence link is uninitialized", nil)
		return
	}

	resp, err := h.client.VerifyOTPAndClaimAccount(r.Context(), &campusv1.VerifyOTPAndClaimAccountRequest{
		ClaimToken:       body.ClaimToken,
		OtpChallengeId:   body.OTPChallengeID,
		OtpCode:          body.OTPCode,
		NewPassword:      body.NewPassword,
		EnableBiometrics: body.EnableBiometrics,
	})
	if err != nil {
		status, code, title, detail := mapGRPCError(err, "CLAIM_COMPLETION_FAILED", "Claim Failed")
		WriteProblem(w, r, status, code, title, detail, nil)
		return
	}

	WriteJSON(w, http.StatusOK, map[string]any{
		"success":       resp.Success,
		"access_token":  resp.AccessToken,
		"refresh_token": resp.RefreshToken,
		"user_id":       resp.UserId,
		"username":      resp.Username,
		"email":         resp.Email,
		"role_code":     resp.RoleCode,
	})
}

// BLOCK_AUTH_HTTP_LOGIN_001
// Purpose: Universal authentication endpoint supporting username, email, or PRN.
func (h *AuthHTTPHandler) HandleLogin(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Identifier string `json:"identifier"`
		Password   string `json:"password"`
		TOTPCode   string `json:"totp_code"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil || body.Identifier == "" || body.Password == "" {
		WriteProblem(w, r, http.StatusBadRequest, "INVALID_LOGIN_REQUEST", "Invalid Login Request", "identifier and password are required", nil)
		return
	}

	if h.client == nil {
		WriteProblem(w, r, http.StatusServiceUnavailable, "SERVICE_UNAVAILABLE", "Auth Service Unavailable", "Backend persistence link is uninitialized", nil)
		return
	}

	resp, err := h.client.Login(r.Context(), &campusv1.LoginRequest{
		Identifier: body.Identifier,
		Password:   body.Password,
		TotpCode:   body.TOTPCode,
	})
	if err != nil {
		status, code, title, detail := mapGRPCError(err, "LOGIN_FAILED", "Invalid Credentials")
		if status == http.StatusInternalServerError {
			status = http.StatusUnauthorized
		}
		WriteProblem(w, r, status, code, title, detail, nil)
		return
	}

	WriteJSON(w, http.StatusOK, map[string]any{
		"access_token":       resp.AccessToken,
		"refresh_token":      resp.RefreshToken,
		"expires_in_seconds": resp.ExpiresInSeconds,
		"user_id":            resp.UserId,
		"username":           resp.Username,
		"full_name":          resp.FullName,
		"role_codes":         resp.RoleCodes,
	})
}

// BLOCK_AUTH_HTTP_REFRESH_001
// Purpose: Refreshes access tokens via rotating family refresh tokens.
func (h *AuthHTTPHandler) HandleRefresh(w http.ResponseWriter, r *http.Request) {
	var body struct {
		RefreshToken string `json:"refresh_token"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil || body.RefreshToken == "" {
		WriteProblem(w, r, http.StatusBadRequest, "INVALID_REFRESH_REQUEST", "Invalid Refresh Request", "refresh_token is required", nil)
		return
	}

	if h.client == nil {
		WriteProblem(w, r, http.StatusServiceUnavailable, "SERVICE_UNAVAILABLE", "Auth Service Unavailable", "Backend persistence link is uninitialized", nil)
		return
	}

	resp, err := h.client.RefreshToken(r.Context(), &campusv1.RefreshTokenRequest{
		RefreshToken: body.RefreshToken,
	})
	if err != nil {
		status, code, title, detail := mapGRPCError(err, "REFRESH_FAILED", "Refresh Failed")
		if status == http.StatusInternalServerError {
			status = http.StatusUnauthorized
		}
		WriteProblem(w, r, status, code, title, detail, nil)
		return
	}

	WriteJSON(w, http.StatusOK, map[string]any{
		"access_token":       resp.AccessToken,
		"refresh_token":      resp.RefreshToken,
		"expires_in_seconds": resp.ExpiresInSeconds,
	})
}

// BLOCK_AUTH_HTTP_REVOKE_001
// Purpose: Revokes active session tokens and family refresh tokens.
func (h *AuthHTTPHandler) HandleRevoke(w http.ResponseWriter, r *http.Request) {
	var body struct {
		RefreshToken string `json:"refresh_token"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil || body.RefreshToken == "" {
		WriteProblem(w, r, http.StatusBadRequest, "INVALID_REVOKE_REQUEST", "Invalid Revocation Request", "refresh_token is required", nil)
		return
	}

	if h.client != nil {
		_, _ = h.client.RevokeSession(r.Context(), &campusv1.RevokeSessionRequest{
			RefreshToken: body.RefreshToken,
		})
	}

	WriteJSON(w, http.StatusOK, map[string]any{
		"revoked": true,
	})
}

// BLOCK_AUTH_HTTP_EXEC_QR_001
// Purpose: Provisions a sealed single-use QR credential for Tier-1 institutional executives.
func (h *AuthHTTPHandler) HandleGenerateExecutiveQR(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Role        string `json:"role"`
		FullName    string `json:"full_name"`
		Email       string `json:"email"`
		PhoneNumber string `json:"phone_number"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil || body.FullName == "" || body.Email == "" || body.PhoneNumber == "" {
		WriteProblem(w, r, http.StatusBadRequest, "INVALID_EXEC_REQUEST", "Invalid Request", "role, full_name, email, and phone_number are required", nil)
		return
	}

	var execRole campusv1.ExecutiveRole
	normalized := strings.ToUpper(strings.ReplaceAll(body.Role, " ", "_"))
	switch {
	case strings.Contains(normalized, "EXECUTIVE_DIRECTOR"):
		execRole = campusv1.ExecutiveRole_EXECUTIVE_ROLE_EXECUTIVE_DIRECTOR
	case strings.Contains(normalized, "DIRECTOR"):
		execRole = campusv1.ExecutiveRole_EXECUTIVE_ROLE_DIRECTOR
	case strings.Contains(normalized, "DEAN"):
		execRole = campusv1.ExecutiveRole_EXECUTIVE_ROLE_DEAN
	case strings.Contains(normalized, "REGISTRAR"):
		execRole = campusv1.ExecutiveRole_EXECUTIVE_ROLE_REGISTRAR
	default:
		execRole = campusv1.ExecutiveRole_EXECUTIVE_ROLE_DIRECTOR
	}

	if h.client == nil {
		WriteProblem(w, r, http.StatusServiceUnavailable, "SERVICE_UNAVAILABLE", "Auth Service Unavailable", "Backend persistence link is uninitialized", nil)
		return
	}

	resp, err := h.client.GenerateExecutiveQR(r.Context(), &campusv1.GenerateExecutiveQRRequest{
		Role:        execRole,
		FullName:    body.FullName,
		Email:       body.Email,
		PhoneNumber: body.PhoneNumber,
	})
	if err != nil {
		status, code, title, detail := mapGRPCError(err, "EXEC_QR_FAILED", "Failed to generate executive QR")
		WriteProblem(w, r, status, code, title, detail, nil)
		return
	}

	WriteJSON(w, http.StatusOK, map[string]any{
		"claim_token":       resp.ClaimToken,
		"sealed_qr_payload": resp.SealedQrPayload,
		"expires_at_unix":   resp.ExpiresAtUnix,
	})
}

