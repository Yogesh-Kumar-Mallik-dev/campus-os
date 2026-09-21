package ui

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

var _ fyne.Widget = (*Switch)(nil)
var _ fyne.Disableable = (*Switch)(nil)
var _ fyne.Focusable = (*Switch)(nil)
var _ desktop.Hoverable = (*Switch)(nil)
var _ desktop.Cursorable = (*Switch)(nil)
var _ desktop.Mouseable = (*Switch)(nil)

// Switch implements a modern toggle switch component mirroring shadcn switch.svelte
// with complete interactive states: Hover, Active/Pressed, Focus ring, and Disabled.
type Switch struct {
	widget.BaseWidget
	Checked   bool
	Hovered   bool
	Focused   bool
	Pressed   bool
	disabled  bool
	OnChanged func(bool)
}

// NewSwitch constructs a new toggle switch.
func NewSwitch(checked bool, onChanged func(bool)) *Switch {
	s := &Switch{
		Checked:   checked,
		OnChanged: onChanged,
	}
	s.ExtendBaseWidget(s)
	return s
}

// SetChecked programmatically updates the switch state.
func (s *Switch) SetChecked(checked bool) {
	if s.disabled {
		return
	}
	s.Checked = checked
	s.Refresh()
	if s.OnChanged != nil {
		s.OnChanged(s.Checked)
	}
}

// Enable marks the switch as interactive.
func (s *Switch) Enable() {
	s.disabled = false
	s.Refresh()
}

// Disable marks the switch as non-interactive with dimmed visual styling.
func (s *Switch) Disable() {
	s.disabled = true
	s.Focused = false
	s.Pressed = false
	s.Refresh()
}

// Disabled reports whether the switch is currently disabled.
func (s *Switch) Disabled() bool {
	return s.disabled
}

// Tapped toggles the checked state on click or touch tap.
func (s *Switch) Tapped(_ *fyne.PointEvent) {
	if s.disabled {
		return
	}
	s.SetChecked(!s.Checked)
}

// MouseIn handles hover entrance.
func (s *Switch) MouseIn(_ *desktop.MouseEvent) {
	if s.disabled {
		return
	}
	s.Hovered = true
	s.Refresh()
}

// MouseMoved handles hover tracking.
func (s *Switch) MouseMoved(_ *desktop.MouseEvent) {}

// MouseOut handles hover exit.
func (s *Switch) MouseOut() {
	s.Hovered = false
	s.Refresh()
}

// MouseDown handles active press state.
func (s *Switch) MouseDown(_ *desktop.MouseEvent) {
	if s.disabled {
		return
	}
	s.Pressed = true
	s.Refresh()
}

// MouseUp handles active press release.
func (s *Switch) MouseUp(_ *desktop.MouseEvent) {
	s.Pressed = false
	s.Refresh()
}

// Cursor returns pointer when active, or default when disabled.
func (s *Switch) Cursor() desktop.Cursor {
	if s.disabled {
		return desktop.DefaultCursor
	}
	return desktop.PointerCursor
}

// FocusGained activates the keyboard focus ring.
func (s *Switch) FocusGained() {
	if s.disabled {
		return
	}
	s.Focused = true
	s.Refresh()
}

// FocusLost deactivates the keyboard focus ring.
func (s *Switch) FocusLost() {
	s.Focused = false
	s.Refresh()
}

// TypedRune handles keyboard spacebar toggle for accessibility.
func (s *Switch) TypedRune(r rune) {
	if s.disabled {
		return
	}
	if r == ' ' {
		s.SetChecked(!s.Checked)
	}
}

// TypedKey handles Enter / Return key toggle for accessibility.
func (s *Switch) TypedKey(k *fyne.KeyEvent) {
	if s.disabled {
		return
	}
	if k.Name == fyne.KeySpace || k.Name == fyne.KeyReturn || k.Name == fyne.KeyEnter {
		s.SetChecked(!s.Checked)
	}
}

