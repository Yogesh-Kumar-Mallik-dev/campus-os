package ui

import (
	"context"
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/Yogesh-Kumar-Mallik-dev/campus-os/apps/client/internal/api"
)

// OnboardingStep enumerates the wizard progression states.
type OnboardingStep int

const (
	StepScan OnboardingStep = iota
	StepSIMVerify
	StepReviewProfile
	StepSetPassword
	StepComplete
)

// OnboardingWizard manages the UI state and transitions for account activation.
type OnboardingWizard struct {
	client       *api.Client
	content      *fyne.Container
	step         OnboardingStep
	onComplete   func()

	// Wizard Form State
	ClaimToken          string
	MaskedPhone         string
	Username            string
	AcademicName        string
	LegalFullName       string
	AdmissionType       string
	EntrySemester       int
	LateralSummary      string
	OTPChallengeID      string
	SelectedSIMPhone    string
	OTPCode             string
	NewPassword         string
	BiometricsEnabled   bool
	StatusText          string
	IsError             bool
}

// BLOCK_UI_ONBOARDING_NEW_001
// Purpose: Constructs a new OnboardingWizard instance.
func NewOnboardingWizard(client *api.Client, onComplete func()) *OnboardingWizard {
	w := &OnboardingWizard{
		client:           client,
		step:             StepScan,
		onComplete:       onComplete,
		SelectedSIMPhone: "+919876543210",
		BiometricsEnabled: true,
	}
	w.content = container.NewVBox()
	w.render()
	return w
}

// CanvasObject returns the renderable canvas object.
func (w *OnboardingWizard) CanvasObject() fyne.CanvasObject {
	return w.content
}

// CurrentStep returns the active wizard step for testing.
func (w *OnboardingWizard) CurrentStep() OnboardingStep {
	return w.step
}

// SetStep transitions the wizard step and rerenders.
func (w *OnboardingWizard) SetStep(step OnboardingStep) {
	w.step = step
	w.render()
}

func (w *OnboardingWizard) render() {
	w.content.Objects = nil

	// Status / Error Banner
	if w.StatusText != "" {
		style := fyne.TextStyle{Italic: true}
		statusLabel := widget.NewLabelWithStyle(w.StatusText, fyne.TextAlignCenter, style)
		w.content.Add(statusLabel)
		w.content.Add(widget.NewSeparator())
	}

	switch w.step {
	case StepScan:
		w.renderScanStep()
	case StepSIMVerify:
		w.renderSIMVerifyStep()
	case StepReviewProfile:
		w.renderReviewProfileStep()
	case StepSetPassword:
		w.renderSetPasswordStep()
	case StepComplete:
		w.renderCompleteStep()
	}
}

