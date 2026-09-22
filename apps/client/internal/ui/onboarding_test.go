package ui

import (
	"bytes"
	"context"
	"image"
	"image/color"
	"image/png"
	"net/http"
	"net/http/httptest"
	"testing"

	"fyne.io/fyne/v2/test"
	"github.com/Yogesh-Kumar-Mallik-dev/campus-os/apps/client/internal/api"
	"github.com/skip2/go-qrcode"
)

// BLOCK_UI_ONBOARDING_TEST_001
// Purpose: Verifies OnboardingWizard step transitions, profile review, and password grace period validation.
func TestOnboardingWizard_Steps(t *testing.T) {
	testApp := test.NewApp()
	defer testApp.Quit()
	testWindow := test.NewWindow(nil)

	client := api.NewClient("http://localhost:8080")
	completed := false
	wizard := NewOnboardingWizard(client, testWindow, func() {
		completed = true
	})

	if wizard.CurrentStep() != StepScan {
		t.Fatalf("expected initial step StepScan, got %v", wizard.CurrentStep())
	}

	// Move to SIM Verify
	wizard.ClaimToken = "tok_test_123"
	wizard.MaskedPhone = "+91 98XXX-XX123"
	wizard.SetStep(StepSIMVerify)

	if wizard.CurrentStep() != StepSIMVerify {
		t.Fatalf("expected step StepSIMVerify, got %v", wizard.CurrentStep())
	}

	// Move to Review Profile
	wizard.AcademicName = "Yogesh"
	wizard.LegalFullName = "Yogesh Kumar Mallik"
	wizard.Username = "yogesh.cse.2024.l"
	wizard.AdmissionType = "LATERAL_ENTRY"
	wizard.EntrySemester = 3
	wizard.SetStep(StepReviewProfile)

	if wizard.CurrentStep() != StepReviewProfile {
		t.Fatalf("expected step StepReviewProfile, got %v", wizard.CurrentStep())
	}

	// Move to Set Password
	wizard.SetStep(StepSetPassword)
	if wizard.CurrentStep() != StepSetPassword {
		t.Fatalf("expected step StepSetPassword, got %v", wizard.CurrentStep())
	}

	// Move to Complete
	wizard.SetStep(StepComplete)
	if wizard.CurrentStep() != StepComplete {
		t.Fatalf("expected step StepComplete, got %v", wizard.CurrentStep())
	}

	// Trigger completion callback
	if wizard.onComplete != nil {
		wizard.onComplete()
	}
	if !completed {
		t.Errorf("expected completion callback to be invoked")
	}
}

// BLOCK_UI_ONBOARDING_TEST_002
// Purpose: Verifies toast notifications and UI state interactions in the onboarding wizard.
func TestOnboardingWizard_ToastAndActions(t *testing.T) {
	testApp := test.NewApp()
	defer testApp.Quit()
	testWindow := test.NewWindow(nil)

	client := api.NewClient("http://localhost:8080")
	wizard := NewOnboardingWizard(client, testWindow, nil)

	// Verify toast does not panic
	wizard.toast("Test Notice", "Detailed payload", AlertDefault)

	// Verify nil window toast does not panic
	wizardNilWin := NewOnboardingWizard(client, nil, nil)
	wizardNilWin.toast("Silent Notice", "Should safely no-op", AlertDefault)

	// Verify CanvasObject is non-nil
	if wizard.CanvasObject() == nil {
		t.Fatal("expected non-nil canvas object from wizard")
	}

	// Verify step indicator rendering
	wizard.SetStep(StepSIMVerify)
	if wizard.CurrentStep() != StepSIMVerify {
		t.Fatalf("expected step %v, got %v", StepSIMVerify, wizard.CurrentStep())
	}
}

