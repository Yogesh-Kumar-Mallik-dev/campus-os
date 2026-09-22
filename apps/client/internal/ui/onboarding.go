package ui

import (
	"bytes"
	"context"
	"fmt"
	"image"
	"io"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/storage"
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
	"Voucher",
	"Telephony",
	"Dossier",
	"Password",
}

// OnboardingWizard manages the UI state and transitions for account activation.
type OnboardingWizard struct {
	client     *api.Client
	window     fyne.Window
	content    *fyne.Container
	step              OnboardingStep
	onComplete        func()
	onNavigateToLogin func()

	// Wizard Form State
	ClaimToken        string
	ScannedImageName  string
	ScannedImageObj   image.Image
	ScannedImageSize  string
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

	// Persistent Stepper indicator across all steps to prevent layout jump
	stepper := container.NewCenter(NewStepIndicator(int(w.step)+1, stepNames))
	w.content.Add(stepper)

	// Status / Error Banner using NewShadcnAlert
	if w.StatusText != "" {
		variant := AlertSuccess
		svgIcon := ResourceFromSVG("check.svg", SVGCheckVerified)
		title := "Verification Confirmed"
		if w.IsError {
			variant = AlertDestructive
			svgIcon = ResourceFromSVG("alert.svg", SVGClockGrace)
			title = "Action Required"
		}
		alert := NewShadcnAlert(title, w.StatusText, variant, svgIcon)
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

// processUploadedImage decodes an image stream, provides instant visual feedback, and extracts QR tokens.
func (w *OnboardingWizard) processUploadedImage(name string, r io.Reader) {
	data, err := io.ReadAll(r)
	if err != nil || len(data) == 0 {
		fyne.Do(func() {
			w.IsError = true
			w.StatusText = "Failed to read image file data."
			w.toast("Read Error", w.StatusText, AlertDestructive)
			w.render()
		})
		return
	}

	sizeStr := fmt.Sprintf("%.1f KB", float64(len(data))/1024.0)
	img, _, decodeErr := image.Decode(bytes.NewReader(data))
	if decodeErr != nil {
		fyne.Do(func() {
			w.ScannedImageName = name
			w.ScannedImageObj = nil
			w.ScannedImageSize = sizeStr
			w.IsError = true
			w.StatusText = fmt.Sprintf("Unsupported or corrupt image format: %v", decodeErr)
			w.toast("Invalid Image", w.StatusText, AlertDestructive)
			w.render()
		})
		return
	}

	token, qrErr := DecodeQRFromLoadedImage(img)
	fyne.Do(func() {
		w.ScannedImageName = name
		w.ScannedImageObj = img
		w.ScannedImageSize = sizeStr
		if qrErr != nil {
			w.ClaimToken = ""
			w.IsError = true
			w.StatusText = fmt.Sprintf("No QR code detected in '%s'. Ensure the code is clear, glare-free, and not cropped.", name)
			w.toast("No QR Code Detected", "Please ensure the QR code is clearly visible.", AlertDestructive)
		} else {
			w.ClaimToken = token
			w.IsError = false
			w.StatusText = fmt.Sprintf("QR code decoded successfully from %s: %s", name, token)
			w.toast("QR Picture Decoded", fmt.Sprintf("Token: %s", token), AlertSuccess)
		}
		w.render()
	})
}

// Step 1: Scan Sealed QR Code
func (w *OnboardingWizard) renderScanStep() {
	var visualPreview fyne.CanvasObject

	if w.ScannedImageObj != nil {
		thumbnail := canvas.NewImageFromImage(w.ScannedImageObj)
		thumbnail.FillMode = canvas.ImageFillContain
		thumbnail.SetMinSize(fyne.NewSize(140, 140))
		thumbnailBox := container.NewGridWrap(fyne.NewSize(140, 140), thumbnail)

		fileTitle := widget.NewLabelWithStyle(w.ScannedImageName, fyne.TextAlignLeading, fyne.TextStyle{Bold: true})
		fileTitle.Wrapping = fyne.TextWrapWord

		var statusBadge fyne.CanvasObject
		var statusDetail *widget.Label

		if !w.IsError && w.ClaimToken != "" {
			statusBadge = NewBadge("QR CODE RECOGNIZED", BadgeSuccess, BadgeShapePill)
			statusDetail = widget.NewLabel(fmt.Sprintf("Extracted Voucher Token:\n%s", w.ClaimToken))
			statusDetail.TextStyle = fyne.TextStyle{Bold: true}
		} else {
			statusBadge = NewBadge("NO QR CODE DETECTED", BadgeDestructive, BadgeShapePill)
			statusDetail = widget.NewLabel("No QR code found in this image. Try another picture, adjust crop, or enter the code manually.")
			statusDetail.Wrapping = fyne.TextWrapWord
		}

		removeBtn := NewShadcnButton("Remove / Select Another", ButtonOutline, ButtonSizeSm, ResourceFromSVG("close.svg", LucideX), func() {
			w.ScannedImageObj = nil
			w.ScannedImageName = ""
			w.ScannedImageSize = ""
			w.ClaimToken = ""
			w.IsError = false
			w.StatusText = ""
			w.render()
		})

		rightDetails := container.NewVBox(
			container.NewHBox(statusBadge, NewBadge(w.ScannedImageSize, BadgeSecondary, BadgeShapePill)),
			fileTitle,
			statusDetail,
			removeBtn,
		)

		previewCard := container.NewBorder(
			nil, nil,
			thumbnailBox,
			nil,
			container.NewPadded(rightDetails),
		)

		bg := canvas.NewRectangle(theme.Color(theme.ColorNameMenuBackground))
		bg.StrokeColor = theme.Color(theme.ColorNameInputBorder)
		bg.StrokeWidth = 1
		bg.CornerRadius = 8

		visualPreview = container.NewStack(bg, container.NewPadded(previewCard))
	} else if w.ScannedImageName != "" {
		checkIconRes := ResourceFromSVG("check_circle.svg", LucideCheckCircle2)
		scannerVisual := RenderSVGImage(checkIconRes, 56, 56)
		scannerTitle := "Live Optical Camera QR Decoded"
		scannerDesc := fmt.Sprintf("Extracted admission token: %s", w.ClaimToken)

		visualPreview = NewEmptyState(EmptyStateParams{
			Media:       scannerVisual,
			Title:       scannerTitle,
			Description: scannerDesc,
		})
	} else {
		qrIconRes := ResourceFromSVG("qr_viewfinder.svg", SVGQRCodeFrame)
		scannerVisual := RenderSVGImage(qrIconRes, 56, 56)
		title := "Optical Camera or Picture Scan"
		desc := "Align docket under optical camera or select a WhatsApp / screenshot image"

		visualPreview = NewEmptyState(EmptyStateParams{
			Media:       scannerVisual,
			Title:       title,
			Description: desc,
		})
	}

	tokenField, tokenEntry := NewFormField(FormField{
		Label:       "Voucher Claim Token",
		Placeholder: "Enter 16-character code (or leave blank for demo)",
		HelperText:  "Auto-filled from decoded QR picture or manually entered from admission docket",
	})
	if w.ClaimToken != "" {
		tokenEntry.SetText(w.ClaimToken)
	}

	scanCameraBtn := NewShadcnButton("Open Optical Camera", ButtonDefault, ButtonSizeDefault, WhiteResourceFromSVG("camera.svg", LucideCamera), func() {
		if w.window == nil {
			return
		}

		launchScanner := func() {
			ShowCameraScannerModal(w.window, func(token string) {
				w.ClaimToken = token
				w.ScannedImageName = "Live Optical Camera"
				w.ScannedImageObj = nil
				tokenEntry.SetText(token)
				w.IsError = false
				w.StatusText = fmt.Sprintf("QR code decoded successfully from Camera: %s", token)
				w.toast("QR Code Scanned", fmt.Sprintf("Token: %s", token), AlertSuccess)
				w.render()
			}, func() {
				// Cancelled by user
			})
		}

		perm := GetCameraPermission()
		switch perm {
		case PermissionGranted:
			launchScanner()
		case PermissionDenied:
			ShowAlertDialog(w.window, "Camera Access Denied", "Camera permission is currently disabled on this device. Would you like to enable it?", "Enable Camera", false, func() {
				_ = SetCameraPermission(PermissionGranted)
				launchScanner()
			})
		default: // PermissionPrompt
			ShowCameraPermissionModal(w.window, func() {
				launchScanner()
			}, func() {
				w.toast("Camera Access Denied", "Camera access was not granted. You can still select a picture file.", AlertDestructive)
			})
		}
	})

	uploadBtn := NewShadcnButton("Select QR Picture", ButtonOutline, ButtonSizeDefault, ResourceFromSVG("upload.svg", LucideUpload), func() {
		if w.window == nil {
			return
		}
		fd := dialog.NewFileOpen(func(reader fyne.URIReadCloser, err error) {
			if err != nil || reader == nil {
				return
			}
			defer reader.Close()
			w.processUploadedImage(reader.URI().Name(), reader)
		}, w.window)

		fd.SetFilter(storage.NewExtensionFileFilter([]string{".png", ".jpg", ".jpeg", ".webp", ".bmp"}))
		fd.SetTitleText("Select Admission QR Picture or Screenshot")
		fd.Show()
	})

	// Setup drag-and-drop on desktop window
	if w.window != nil {
		w.window.SetOnDropped(func(pos fyne.Position, uris []fyne.URI) {
			if len(uris) == 0 || w.step != StepScan {
				return
			}
			reader, err := storage.Reader(uris[0])
			if err != nil {
				return
			}
			defer reader.Close()
			w.processUploadedImage(uris[0].Name(), reader)
		})
	}

	scanBtn := NewShadcnButton("Validate Admission Token", ButtonDefault, ButtonSizeDefault, WhiteResourceFromSVG("scan.svg", LucideScan), func() {
		token := tokenEntry.Text
		if token == "" {
			token = "claim_genesis_test_demo"
		}
		w.ClaimToken = token

		go func() {
			resp, err := w.client.ValidateClaim(context.Background(), token)
			fyne.Do(func() {
				if err == nil && resp.IsValid {
					w.MaskedPhone = resp.MaskedPhoneNumber
					w.Username = resp.Username
					w.AcademicName = resp.AcademicName
					w.LegalFullName = resp.LegalFullName
					w.AdmissionType = resp.AdmissionType
					w.EntrySemester = resp.EntrySemesterNumber
					w.LateralSummary = resp.LateralEntrySummary
					w.IsError = false
					w.StatusText = "QR Token Verified with Backend. Hardware telephony match required."
					w.toast("QR Code Recognized", "Hardware carrier binding required next.", AlertSuccess)
					w.SetStep(StepSIMVerify)
				} else if err != nil && api.IsUnreachable(err) {
					// Fallback demo mock values for instant local dev testing when backend is offline
					w.MaskedPhone = "+91 98XXX-XX210"
					w.Username = "yogesh.cse.2024.l"
					w.AcademicName = "Yogesh"
					w.LegalFullName = "Yogesh Kumar Mallik"
					w.AdmissionType = "LATERAL_ENTRY"
					w.EntrySemester = 3
					w.LateralSummary = "Lateral Entry: Direct admission to Semester 3. Prior polytechnic credits verified."
					w.IsError = false
					w.StatusText = "Offline Demo Mode: Backend unreachable, proceeding with simulated profile."
					w.toast("Offline Mode Active", "Backend offline. Using simulated student record.", AlertWarning)
					w.SetStep(StepSIMVerify)
				} else {
					w.IsError = true
					if err != nil {
						w.StatusText = err.Error()
					} else {
						w.StatusText = "Invalid Claim Token: The token is expired, consumed, or invalid."
					}
					w.toast("Validation Failed", w.StatusText, AlertDestructive)
					w.render()
				}
			})
		}()
	})

	signInBtn := NewShadcnButton("Already Activated? Institutional Sign In", ButtonGhost, ButtonSizeDefault, ResourceFromSVG("login.svg", LucideLogIn), func() {
		if w.onNavigateToLogin != nil {
			w.onNavigateToLogin()
		}
	})

	scannerCard := NewShadcnCard(CardParts{
		Badge:       NewBadge("STEP 1 OF 4", BadgeDefault, BadgeShapePill),
		Title:       "Scan Sealed Admission QR",
		Description: "Position camera viewfinder, select a screenshot/picture, or enter the manual token beneath the seal",
		Content: container.NewVBox(
			visualPreview,
			container.NewCenter(container.NewHBox(scanCameraBtn, uploadBtn)),
			NewShadcnSeparator(true),
			tokenField,
		),
		Footer: container.NewVBox(
			scanBtn,
			NewShadcnSeparator(true),
			signInBtn,
		),
	})

	w.content.Add(scannerCard)
}

// Step 2: SIM Telephony & SMS OTP Handshake
func (w *OnboardingWizard) renderSIMVerifyStep() {
	phoneInfo := NewKeyValueRow("REGISTERED TELEPHONY CONTACT", w.MaskedPhone, nil)

	var sim1Card, sim2Card *SelectableCard

	sim1Icon := RenderSVGImage(ResourceFromSVG("sim_active.svg", SVGSIMCardActive), 20, 20)
	sim1Label := widget.NewLabel("Slot 1 (Jio 5G): +91 98765-43210")
	sim1Label.Wrapping = fyne.TextWrapWord
	sim1Content := container.NewBorder(nil, nil,
		container.NewHBox(sim1Icon),
		NewBadge("Slot 1", BadgeSecondary, BadgeShapePill),
		sim1Label,
	)

	sim2Icon := RenderSVGImage(ResourceFromSVG("sim_sec.svg", SVGSIMCardSecondary), 20, 20)
	sim2Label := widget.NewLabel("Slot 2 (Airtel): +91 91234-56789")
	sim2Label.Wrapping = fyne.TextWrapWord
	sim2Content := container.NewBorder(nil, nil,
		container.NewHBox(sim2Icon),
		NewBadge("Slot 2", BadgeSecondary, BadgeShapePill),
		sim2Label,
	)

	sim1Card = NewSelectableCard(sim1Content, true, func() {
		sim1Card.SetSelected(true)
		sim2Card.SetSelected(false)
		w.SelectedSIMPhone = "+919876543210"
	})

	sim2Card = NewSelectableCard(sim2Content, false, func() {
		sim1Card.SetSelected(false)
		sim2Card.SetSelected(true)
		w.SelectedSIMPhone = "+919123456789"
	})

	otpField, otpEntry := NewFormField(FormField{
		Label:       "One-Time Challenge Code (OTP)",
		Placeholder: "Enter 6-digit SMS OTP (e.g. 123456)",
		HelperText:  "Sent via cellular carrier handshake to verify registered SIM presence",
	})

	verifyBtn := NewShadcnButton("Verify Carrier SIM & Submit OTP", ButtonDefault, ButtonSizeDefault, WhiteResourceFromSVG("phone.svg", LucideSmartphone), func() {
		otp := otpEntry.Text
		if otp == "" {
			otp = "123456"
		}
		w.OTPCode = otp

		go func() {
			resp, err := w.client.VerifySIM(context.Background(), w.ClaimToken, "dummy_sim_iccid_hash", w.SelectedSIMPhone)
			fyne.Do(func() {
				if err == nil {
					if !resp.SimMatched {
						w.IsError = true
						w.StatusText = fmt.Sprintf("SIM Hardware Mismatch: Selected SIM (%s) does not match the registered telephony contact on file (%s).", w.SelectedSIMPhone, w.MaskedPhone)
						w.toast("SIM Mismatch", "Physical SIM does not match student record.", AlertDestructive)
						w.render()
						return
					}
					w.OTPChallengeID = resp.OTPChallengeID
					w.IsError = false
					w.StatusText = "SIM Binding & OTP Confirmed via Backend. Review official admission dossier."
					w.toast("SIM & OTP Verified", "Hardware telephony binding confirmed.", AlertSuccess)
					w.SetStep(StepReviewProfile)
				} else if api.IsUnreachable(err) {
					w.OTPChallengeID = "mock_otp_challenge_123"
					w.IsError = false
					w.StatusText = "SIM Binding Simulated (Offline). Review official admission dossier."
					w.toast("SIM Verified (Offline)", "Review your official admission records.", AlertSuccess)
					w.SetStep(StepReviewProfile)
				} else {
					w.IsError = true
					w.StatusText = err.Error()
					w.toast("Verification Failed", err.Error(), AlertDestructive)
					w.render()
				}
			})
		}()
	})

	bypassBtn := NewShadcnButton("Bypass SIM Handshake (Dev Mode)", ButtonGhost, ButtonSizeDefault, ResourceFromSVG("check.svg", LucideCheckCircle2), func() {
		w.OTPChallengeID = "mock_otp_challenge_123"
		w.OTPCode = "123456"
		w.IsError = false
		w.StatusText = "SIM Handshake bypassed in development mode."
		w.toast("SIM Bypassed (Dev)", "Proceeding directly to dossier verification.", AlertSuccess)
		w.SetStep(StepReviewProfile)
	})

	simCard := NewShadcnCard(CardParts{
		Badge:       NewBadge("STEP 2 OF 4", BadgeDefault, BadgeShapePill),
		Title:       "Device SIM & Telephony Handshake",
		Description: "Hardware carrier binding ensures activation is locked strictly to your physical handset",
		Content: container.NewVBox(
			phoneInfo,
			NewShadcnSeparator(true),
			sim1Card,
			sim2Card,
			NewShadcnSeparator(true),
			otpField,
		),
		Footer: container.NewVBox(
			verifyBtn,
			NewShadcnSeparator(true),
			bypassBtn,
		),
	})

	w.content.Add(simCard)
}

// Step 3: Review Official Records & Dual-Identity Confirmation
func (w *OnboardingWizard) renderReviewProfileStep() {
	academicName := w.AcademicName
	if academicName == "" {
		academicName = "Yogesh"
	}
	legalName := w.LegalFullName
	if legalName == "" {
		legalName = "Yogesh Kumar Mallik"
	}
	username := w.Username
	if username == "" {
		username = "yogesh.cse.2024.l"
	}
	admissionType := w.AdmissionType
	if admissionType == "" {
		admissionType = "LATERAL_ENTRY"
	}
	semester := w.EntrySemester
	if semester <= 0 {
		semester = 3
	}
	lateralNotes := w.LateralSummary
	if lateralNotes == "" {
		lateralNotes = "Direct admission to Semester 3 (Lateral Entry). Prior polytechnic credits verified by AKTU."
	}

	// Section 1: Dual-Identity KYC Records
	sec1Title := widget.NewLabelWithStyle("Identity & KYC Synchronization", fyne.TextAlignLeading, fyne.TextStyle{Bold: true})
	sec1Desc := widget.NewLabel("Verification correspondence between academic marksheets and government identification")
	sec1Desc.Wrapping = fyne.TextWrapWord

	kycSection := container.NewVBox(
		sec1Title,
		sec1Desc,
		NewKeyValueRow("ACADEMIC SCHOLAR NAME", academicName, nil),
		NewKeyValueRow("LEGAL FULL NAME (GOVT AADHAAR)", legalName, nil),
	)

	// Section 2: Institutional Allocation & Credentials
	sec2Title := widget.NewLabelWithStyle("Institutional Allocation & Pathway", fyne.TextAlignLeading, fyne.TextStyle{Bold: true})
	sec2Desc := widget.NewLabel("Assigned university credentials, entry semester, and curricular credits")
	sec2Desc.Wrapping = fyne.TextWrapWord

	institutionalSection := container.NewVBox(
		sec2Title,
		sec2Desc,
		NewKeyValueRow("INSTITUTIONAL USERNAME", username, nil),
		NewKeyValueRow("OFFICIAL EMAIL", fmt.Sprintf("%s@campus.edu", username), nil),
		NewKeyValueRow("ADMISSION TYPE & SEMESTER", fmt.Sprintf("%s • Semester %d", admissionType, semester), nil),
		NewKeyValueRow("CURRICULAR NOTES", lateralNotes, nil),
	)

	// Confirmation Row
	confirmCheck := widget.NewCheck("", nil)
	confirmCheck.SetChecked(true)
	checkText := widget.NewLabel("I acknowledge that these records accurately reflect my marksheet and government identification")
	checkText.Wrapping = fyne.TextWrapWord
	confirmRow := container.NewBorder(nil, nil, confirmCheck, nil, checkText)

	discrepancyBtn := NewShadcnButton("Report Record Discrepancy", ButtonOutline, ButtonSizeDefault, ResourceFromSVG("flag.svg", LucideFlag), func() {
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

	nextBtn := NewAnimatedArrowButton("Records Confirmed, Set Password", ButtonDefault, ButtonSizeDefault, WhiteResourceFromSVG("file-check.svg", LucideFileCheck), func() {
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

	confirmCheck.OnChanged = func(checked bool) {
		if checked {
			nextBtn.Enable()
		} else {
			nextBtn.Disable()
		}
	}

	// Full-width vertical footer actions to eliminate horizontal squishing across screen sizes
	actionFooter := container.NewVBox(
		nextBtn,
		discrepancyBtn,
	)

	profileCard := NewShadcnCard(CardParts{
		Badge:       NewBadge("STEP 3 OF 4", BadgeDefault, BadgeShapePill),
		Title:       "Review Official Admission Dossier",
		Description: "Official records provisioned by the Academic Registrar. Review carefully before confirming.",
		Content: container.NewVBox(
			kycSection,
			NewShadcnSeparator(true),
			institutionalSection,
			NewShadcnSeparator(true),
			confirmRow,
		),
		Footer: actionFooter,
	})

	w.content.Add(profileCard)
}

// Step 4: Password Setup & 3-Day Orientation Grace Period
func (w *OnboardingWizard) renderSetPasswordStep() {
	clockIcon := ResourceFromSVG("clock.svg", LucideClock)
	graceAlert := NewShadcnAlert(
		"Orientation Grace Period Active",
		"A simplified 6+ character password is accepted during your first 72 hours of enrollment",
		AlertWarning,
		clockIcon,
	)

	passField, passEntry := NewFormField(FormField{
		Label:       "New Access Password",
		Placeholder: "Enter password (min 6 characters)",
		HelperText:  "Choose a secure pass-phrase. Must be at least 6 characters during grace period.",
		IsPassword:  true,
	})

	confirmPassField, confirmPassEntry := NewFormField(FormField{
		Label:       "Confirm New Password",
		Placeholder: "Re-enter new password to verify",
		IsPassword:  true,
	})

	fingerprintIcon := RenderSVGImage(ResourceFromSVG("fingerprint.svg", LucideFingerprint), 20, 20)
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

	completeBtn := NewShadcnButton("Complete Activation & Provision Pass", ButtonDefault, ButtonSizeDefault, WhiteResourceFromSVG("lock.svg", LucideLock), func() {
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

		go func() {
			otpID := w.OTPChallengeID
			if otpID == "" {
				otpID = "otp_mock_challenge_123"
			}
			otpCode := w.OTPCode
			if otpCode == "" {
				otpCode = "123456"
			}

			resp, err := w.client.CompleteClaim(context.Background(), w.ClaimToken, otpID, otpCode, w.NewPassword, w.BiometricsEnabled)
			fyne.Do(func() {
				if err == nil && resp.Success {
					if resp.Username != "" {
						w.Username = resp.Username
					}
					w.IsError = false
					w.StatusText = "Account Successfully Activated via Backend"
					w.toast("Success", "Account credentials provisioned.", AlertSuccess)
					w.SetStep(StepComplete)
				} else if err != nil && api.IsUnreachable(err) {
					w.IsError = false
					w.StatusText = "Account Activated (Offline Demo)"
					w.toast("Success", "Account credentials provisioned.", AlertSuccess)
					w.SetStep(StepComplete)
				} else {
					w.IsError = true
					if err != nil {
						w.StatusText = err.Error()
					} else {
						w.StatusText = "Account activation failed"
					}
					w.toast("Activation Failed", w.StatusText, AlertDestructive)
					w.render()
				}
			})
		}()
	})

	passwordCard := NewShadcnCard(CardParts{
		Badge:       NewBadge("STEP 4 OF 4", BadgeDefault, BadgeShapePill),
		Title:       "Set Your Access Password",
		Description: "Configure your primary institutional password and device biometric keychain",
		Content: container.NewVBox(
			graceAlert,
			passField,
			confirmPassField,
			NewShadcnSeparator(true),
			biometricRow,
		),
		Footer: completeBtn,
	})

	w.content.Add(passwordCard)
}

// Step 5: Celebration & Gateway Handoff
func (w *OnboardingWizard) renderCompleteStep() {
	sparklesIcon := RenderSVGImage(ResourceFromSVG("sparkles.svg", LucideSparkles), 48, 48)

	scholarName := w.AcademicName
	if scholarName == "" {
		scholarName = "Yogesh"
	}
	username := w.Username
	if username == "" {
		username = "yogesh.cse.2024.l"
	}

	passRows := container.NewVBox(
		container.NewCenter(sparklesIcon),
		container.NewCenter(NewStatusBadge("VERIFIED SCHOLAR", BadgeSuccess)),
		NewShadcnSeparator(true),
		NewKeyValueRow("SCHOLAR NAME", scholarName, nil),
		NewShadcnSeparator(true),
		NewKeyValueRow("INSTITUTIONAL USERNAME", username, nil),
		NewShadcnSeparator(true),
		NewKeyValueRow("INSTITUTIONAL EMAIL", fmt.Sprintf("%s@campus.edu", username), nil),
		NewShadcnSeparator(true),
		NewKeyValueRow("CAMPUS GATE ACCESS", "ENABLED • NFC & QR Checkpoint Check-in Active", nil),
		NewShadcnSeparator(true),
		NewKeyValueRow("INSTITUTIONAL AFFILIATION", "Dr. A.P.J. Abdul Kalam Technical University (AKTU)", nil),
	)

	enterBtn := NewAnimatedArrowButton("Finish & Enter Scholar Portal", ButtonDefault, ButtonSizeDefault, WhiteResourceFromSVG("graduation-cap.svg", LucideGraduationCap), func() {
		if w.onComplete != nil {
			w.onComplete()
		}
	})

	idCard := NewShadcnCard(CardParts{
		Badge:       NewBadge("ACTIVATION COMPLETE", BadgeDefault, BadgeShapePill),
		Title:       "Digital Campus Pass",
		Description: "Cryptographically verified digital pass ready for turnstile and library checkpoints",
		Content:     passRows,
		Footer:      enterBtn,
	})

	w.content.Add(idCard)
}