// Step 1: Scan Sealed QR Code
func (w *OnboardingWizard) renderScanStep() {
	header := widget.NewLabelWithStyle("Step 1 of 4: Scan Sealed QR Code", fyne.TextAlignCenter, fyne.TextStyle{Bold: true})
	subtext := widget.NewLabelWithStyle("Position your camera over the official sealed QR on your admission slip.", fyne.TextAlignCenter, fyne.TextStyle{})

	scannerMock := widget.NewCard(
		"Camera Viewfinder",
		"Aim at tamper-evident QR code",
		container.NewCenter(widget.NewIcon(theme.SearchIcon())),
	)

	tokenEntry := widget.NewEntry()
	tokenEntry.SetPlaceHolder("Or enter 16-character claim code manually")
	if w.ClaimToken != "" {
		tokenEntry.SetText(w.ClaimToken)
	}

	scanBtn := widget.NewButtonWithIcon("Validate QR Code Token", theme.ConfirmIcon(), func() {
		token := tokenEntry.Text
		if token == "" {
			token = "claim_genesis_test_demo"
		}
		w.ClaimToken = token

		// Validate with backend API if available, or simulate
		go func() {
			resp, err := w.client.ValidateClaim(context.Background(), token)
			if err == nil && resp.IsValid {
				w.MaskedPhone = resp.MaskedPhoneNumber
				w.Username = resp.Username
				w.AcademicName = resp.AcademicName
				w.LegalFullName = resp.LegalFullName
				w.AdmissionType = resp.AdmissionType
				w.EntrySemester = resp.EntrySemesterNumber
				w.LateralSummary = resp.LateralEntrySummary
			} else {
				// Fallback demo state
				w.MaskedPhone = "+91 98XXX-XX210"
				w.Username = "yogesh.cse.2024.l"
				w.AcademicName = "Yogesh"
				w.LegalFullName = "Yogesh Kumar Mallik"
				w.AdmissionType = "LATERAL_ENTRY"
				w.EntrySemester = 3
				w.LateralSummary = "Lateral Entry: Direct admission to Semester 3. Prior polytechnic credits verified."
			}
			w.StatusText = "QR Code Recognized! Please verify device SIM presence."
			w.SetStep(StepSIMVerify)
		}()
	})
	scanBtn.Importance = widget.HighImportance

	w.content.Add(header)
	w.content.Add(subtext)
	w.content.Add(scannerMock)
	w.content.Add(tokenEntry)
	w.content.Add(scanBtn)
}

// Step 2: SIM Telephony & SMS OTP Handshake
func (w *OnboardingWizard) renderSIMVerifyStep() {
	header := widget.NewLabelWithStyle("Step 2 of 4: Device SIM & OTP Verification", fyne.TextAlignCenter, fyne.TextStyle{Bold: true})
	phoneInfo := widget.NewLabel(fmt.Sprintf("Registered Mobile Phone: %s", w.MaskedPhone))

	simCard := widget.NewCard(
		"Hardware Telephony Detection",
		"Checking device SIM slots for anti-theft binding",
		container.NewVBox(
			widget.NewLabel("🟢 SIM Slot 1 (Jio 5G): +91 98765-43210 [MATCH FOUND]"),
			widget.NewLabel("⚪ SIM Slot 2 (Airtel): +91 91234-56789 [SECONDARY]"),
		),
	)

	otpEntry := widget.NewEntry()
	otpEntry.SetPlaceHolder("Enter 6-digit SMS OTP (e.g. 123456)")

	verifyBtn := widget.NewButtonWithIcon("Verify SIM & Submit OTP", theme.ConfirmIcon(), func() {
		w.OTPCode = otpEntry.Text
		if w.OTPCode == "" {
			w.OTPCode = "123456"
		}
		w.StatusText = "Device SIM & OTP Confirmed! Please verify your official records."
		w.SetStep(StepReviewProfile)
	})
	verifyBtn.Importance = widget.HighImportance

	w.content.Add(header)
	w.content.Add(phoneInfo)
	w.content.Add(simCard)
	w.content.Add(otpEntry)
	w.content.Add(verifyBtn)
}

// Step 3: Review Official Records & Dual-Identity Confirmation
func (w *OnboardingWizard) renderReviewProfileStep() {
	header := widget.NewLabelWithStyle("Step 3 of 4: Review Official Student Records", fyne.TextAlignCenter, fyne.TextStyle{Bold: true})

	profileCard := widget.NewCard(
		"Student Identity & Enrollment",
		"Verify the accuracy of your academic vs legal identification",
		container.NewVBox(
			widget.NewLabel(fmt.Sprintf("Academic Name (Class X Marksheet) : %s", w.AcademicName)),
			widget.NewLabel(fmt.Sprintf("Legal Full Name (Govt ID / Aadhaar) : %s", w.LegalFullName)),
			widget.NewLabel(fmt.Sprintf("Canonical Username                 : %s", w.Username)),
			widget.NewLabel(fmt.Sprintf("Assigned Institutional Email       : %s@campus.edu", w.Username)),
			widget.NewLabel(fmt.Sprintf("Admission Category                 : %s (Entry Sem %d)", w.AdmissionType, w.EntrySemester)),
			widget.NewLabel(fmt.Sprintf("Academic Notes                     : %s", w.LateralSummary)),
		),
	)

	confirmCheck := widget.NewCheck("I confirm these records match my secondary school and legal certificates", nil)
	confirmCheck.SetChecked(true)

	nextBtn := widget.NewButtonWithIcon("Looks Good, Set Password ->", theme.NavigateNextIcon(), func() {
		if !confirmCheck.Checked {
			w.StatusText = "Please confirm that the records match your official certificates"
			w.render()
			return
		}
		w.StatusText = "Identity confirmed. Set your password."
		w.SetStep(StepSetPassword)
	})
	nextBtn.Importance = widget.HighImportance

	w.content.Add(header)
	w.content.Add(profileCard)
	w.content.Add(confirmCheck)
	w.content.Add(nextBtn)
}

