package ui

import (
	"testing"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/test"
	"fyne.io/fyne/v2/widget"
)

func TestSeparator_Render(t *testing.T) {
	sepHoriz := NewShadcnSeparator(true)
	if sepHoriz == nil {
		t.Fatal("expected non-nil horizontal separator")
	}

	sepVert := NewShadcnSeparator(false)
	if sepVert == nil {
		t.Fatal("expected non-nil vertical separator")
	}
}

func TestSwitch_Interaction(t *testing.T) {
	changed := false
	var lastState bool

	sw := NewSwitch(false, func(b bool) {
		changed = true
		lastState = b
	})

	if sw == nil || sw.Checked {
		t.Fatal("expected switch to initialize unchecked")
	}

	// Test programmatic SetChecked
	sw.SetChecked(true)
	if !changed || !lastState || !sw.Checked {
		t.Fatal("expected switch to update to checked")
	}

	// Test Tap interaction
	changed = false
	sw.Tapped(&fyne.PointEvent{})
	if !changed || lastState || sw.Checked {
		t.Fatal("expected switch tap to toggle to unchecked")
	}

	// Test renderer methods
	renderer := test.WidgetRenderer(sw)
	if renderer == nil {
		t.Fatal("expected non-nil widget renderer for switch")
	}
	renderer.Layout(fyne.NewSize(50, 30))
	minSize := renderer.MinSize()
	if minSize.Width <= 0 || minSize.Height <= 0 {
		t.Errorf("expected valid min size, got %v", minSize)
	}
	if len(renderer.Objects()) != 2 {
		t.Errorf("expected 2 objects (track, knob), got %d", len(renderer.Objects()))
	}
	renderer.Destroy()
}

func TestEmptyState_Render(t *testing.T) {
	media := widget.NewIcon(nil)
	actionBtn := widget.NewButton("Retry", nil)

	emptyObj := NewEmptyState(EmptyStateParams{
		Media:       media,
		Title:       "No Records Found",
		Description: "Your query returned 0 active sessions.",
		Action:      actionBtn,
	})

	if emptyObj == nil {
		t.Fatal("expected non-nil empty state component")
	}

	// Minimal empty state
	emptyMin := NewEmptyState(EmptyStateParams{})
	if emptyMin == nil {
		t.Fatal("expected non-nil minimal empty state")
	}
}

func TestToast_And_AlertDialog(t *testing.T) {
	testApp := test.NewApp()
	defer testApp.Quit()

	win := test.NewWindow(widget.NewLabel("Canvas Root"))
	win.Resize(fyne.NewSize(600, 400))

	// Toast notification test
	toast := ShowToast(win, "Action Saved", "Changes persisted successfully", AlertSuccess, 50*time.Millisecond)
	if toast == nil {
		t.Fatal("expected non-nil toast popup")
	}

	// Different toast variant
	toastWarn := ShowToast(win, "Warning", "Network latency high", AlertWarning, 0)
	if toastWarn == nil {
		t.Fatal("expected non-nil warning toast")
	}

	// Modal alert dialog test
	confirmed := false
	dlg := ShowAlertDialog(win, "Confirm Action", "Are you sure you want to proceed?", "Confirm", true, func() {
		confirmed = true
	})
	if dlg == nil {
		t.Fatal("expected non-nil alert dialog popup")
	}

	// Hide dialog
	dlg.Hide()
	if confirmed {
		t.Fatal("confirmed should not be called on dismiss without button click")
	}
}
