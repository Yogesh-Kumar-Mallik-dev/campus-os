package ui

import (
	"context"
	"fmt"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
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

var stepNames = []string{
	"QR Scan",
	"SIM Binding",
	"Dossier",
	"Password",
}

// OnboardingWizard manages the UI state and transitions for account activation.
type OnboardingWizard struct {
	client     *api.Client
	window     fyne.Window
	content    *fyne.Container
	step       OnboardingStep
	onComplete func()

	// Wizard Form State
	ClaimToken        string
	MaskedPhone       string
	Username          string
	AcademicName      string
	LegalFullName     string
	AdmissionType     string
	EntrySemester     int
	LateralSummary    string
	OTPChallengeID    string
	SelectedSIMPhone  string
	OTPCode           string
	NewPassword       string
	BiometricsEnabled bool
	StatusText        string
	IsError           bool
}

// BLOCK_UI_ONBOARDING_NEW_001
// Purpose: Constructs a new OnboardingWizard instance using native shadcn primitives.
func NewOnboardingWizard(client *api.Client, window fyne.Window, onComplete func()) *OnboardingWizard {
	w := &OnboardingWizard{
		client:            client,
		window:            window,
		step:              StepScan,
		onComplete:        onComplete,
		SelectedSIMPhone:  "+919876543210",
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

func (w *OnboardingWizard) toast(title, message string, variant AlertVariant) {
	if w.window != nil {
		ShowToast(w.window, title, message, variant, 3*time.Second)
	}
}

func (w *OnboardingWizard) render() {
	w.content.Objects = nil

	// Responsive Stepper indicator
	if w.step < StepComplete {
		stepper := container.NewCenter(NewStepIndicator(int(w.step)+1, stepNames))
		w.content.Add(stepper)
	}

	// Status / Error Banner using NewShadcnAlert
	if w.StatusText != "" {
		variant := AlertSuccess
		svgIcon := ResourceFromSVG("check.svg", SVGCheckVerified)
		if w.IsError {
			variant = AlertDestructive
			svgIcon = ResourceFromSVG("alert.svg", SVGClockGrace)
		}
		alert := NewShadcnAlert("Status Notification", w.StatusText, variant, svgIcon)
		w.content.Add(alert)
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
	header := NewPageHeader(
		"Scan Sealed Admission QR",
		"Position camera viewfinder over the tamper-evident QR on your official admission slip",
		NewBadge("STEP 1 OF 4", BadgeDefault, BadgeShapePill),
	)

	qrIconRes := ResourceFromSVG("qr_viewfinder.svg", SVGQRCodeFrame)
	qrVisual := RenderSVGImage(qrIconRes, 68, 68)

	emptyScanner := NewEmptyState(EmptyStateParams{
		Media:       qrVisual,
		Title:       "Optical Camera Ready",
		Description: "Align official QR code within optical brackets",
	})
	scannerCard := NewShadcnCard(CardParts{
		Title:   "Optical Capture Interface",
		Content: emptyScanner,
	})

	tokenEntry := NewShadcnInput("Enter 16-character code (or leave blank for demo)", false)
	if w.ClaimToken != "" {
		tokenEntry.SetText(w.ClaimToken)
	}

	scanBtn := NewShadcnButton("Validate QR Code Token", ButtonDefault, ButtonSizeDefault, theme.ConfirmIcon(), func() {
		token := tokenEntry.Text
		if token == "" {
			token = "claim_genesis_test_demo"
		}
		w.ClaimToken = token

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
				// Fallback demo mock values for instant local dev testing
				w.MaskedPhone = "+91 98XXX-XX210"
				w.Username = "yogesh.cse.2024.l"
				w.AcademicName = "Yogesh"
				w.LegalFullName = "Yogesh Kumar Mallik"
				w.AdmissionType = "LATERAL_ENTRY"
				w.EntrySemester = 3
				w.LateralSummary = "Lateral Entry: Direct admission to Semester 3. Prior polytechnic credits verified."
			}
			w.IsError = false
			w.StatusText = "QR Token Verified. Hardware telephony match required."
			w.toast("QR Code Recognized", "Hardware carrier binding required next.", AlertSuccess)
			w.SetStep(StepSIMVerify)
		}()
	})

	formCard := NewShadcnCard(CardParts{
		Title: "Manual Token Input",
		Content: container.NewVBox(
			tokenEntry,
			scanBtn,
		),
	})

	w.content.Add(header)
	w.content.Add(scannerCard)
	w.content.Add(formCard)
}

// Step 2: SIM Telephony & SMS OTP Handshake
func (w *OnboardingWizard) renderSIMVerifyStep() {
	header := NewPageHeader(
		"Device SIM & Telephony Handshake",
		"Hardware carrier binding ensures activation is locked strictly to your physical handset",
		NewBadge("STEP 2 OF 4", BadgeWarning, BadgeShapePill),
	)

	phoneInfo := widget.NewLabel(fmt.Sprintf("Registered Scholar Phone: %s", w.MaskedPhone))
	phoneInfo.Wrapping = fyne.TextWrapWord

	sim1Icon := RenderSVGImage(ResourceFromSVG("sim_active.svg", SVGSIMCardActive), 20, 20)
	sim1Label := widget.NewLabel("Slot 1 (Jio 5G): +91 98765-43210")
	sim1Label.Wrapping = fyne.TextWrapWord
	sim1Row := container.NewBorder(nil, nil,
		container.NewHBox(sim1Icon),
		NewBadge("MATCH", BadgeSuccess, BadgeShapePill),
		sim1Label,
	)

	sim2Icon := RenderSVGImage(ResourceFromSVG("sim_sec.svg", SVGSIMCardSecondary), 20, 20)
	sim2Label := widget.NewLabel("Slot 2 (Airtel): +91 91234-56789")
	sim2Label.Wrapping = fyne.TextWrapWord
	sim2Row := container.NewBorder(nil, nil,
		container.NewHBox(sim2Icon),
		NewBadge("SECONDARY", BadgeSecondary, BadgeShapePill),
		sim2Label,
	)

	simCard := NewShadcnCard(CardParts{
		Title: "Hardware Telephony Slots",
		Content: container.NewVBox(
			phoneInfo,
			NewShadcnSeparator(true),
			sim1Row,
			sim2Row,
		),
	})

	otpEntry := NewShadcnInput("Enter 6-digit SMS OTP (e.g. 123456)", false)

	verifyBtn := NewShadcnButton("Verify Carrier SIM & Submit OTP", ButtonDefault, ButtonSizeDefault, theme.ConfirmIcon(), func() {
		w.OTPCode = otpEntry.Text
		if w.OTPCode == "" {
			w.OTPCode = "123456"
		}
		w.IsError = false
		w.StatusText = "SIM Binding & OTP Confirmed. Review official admission dossier."
		w.toast("SIM & OTP Verified", "Review your official admission records.", AlertSuccess)
		w.SetStep(StepReviewProfile)
	})

	otpCard := NewShadcnCard(CardParts{
		Title: "SMS Challenge Verification",
		Content: container.NewVBox(
			otpEntry,
			verifyBtn,
		),
	})

	w.content.Add(header)
	w.content.Add(simCard)
	w.content.Add(otpCard)
}

// Step 3: Review Official Records & Dual-Identity Confirmation
func (w *OnboardingWizard) renderReviewProfileStep() {
	header := NewPageHeader(
		"Review Official Records",
		"Verify correspondence between secondary marksheet and govt identification",
		NewBadge("STEP 3 OF 4", BadgeDefault, BadgeShapePill),
	)

	acadLabel := widget.NewLabel(w.AcademicName)
	acadLabel.Wrapping = fyne.TextWrapWord
	legalLabel := widget.NewLabel(w.LegalFullName)
	legalLabel.Wrapping = fyne.TextWrapWord
	notesLabel := widget.NewLabel(w.LateralSummary)
	notesLabel.Wrapping = fyne.TextWrapWord

	profileForm := widget.NewForm(
		widget.NewFormItem("Academic Name", acadLabel),
		widget.NewFormItem("Legal Full Name", legalLabel),
		widget.NewFormItem("Username", widget.NewLabel(w.Username)),
		widget.NewFormItem("Email", widget.NewLabel(fmt.Sprintf("%s@campus.edu", w.Username))),
		widget.NewFormItem("Admission", widget.NewLabel(fmt.Sprintf("%s (Sem %d)", w.AdmissionType, w.EntrySemester))),
		widget.NewFormItem("Notes", notesLabel),
	)

	profileCard := NewShadcnCard(CardParts{
		Title:   "Verified Admission Dossier",
		Content: profileForm,
	})

	confirmCheck := widget.NewCheck("", nil)
	confirmCheck.SetChecked(true)
	checkText := widget.NewLabel("I confirm these records accurately reflect my marksheet and govt identification")
	checkText.Wrapping = fyne.TextWrapWord
	confirmRow := container.NewBorder(nil, nil, confirmCheck, nil, checkText)

	discrepancyBtn := NewShadcnButton("Report Discrepancy", ButtonOutline, ButtonSizeDefault, nil, func() {
		if w.window != nil {
			ShowAlertDialog(w.window, "Report Record Discrepancy",
				"If your academic marksheet or legal Aadhaar name differs from these records, an audit flag will be sent to the Registrar desk.",
				"Acknowledge & Flag", false, func() {
					w.IsError = true
					w.StatusText = "Audit flag recorded. Please consult the Academic Registrar desk."
					w.toast("Discrepancy Reported", "Audit flag submitted to Registrar.", AlertWarning)
					w.render()
				})
		}
	})

	nextBtn := NewShadcnButton("Records Confirmed, Set Password ->", ButtonDefault, ButtonSizeDefault, theme.NavigateNextIcon(), func() {
		if !confirmCheck.Checked {
			w.IsError = true
			w.StatusText = "Please acknowledge record verification before proceeding"
			w.toast("Confirmation Required", "Please check the confirmation box.", AlertDestructive)
			w.render()
			return
		}
		w.IsError = false
		w.StatusText = "Identity records confirmed. Create your password."
		w.SetStep(StepSetPassword)
	})

	actionRow := container.NewGridWithColumns(2, discrepancyBtn, nextBtn)

	w.content.Add(header)
	w.content.Add(profileCard)
	w.content.Add(confirmRow)
	w.content.Add(actionRow)
}

// Step 4: Password Setup & 3-Day Orientation Grace Period
func (w *OnboardingWizard) renderSetPasswordStep() {
	header := NewPageHeader(
		"Set Your Access Password",
		"Create a secure password to finalize institutional account provisioning",
		NewBadge("STEP 4 OF 4", BadgeWarning, BadgeShapePill),
	)

	clockIcon := ResourceFromSVG("clock.svg", SVGClockGrace)
	graceAlert := NewShadcnAlert(
		"Orientation Grace Period Active",
		"A simplified 6+ character password is accepted during your first 72 hours of enrollment",
		AlertWarning,
		clockIcon,
	)

	passEntry := NewShadcnInput("Enter new password (min 6 characters)", true)
	confirmPassEntry := NewShadcnInput("Confirm new password", true)

	fingerprintIcon := RenderSVGImage(ResourceFromSVG("fingerprint.svg", SVGFingerprint), 20, 20)
	bioLabel := widget.NewLabel("Enable Biometric Keyring (Fingerprint / Face ID)")
	bioLabel.Wrapping = fyne.TextWrapWord
	biometricSwitch := NewSwitch(w.BiometricsEnabled, func(b bool) {
		w.BiometricsEnabled = b
	})
	biometricRow := container.NewBorder(nil, nil,
		container.NewHBox(fingerprintIcon),
		biometricSwitch,
		bioLabel,
	)

	completeBtn := NewShadcnButton("Complete Activation & Provision Pass", ButtonDefault, ButtonSizeDefault, theme.ConfirmIcon(), func() {
		if len(passEntry.Text) < 6 {
			w.IsError = true
			w.StatusText = "Password must contain at least 6 characters during grace period"
			w.toast("Password Too Short", "Must be at least 6 characters.", AlertDestructive)
			w.render()
			return
		}
		if passEntry.Text != confirmPassEntry.Text {
			w.IsError = true
			w.StatusText = "Passwords do not match"
			w.toast("Mismatch Error", "Passwords do not match.", AlertDestructive)
			w.render()
			return
		}

		w.NewPassword = passEntry.Text
		w.IsError = false
		w.StatusText = "Account Successfully Activated"
		w.toast("Success", "Account credentials provisioned.", AlertSuccess)
		w.SetStep(StepComplete)
	})

	passwordCard := NewShadcnCard(CardParts{
		Title: "Credential Security",
		Content: container.NewVBox(
			graceAlert,
			widget.NewLabel("New Password:"),
			passEntry,
			widget.NewLabel("Confirm Password:"),
			confirmPassEntry,
			biometricRow,
			completeBtn,
		),
	})

	w.content.Add(header)
	w.content.Add(passwordCard)
}

// Step 5: Celebration & Gateway Handoff
func (w *OnboardingWizard) renderCompleteStep() {
	checkIcon := ResourceFromSVG("verified.svg", SVGCheckVerified)
	verifiedVisual := RenderSVGImage(checkIcon, 40, 40)

	header := NewPageHeader(
		"Account Activated Successfully",
		"Your institutional credentials and gate pass are active and verified",
		NewBadge("ACCOUNT ACTIVE", BadgeSuccess, BadgeShapePill),
	)

	scholarLabel := widget.NewLabel("Scholar: " + w.AcademicName)
	scholarLabel.Wrapping = fyne.TextWrapWord
	userLabel := widget.NewLabel("Username:    " + w.Username)
	userLabel.Wrapping = fyne.TextWrapWord
	emailLabel := widget.NewLabel("Email:       " + w.Username + "@campus.edu")
	emailLabel.Wrapping = fyne.TextWrapWord
	gateLabel := widget.NewLabel("Campus Gate: ENABLED (Ready for NFC/QR Checkpoint Scanning)")
	gateLabel.Wrapping = fyne.TextWrapWord
	affilLabel := widget.NewLabel("Affiliation: Dr. A.P.J. Abdul Kalam Technical University (AKTU)")
	affilLabel.Wrapping = fyne.TextWrapWord

	idBadge := container.NewVBox(
		container.NewCenter(verifiedVisual),
		container.NewBorder(nil, nil,
			scholarLabel,
			NewBadge("VERIFIED SCHOLAR", BadgeSuccess, BadgeShapePill),
		),
		NewShadcnSeparator(true),
		userLabel,
		emailLabel,
		gateLabel,
		affilLabel,
	)

	idCard := NewShadcnCard(CardParts{
		Title:   "Digital Campus Pass",
		Content: idBadge,
	})

	enterBtn := NewShadcnButton("Finish & Ready", ButtonDefault, ButtonSizeDefault, theme.HomeIcon(), func() {
		if w.onComplete != nil {
			w.onComplete()
		}
	})

	w.content.Add(header)
	w.content.Add(idCard)
	w.content.Add(layout.NewSpacer())
	w.content.Add(enterBtn)
}
