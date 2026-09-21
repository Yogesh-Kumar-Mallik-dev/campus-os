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

	// Test Hover states
	sw.MouseIn(nil)
	if !sw.Hovered {
		t.Fatal("expected switch to be hovered after MouseIn")
	}
	sw.MouseOut()
	if sw.Hovered {
		t.Fatal("expected switch not to be hovered after MouseOut")
	}

	// Test Focus and Keyboard navigation accessibility
	sw.FocusGained()
	if !sw.Focused {
		t.Fatal("expected switch to be focused after FocusGained")
	}
	// Spacebar toggles switch
	sw.TypedRune(' ')
	if !sw.Checked {
		t.Fatal("expected spacebar key to toggle switch to checked")
	}
	// Enter key toggles switch
	sw.TypedKey(&fyne.KeyEvent{Name: fyne.KeyReturn})
	if sw.Checked {
		t.Fatal("expected enter key to toggle switch to unchecked")
	}
	sw.FocusLost()
	if sw.Focused {
		t.Fatal("expected switch not to be focused after FocusLost")
	}

	// Test Disabled state
	sw.Disable()
	if !sw.Disabled() {
		t.Fatal("expected switch to report disabled")
	}
	// Tapping while disabled must NOT toggle
	sw.Tapped(&fyne.PointEvent{})
	if sw.Checked {
		t.Fatal("expected disabled switch to ignore tap")
	}
	// Keyboard while disabled must NOT toggle
	sw.TypedRune(' ')
	if sw.Checked {
		t.Fatal("expected disabled switch to ignore keyboard")
	}
	// Re-enable
	sw.Enable()
	if sw.Disabled() {
		t.Fatal("expected switch to be re-enabled")
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
	if len(renderer.Objects()) != 3 {
		t.Errorf("expected 3 objects (focusRing, track, knob), got %d", len(renderer.Objects()))
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

	emptyObj.Resize(fyne.NewSize(400, 200))
	if emptyObj.Size().Width != 400 {
		t.Errorf("expected empty state to expand to width 400, got %f", emptyObj.Size().Width)
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

func TestSelectableCard_Interaction(t *testing.T) {
	tapped := false
	card := NewSelectableCard(widget.NewLabel("Card Content"), false, func() {
		tapped = true
	})

	if card == nil || card.Selected {
		t.Fatal("expected card to initialize unselected")
	}

	// Tap interaction
	card.Tapped(nil)
	if !tapped {
		t.Fatal("expected tapped handler to be called")
	}
	tapped = false

	// Hover interaction
	card.MouseIn(nil)
	if !card.Hovered {
		t.Fatal("expected card to be hovered")
	}
	card.MouseOut()
	if card.Hovered {
		t.Fatal("expected card not to be hovered")
	}

	// Active press interaction
	card.MouseDown(nil)
	if !card.Pressed {
		t.Fatal("expected card to be pressed")
	}
	card.MouseUp(nil)
	if card.Pressed {
		t.Fatal("expected card not to be pressed")
	}

	// Focus and Keyboard accessibility
	card.FocusGained()
	if !card.Focused {
		t.Fatal("expected card to be focused")
	}
	card.TypedRune(' ')
	if !tapped {
		t.Fatal("expected space key to activate card")
	}
	tapped = false

	card.TypedKey(&fyne.KeyEvent{Name: fyne.KeyReturn})
	if !tapped {
		t.Fatal("expected return key to activate card")
	}
	tapped = false
	card.FocusLost()
	if card.Focused {
		t.Fatal("expected card not to be focused after FocusLost")
	}

	// Selection state
	card.SetSelected(true)
	if !card.Selected {
		t.Fatal("expected card to be selected")
	}

	// Disabled state
	card.Disable()
	if !card.Disabled() {
		t.Fatal("expected card to report disabled")
	}
	card.Tapped(nil)
	if tapped {
		t.Fatal("expected disabled card to ignore tap")
	}
	card.TypedRune(' ')
	if tapped {
		t.Fatal("expected disabled card to ignore key")
	}

	card.Enable()
	if card.Disabled() {
		t.Fatal("expected card to be re-enabled")
	}

	// Renderer test
	renderer := test.WidgetRenderer(card)
	if renderer == nil {
		t.Fatal("expected non-nil renderer for selectable card")
	}
	renderer.Layout(fyne.NewSize(200, 100))
	minSize := renderer.MinSize()
	if minSize.Width <= 0 || minSize.Height <= 0 {
		t.Errorf("expected valid min size for selectable card, got %v", minSize)
	}
	renderer.Destroy()
}