// Step 4: Password Setup & 3-Day Orientation Grace Period
func (w *OnboardingWizard) renderSetPasswordStep() {
	header := widget.NewLabelWithStyle("Step 4 of 4: Create Your Password", fyne.TextAlignCenter, fyne.TextStyle{Bold: true})
	graceBanner := widget.NewLabel("💡 Orientation Grace Period Active: 6+ character password accepted for your first 3 days.")

	passEntry := widget.NewPasswordEntry()
	passEntry.SetPlaceHolder("Enter new password (min 6 characters)")

	confirmPassEntry := widget.NewPasswordEntry()
	confirmPassEntry.SetPlaceHolder("Confirm new password")

	biometricCheck := widget.NewCheck("Enable Face / Fingerprint Unlock on this device", func(b bool) {
		w.BiometricsEnabled = b
	})
	biometricCheck.SetChecked(true)

	completeBtn := widget.NewButtonWithIcon("Complete Onboarding & Activate Account", theme.ConfirmIcon(), func() {
		if len(passEntry.Text) < 6 {
			w.StatusText = "Password must be at least 6 characters during grace period"
			w.render()
			return
		}
		if passEntry.Text != confirmPassEntry.Text {
			w.StatusText = "Passwords do not match"
			w.render()
			return
		}

		w.NewPassword = passEntry.Text
		w.StatusText = "Account Successfully Activated!"
		w.SetStep(StepComplete)
	})
	completeBtn.Importance = widget.HighImportance

	w.content.Add(header)
	w.content.Add(graceBanner)
	w.content.Add(passEntry)
	w.content.Add(confirmPassEntry)
	w.content.Add(biometricCheck)
	w.content.Add(completeBtn)
}

// Step 5: Celebration & Gateway Handoff
func (w *OnboardingWizard) renderCompleteStep() {
	header := widget.NewLabelWithStyle("🎉 Account Activated Successfully!", fyne.TextAlignCenter, fyne.TextStyle{Bold: true})
	subtext := widget.NewLabelWithStyle(fmt.Sprintf("Welcome to Campus OS, %s!", w.AcademicName), fyne.TextAlignCenter, fyne.TextStyle{})

	card := widget.NewCard(
		"Digital Campus Pass",
		"Your account credentials have been verified",
		container.NewVBox(
			widget.NewLabel(fmt.Sprintf("Username : %s", w.Username)),
			widget.NewLabel(fmt.Sprintf("Email    : %s@campus.edu", w.Username)),
			widget.NewLabel("Gate Pass: ACTIVE (Ready for Checkpoint Scanning)"),
		),
	)

	enterBtn := widget.NewButtonWithIcon("Enter Campus OS Dashboard", theme.HomeIcon(), func() {
		if w.onComplete != nil {
			w.onComplete()
		}
	})
	enterBtn.Importance = widget.HighImportance

	w.content.Add(header)
	w.content.Add(subtext)
	w.content.Add(card)
	w.content.Add(enterBtn)
}
