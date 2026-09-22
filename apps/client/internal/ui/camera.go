package ui

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"image"
	"image/color"
	"image/jpeg"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"sync"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

// CameraPermissionState defines user permission states for local webcam capture.
type CameraPermissionState string

const (
	PermissionPrompt  CameraPermissionState = "prompt"
	PermissionGranted CameraPermissionState = "granted"
	PermissionDenied  CameraPermissionState = "denied"
)

type preferencesData struct {
	CameraPermission CameraPermissionState `json:"camera_permission"`
}

var (
	prefMutex        sync.Mutex
	customPrefPath   string
	mockCameraRunner func(ctx context.Context, onFrame func(image.Image), onQRFound func(string)) error
)

func getPreferencesFilePath() string {
	if customPrefPath != "" {
		return customPrefPath
	}
	configDir, err := os.UserConfigDir()
	if err != nil {
		configDir = os.TempDir()
	}
	return filepath.Join(configDir, "campus-os", "preferences.json")
}

// GetCameraPermission loads the locally stored camera permission preference.
func GetCameraPermission() CameraPermissionState {
	prefMutex.Lock()
	defer prefMutex.Unlock()

	filePath := getPreferencesFilePath()
	data, err := os.ReadFile(filePath)
	if err != nil {
		return PermissionPrompt
	}

	var prefs preferencesData
	if err := json.Unmarshal(data, &prefs); err != nil || prefs.CameraPermission == "" {
		return PermissionPrompt
	}

	return prefs.CameraPermission
}

// SetCameraPermission updates and saves the camera permission preference locally.
func SetCameraPermission(state CameraPermissionState) error {
	prefMutex.Lock()
	defer prefMutex.Unlock()

	filePath := getPreferencesFilePath()
	if err := os.MkdirAll(filepath.Dir(filePath), 0700); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	var prefs preferencesData
	if data, err := os.ReadFile(filePath); err == nil {
		_ = json.Unmarshal(data, &prefs)
	}

	prefs.CameraPermission = state
	encoded, err := json.MarshalIndent(prefs, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to encode preferences: %w", err)
	}

	return os.WriteFile(filePath, encoded, 0600)
}

// SetCustomPreferencesPathForTest overrides preferences file path for testing.
func SetCustomPreferencesPathForTest(path string) {
	prefMutex.Lock()
	defer prefMutex.Unlock()
	customPrefPath = path
}

// SetMockCameraRunnerForTest overrides camera capture loop for automated testing.
func SetMockCameraRunnerForTest(runner func(ctx context.Context, onFrame func(image.Image), onQRFound func(string)) error) {
	prefMutex.Lock()
	defer prefMutex.Unlock()
	mockCameraRunner = runner
}

func getMockCameraRunner() func(ctx context.Context, onFrame func(image.Image), onQRFound func(string)) error {
	prefMutex.Lock()
	defer prefMutex.Unlock()
	return mockCameraRunner
}

// BLOCK_UI_CAMERA_MODAL_001
// Purpose: Renders a thematic modal asking for webcam permissions, in sync with Design System 2026.
func ShowCameraPermissionModal(window fyne.Window, onAllow func(), onDeny func()) *widget.PopUp {
	var popup *widget.PopUp

	cameraIcon := ResourceFromSVG("camera.svg", LucideCamera)
	iconVisual := RenderSVGImage(cameraIcon, 48, 48)

	titleLabel := widget.NewLabelWithStyle("Camera Access Required", fyne.TextAlignCenter, fyne.TextStyle{Bold: true})
	descLabel := widget.NewLabel(
		"Campus OS needs access to your local webcam to scan physical QR admission and executive credential dockets.\n\n" +
			"Live frames are analyzed strictly in-memory on your device and are never recorded or transmitted.",
	)
	descLabel.Wrapping = fyne.TextWrapWord

	rememberCheck := widget.NewCheck("Remember my choice on this device", nil)
	rememberCheck.SetChecked(true)

	denyBtn := NewShadcnButton("Don't Allow", ButtonOutline, ButtonSizeDefault, nil, func() {
		if rememberCheck.Checked {
			_ = SetCameraPermission(PermissionDenied)
		}
		if popup != nil {
			popup.Hide()
		}
		if onDeny != nil {
			onDeny()
		}
	})

	allowBtn := NewShadcnButton("Allow Camera Access", ButtonDefault, ButtonSizeDefault, ResourceFromSVG("check.svg", LucideCheckCircle2), func() {
		if rememberCheck.Checked {
			_ = SetCameraPermission(PermissionGranted)
		}
		if popup != nil {
			popup.Hide()
		}
		if onAllow != nil {
			onAllow()
		}
	})

	actionButtons := container.NewHBox(
		layout.NewSpacer(),
		denyBtn,
		allowBtn,
	)

	contentBox := container.NewVBox(
		container.NewCenter(iconVisual),
		container.NewCenter(titleLabel),
		widget.NewSeparator(),
		descLabel,
		rememberCheck,
	)

	card := NewShadcnCard(CardParts{
		Content: contentBox,
		Footer:  actionButtons,
	})

	cardMin := card.MinSize()
	targetWidth := float32(480)
	targetHeight := cardMin.Height
	if targetHeight < 380 {
		targetHeight = 380
	}

	popup = widget.NewModalPopUp(container.NewGridWrap(fyne.NewSize(targetWidth, targetHeight), card), window.Canvas())
	popup.Show()
	return popup
}