func (s *Switch) CreateRenderer() fyne.WidgetRenderer {
	focusRing := canvas.NewRectangle(color.Transparent)
	focusRing.CornerRadius = 14
	focusRing.StrokeWidth = 2
	focusRing.StrokeColor = theme.Color(theme.ColorNameFocus)
	focusRing.Hide()

	track := canvas.NewRectangle(color.NRGBA{R: 40, G: 46, B: 62, A: 255})
	track.CornerRadius = 11
	track.StrokeWidth = 1
	track.StrokeColor = theme.Color(theme.ColorNameInputBorder)

	knob := canvas.NewCircle(color.NRGBA{R: 255, G: 255, B: 255, A: 255})
	knob.Resize(fyne.NewSize(18, 18))

	r := &switchRenderer{
		sw:        s,
		focusRing: focusRing,
		track:     track,
		knob:      knob,
		objects:   []fyne.CanvasObject{focusRing, track, knob},
	}
	r.Refresh()
	return r
}

type switchRenderer struct {
	sw        *Switch
	focusRing *canvas.Rectangle
	track     *canvas.Rectangle
	knob      *canvas.Circle
	objects   []fyne.CanvasObject
}

func (r *switchRenderer) Layout(size fyne.Size) {
	// Focus ring wraps around the 44x24 track with 2px padding
	r.focusRing.Resize(fyne.NewSize(48, 28))
	r.focusRing.Move(fyne.NewPos(0, 0))

	r.track.Resize(fyne.NewSize(44, 24))
	r.track.Move(fyne.NewPos(2, 2))

	knobW := float32(18)
	if r.sw.Pressed && !r.sw.disabled {
		knobW = 21 // Subtle tactile stretch on active press
	}
	r.knob.Resize(fyne.NewSize(knobW, 18))

	if r.sw.Checked {
		if r.sw.Pressed && !r.sw.disabled {
			r.knob.Move(fyne.NewPos(22, 5))
		} else {
			r.knob.Move(fyne.NewPos(25, 5))
		}
	} else {
		r.knob.Move(fyne.NewPos(5, 5))
	}
}

func (r *switchRenderer) MinSize() fyne.Size {
	return fyne.NewSize(48, 28)
}

func (r *switchRenderer) Refresh() {
	if r.sw.Focused && !r.sw.disabled {
		r.focusRing.StrokeColor = theme.Color(theme.ColorNameFocus)
		r.focusRing.Show()
	} else {
		r.focusRing.Hide()
	}

	if r.sw.disabled {
		// Disabled state: dimmed contrast
		if r.sw.Checked {
			r.track.FillColor = color.NRGBA{R: 120, G: 50, B: 46, A: 130}
		} else {
			r.track.FillColor = color.NRGBA{R: 28, G: 34, B: 44, A: 130}
		}
		r.track.StrokeColor = color.NRGBA{R: 46, G: 56, B: 68, A: 100}
		r.knob.FillColor = color.NRGBA{R: 140, G: 150, B: 162, A: 150}
	} else if r.sw.Checked {
		// Active Checked state
		if r.sw.Pressed {
			r.track.FillColor = color.NRGBA{R: 215, G: 75, B: 68, A: 255}
		} else if r.sw.Hovered {
			r.track.FillColor = color.NRGBA{R: 255, G: 110, B: 102, A: 255}
		} else {
			r.track.FillColor = color.NRGBA{R: 244, G: 90, B: 81, A: 255}
		}
		r.track.StrokeColor = color.NRGBA{R: 255, G: 120, B: 112, A: 255}
		r.knob.FillColor = color.NRGBA{R: 255, G: 255, B: 255, A: 255}
	} else {
		// Inactive Unchecked state
		if r.sw.Pressed {
			r.track.FillColor = color.NRGBA{R: 54, G: 64, B: 80, A: 255}
		} else if r.sw.Hovered {
			r.track.FillColor = color.NRGBA{R: 48, G: 58, B: 74, A: 255}
		} else {
			r.track.FillColor = color.NRGBA{R: 35, G: 42, B: 54, A: 255}
		}
		if r.sw.Hovered {
			r.track.StrokeColor = color.NRGBA{R: 70, G: 84, B: 104, A: 255}
		} else {
			r.track.StrokeColor = color.NRGBA{R: 46, G: 56, B: 68, A: 255}
		}
		r.knob.FillColor = color.NRGBA{R: 255, G: 255, B: 255, A: 255}
	}

	r.Layout(r.sw.Size())
	r.focusRing.Refresh()
	r.track.Refresh()
	r.knob.Refresh()
}

func (r *switchRenderer) Objects() []fyne.CanvasObject {
	return r.objects
}

func (r *switchRenderer) Destroy() {}