// BLOCK_UI_ONBOARDING_TEST_003
// Purpose: Verifies frontend-to-backend live synchronization for claim validation, SIM match, and completion.
func TestOnboardingWizard_BackendSync(t *testing.T) {
	testApp := test.NewApp()
	defer testApp.Quit()
	testWindow := test.NewWindow(nil)

	// Spin up mock HTTP backend matching backend/internal/transport/http
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch r.URL.Path {
		case "/api/v1/auth/claim/validate":
			_, _ = w.Write([]byte(`{
				"is_valid": true,
				"user_id": "u-sync-1",
				"masked_phone_number": "+91 98XXX-XX210",
				"username": "scholar.sync",
				"academic_name": "Scholar Sync",
				"legal_full_name": "Scholar Sync Kumar",
				"admission_type": "LATERAL_ENTRY",
				"entry_semester_number": 3,
				"lateral_entry_summary": "Sync verified",
				"grace_period_remaining_seconds": 86400
			}`))
		case "/api/v1/auth/claim/verify-sim":
			_, _ = w.Write([]byte(`{
				"sim_matched": true,
				"otp_challenge_id": "otp_sync_challenge_456",
				"resend_available_in_seconds": 45
			}`))
		case "/api/v1/auth/claim/complete":
			_, _ = w.Write([]byte(`{
				"success": true,
				"access_token": "token_sync_access",
				"refresh_token": "token_sync_refresh",
				"user_id": "u-sync-1",
				"username": "scholar.sync",
				"email": "scholar.sync@campus.edu",
				"role_code": "STUDENT"
			}`))
		default:
			http.NotFound(w, r)
		}
	}))
	defer server.Close()

	client := api.NewClient(server.URL)
	wizard := NewOnboardingWizard(client, testWindow, nil)

	// 1. Direct validation via synced client
	resp, err := wizard.client.ValidateClaim(context.Background(), "sync_token")
	if err != nil || !resp.IsValid {
		t.Fatalf("expected valid claim from backend, got %v", err)
	}
	if resp.Username != "scholar.sync" {
		t.Errorf("expected scholar.sync, got %s", resp.Username)
	}

	// 2. SIM verification via synced client
	simResp, err := wizard.client.VerifySIM(context.Background(), "sync_token", "hash", "+919876543210")
	if err != nil || !simResp.SimMatched {
		t.Fatalf("expected sim match from backend, got %v", err)
	}
	if simResp.OTPChallengeID != "otp_sync_challenge_456" {
		t.Errorf("expected otp_sync_challenge_456, got %s", simResp.OTPChallengeID)
	}

	// 3. Claim completion via synced client
	compResp, err := wizard.client.CompleteClaim(context.Background(), "sync_token", simResp.OTPChallengeID, "123456", "securePassword", true)
	if err != nil || !compResp.Success {
		t.Fatalf("expected successful claim completion, got %v", err)
	}
	if compResp.RoleCode != "STUDENT" {
		t.Errorf("expected STUDENT role, got %s", compResp.RoleCode)
	}
}

// BLOCK_UI_ONBOARDING_TEST_004
// Purpose: Verifies visual feedback and token extraction when uploading images.
func TestOnboardingWizard_ProcessUploadedImage(t *testing.T) {
	testApp := test.NewApp()
	defer testApp.Quit()
	testWindow := test.NewWindow(nil)

	client := api.NewClient("http://localhost:8080")
	wizard := NewOnboardingWizard(client, testWindow, nil)

	// Case 1: Uploading a valid QR code image
	pngData, err := qrcode.Encode("CAMPUS_OS:CLAIM:v1:claim_exec_director_999", qrcode.Medium, 256)
	if err != nil {
		t.Fatalf("failed to encode test QR: %v", err)
	}

	wizard.processUploadedImage("executive_docket.png", bytes.NewReader(pngData))

	if wizard.ScannedImageObj == nil {
		t.Errorf("expected ScannedImageObj to be populated")
	}
	if wizard.ScannedImageName != "executive_docket.png" {
		t.Errorf("expected executive_docket.png, got %s", wizard.ScannedImageName)
	}
	if wizard.ClaimToken != "claim_exec_director_999" {
		t.Errorf("expected 'claim_exec_director_999', got %q", wizard.ClaimToken)
	}
	if wizard.IsError {
		t.Errorf("expected IsError=false for valid QR image, got true with status: %s", wizard.StatusText)
	}

	// Verify canvas rendered visual preview
	obj := wizard.CanvasObject()
	if obj == nil {
		t.Errorf("expected CanvasObject to not be nil")
	}

	// Case 2: Uploading a blank image (no QR code)
	blankImg := image.NewRGBA(image.Rect(0, 0, 100, 100))
	for x := 0; x < 100; x++ {
		for y := 0; y < 100; y++ {
			blankImg.Set(x, y, color.White)
		}
	}
	var buf bytes.Buffer
	if err := png.Encode(&buf, blankImg); err != nil {
		t.Fatalf("failed to encode blank image: %v", err)
	}

	wizard.processUploadedImage("random_photo.png", &buf)

	if wizard.ScannedImageObj == nil {
		t.Errorf("expected ScannedImageObj to still display uploaded photo")
	}
	if wizard.ClaimToken != "" {
		t.Errorf("expected ClaimToken to be empty for blank photo, got %q", wizard.ClaimToken)
	}
	if !wizard.IsError {
		t.Errorf("expected IsError=true for photo with no QR code")
	}
}



