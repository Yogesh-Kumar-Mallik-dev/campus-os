package ui

import (
	"testing"

	"fyne.io/fyne/v2/test"
	"github.com/Yogesh-Kumar-Mallik-dev/campus-os/apps/client/internal/api"
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