// DetectCameraDevice returns the first active v4l2 device path or empty if unavailable.
func DetectCameraDevice() string {
	candidates := []string{"/dev/video0", "/dev/video1", "/dev/video2"}
	for _, dev := range candidates {
		if _, err := os.Stat(dev); err == nil {
			return dev
		}
	}
	return ""
}

// BLOCK_UI_CAMERA_SCANNER_001
// Purpose: Displays a live optical camera viewfinder modal with real-time QR decoding.
func ShowCameraScannerModal(window fyne.Window, onTokenScanned func(token string), onCancel func()) *widget.PopUp {
	var popup *widget.PopUp

	device := DetectCameraDevice()
	mockRunner := getMockCameraRunner()
	if device == "" && mockRunner == nil {
		ShowToast(window, "No Camera Found", "No local webcam device was detected. Please connect a camera or upload a picture.", AlertDestructive, 4*time.Second)
		if onCancel != nil {
			onCancel()
		}
		return nil
	}

	ctx, cancel := context.WithCancel(context.Background())

	// Viewfinder frame image
	frameImg := canvas.NewImageFromImage(nil)
	frameImg.FillMode = canvas.ImageFillContain
	frameImg.SetMinSize(fyne.NewSize(480, 270))

	// Dark placeholder backdrop
	bgRect := canvas.NewRectangle(color.NRGBA{R: 15, G: 23, B: 42, A: 255})
	bgRect.SetMinSize(fyne.NewSize(480, 270))

	statusLabel := widget.NewLabelWithStyle("Initializing camera sensor...", fyne.TextAlignCenter, fyne.TextStyle{Italic: true})

	closeBtn := NewShadcnButton("Cancel / Switch to File", ButtonOutline, ButtonSizeDefault, nil, func() {
		cancel()
		if popup != nil {
			popup.Hide()
		}
		if onCancel != nil {
			onCancel()
		}
	})

	viewfinderStack := container.NewStack(
		bgRect,
		frameImg,
	)

	contentBox := container.NewVBox(
		container.NewHBox(
			RenderSVGImage(ResourceFromSVG("cam.svg", LucideCamera), 24, 24),
			widget.NewLabelWithStyle("Live Optical QR Scanner", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
			layout.NewSpacer(),
			NewBadge("LIVE FEED", BadgeSecondary, BadgeShapePill),
		),
		widget.NewSeparator(),
		viewfinderStack,
		statusLabel,
		container.NewBorder(nil, nil, nil, closeBtn),
	)

	card := NewShadcnCard(CardParts{
		Content: contentBox,
	})

	cardMin := card.MinSize()
	targetWidth := float32(520)
	targetHeight := cardMin.Height
	if targetHeight < 440 {
		targetHeight = 440
	}

	popup = widget.NewModalPopUp(container.NewGridWrap(fyne.NewSize(targetWidth, targetHeight), card), window.Canvas())
	popup.Show()

	// Launch live camera capture routine
	go func() {
		defer cancel()

		runner := getMockCameraRunner()
		if runner != nil {
			_ = runner(ctx, func(img image.Image) {
				fyne.Do(func() {
					frameImg.Image = img
					frameImg.Refresh()
				})
			}, func(tok string) {
				fyne.Do(func() {
					if popup != nil {
						popup.Hide()
					}
					if onTokenScanned != nil {
						onTokenScanned(tok)
					}
				})
			})
			return
		}

		cmd := exec.CommandContext(ctx, "ffmpeg",
			"-f", "v4l2",
			"-i", device,
			"-vf", "scale=640:360",
			"-r", "10",
			"-f", "image2pipe",
			"-vcodec", "mjpeg",
			"-",
		)

		stdout, err := cmd.StdoutPipe()
		if err != nil {
			fyne.Do(func() {
				statusLabel.SetText("Failed to access camera stream.")
			})
			return
		}

		if err := cmd.Start(); err != nil {
			fyne.Do(func() {
				statusLabel.SetText("Failed to start camera process.")
			})
			return
		}
		defer func() {
			if cmd.Process != nil {
				_ = cmd.Process.Kill()
			}
		}()

		fyne.Do(func() {
			statusLabel.SetText("Align docket QR code within the viewfinder")
		})

		buf := make([]byte, 8192)
		var frameBuf bytes.Buffer
		frameCount := 0

		for {
			select {
			case <-ctx.Done():
				return
			default:
			}

			n, err := stdout.Read(buf)
			if n > 0 {
				frameBuf.Write(buf[:n])
				data := frameBuf.Bytes()

				// Detect JPEG boundaries
				if idx := bytes.Index(data, []byte{0xFF, 0xD9}); idx != -1 {
					jpegBytes := data[:idx+2]
					remaining := data[idx+2:]

					img, decodeErr := jpeg.Decode(bytes.NewReader(jpegBytes))
					if decodeErr == nil {
						frameCount++
						fyne.Do(func() {
							frameImg.Image = img
							frameImg.Refresh()
						})

						// QR scanning
						token, err := DecodeQRFromLoadedImage(img)
						if err == nil && token != "" {
							cancel()
							fyne.Do(func() {
								if popup != nil {
									popup.Hide()
								}
								if onTokenScanned != nil {
									onTokenScanned(token)
								}
							})
							return
						}
					}

					frameBuf.Reset()
					frameBuf.Write(remaining)
				}
			}

			if err != nil {
				if err != io.EOF && ctx.Err() == nil {
					statusLabel.SetText("Camera stream stopped.")
				}
				break
			}
		}
	}()

	return popup
}
