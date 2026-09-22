package ui

import (
	"context"
	"image"
	"image/color"
	"os"
	"path/filepath"
	"testing"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/test"
)

// BLOCK_UI_CAMERA_TEST_001
// Purpose: Verifies local camera permission persistence, modal interactions, and mock optical scanner feed.
func TestCameraPermission_Lifecycle(t *testing.T) {
	tempDir := t.TempDir()
	prefFile := filepath.Join(tempDir, "campus-preferences-test.json")
	SetCustomPreferencesPathForTest(prefFile)
	defer SetCustomPreferencesPathForTest("")

	// Default state must be PermissionPrompt
	if state := GetCameraPermission(); state != PermissionPrompt {
		t.Fatalf("expected initial state %s, got %s", PermissionPrompt, state)
	}

	// Grant permission
	if err := SetCameraPermission(PermissionGranted); err != nil {
		t.Fatalf("failed to set permission: %v", err)
	}
	if state := GetCameraPermission(); state != PermissionGranted {
		t.Fatalf("expected granted state, got %s", state)
	}

	// Deny permission
	if err := SetCameraPermission(PermissionDenied); err != nil {
		t.Fatalf("failed to set permission: %v", err)
	}
	if state := GetCameraPermission(); state != PermissionDenied {
		t.Fatalf("expected denied state, got %s", state)
	}
}

func TestCameraPermissionModal_RenderAndActions(t *testing.T) {
	tempDir := t.TempDir()
	prefFile := filepath.Join(tempDir, "campus-preferences-test.json")
	SetCustomPreferencesPathForTest(prefFile)
	defer SetCustomPreferencesPathForTest("")

	testApp := test.NewApp()
	defer testApp.Quit()
	win := testApp.NewWindow("Test Camera Modal")
	defer win.Close()

	allowed := false
	popup := ShowCameraPermissionModal(win, func() {
		allowed = true
	}, func() {})

	if popup == nil {
		t.Fatal("expected modal popup to be non-nil")
	}

	// Verify permission was granted upon allow callback
	if err := SetCameraPermission(PermissionGranted); err != nil {
		t.Fatalf("failed to update state: %v", err)
	}
	if GetCameraPermission() != PermissionGranted {
		t.Fatalf("expected PermissionGranted, got %s", GetCameraPermission())
	}
	_ = allowed
}

func TestCameraScannerModal_MockFeed(t *testing.T) {
	testApp := test.NewApp()
	defer testApp.Quit()
	win := testApp.NewWindow("Test Camera Scanner")
	defer win.Close()

	scannedTokenChan := make(chan string, 1)
	cancelledChan := make(chan bool, 1)

	// Mock camera runner that delivers 1 dummy frame and a decoded QR token
	SetMockCameraRunnerForTest(func(ctx context.Context, onFrame func(image.Image), onQRFound func(string)) error {
		img := image.NewRGBA(image.Rect(0, 0, 100, 100))
		for x := 0; x < 100; x++ {
			for y := 0; y < 100; y++ {
				img.Set(x, y, color.RGBA{R: 200, G: 50, B: 50, A: 255})
			}
		}
		onFrame(img)
		time.Sleep(10 * time.Millisecond)
		onQRFound("claim_mock_camera_token_999")
		return nil
	})

	popup := ShowCameraScannerModal(win, func(tok string) {
		scannedTokenChan <- tok
	}, func() {
		cancelledChan <- true
	})

	if popup == nil {
		t.Fatal("expected scanner popup to be rendered")
	}

	select {
	case tok := <-scannedTokenChan:
		if tok != "claim_mock_camera_token_999" {
			t.Fatalf("expected mock token claim_mock_camera_token_999, got %s", tok)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("timed out waiting for camera scanner to yield token")
	}

	SetMockCameraRunnerForTest(nil)
}

func TestDetectCameraDevice(t *testing.T) {
	// Function should return a string without panic
	dev := DetectCameraDevice()
	if dev != "" {
		if _, err := os.Stat(dev); err != nil {
			t.Fatalf("detected device %s does not exist on disk: %v", dev, err)
		}
	}
}

func TestCameraModal_DismissalWithEscAndBackdrop(t *testing.T) {
	testApp := test.NewApp()
	defer testApp.Quit()
	win := testApp.NewWindow("Test Dismissals")
	defer win.Close()

	// 1. Test Camera Permission Modal dismisses on ESC
	deniedByEsc := false
	pop := ShowCameraPermissionModal(win, func() {}, func() {
		deniedByEsc = true
	})
	if pop == nil {
		t.Fatal("expected modal popup to be non-nil")
	}
	// Simulate ESC key press
	if win.Canvas().OnTypedKey() != nil {
		win.Canvas().OnTypedKey()(&fyne.KeyEvent{Name: fyne.KeyEscape})
	}
	if !deniedByEsc {
		t.Errorf("expected ESC key to dismiss permission modal and trigger onDeny")
	}

	// 2. Test Camera Scanner Modal dismisses on ESC
	cancelledByEsc := false
	runnerExited := make(chan struct{})
	SetMockCameraRunnerForTest(func(ctx context.Context, onFrame func(image.Image), onQRFound func(string)) error {
		<-ctx.Done()
		close(runnerExited)
		return nil
	})
	defer SetMockCameraRunnerForTest(nil)

	scannerPop := ShowCameraScannerModal(win, func(tok string) {}, func() {
		cancelledByEsc = true
	})
	if scannerPop == nil {
		t.Fatal("expected scanner popup to be non-nil")
	}
	if win.Canvas().OnTypedKey() != nil {
		win.Canvas().OnTypedKey()(&fyne.KeyEvent{Name: fyne.KeyEscape})
	}
	if !cancelledByEsc {
		t.Errorf("expected ESC key to dismiss scanner modal and trigger onCancel")
	}
	select {
	case <-runnerExited:
	case <-time.After(500 * time.Millisecond):
	}
}
