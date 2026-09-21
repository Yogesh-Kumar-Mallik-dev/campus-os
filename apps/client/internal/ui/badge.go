package ui

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

// BadgeVariant enumerates the shadcn badge style variants.
type BadgeVariant int

const (
	BadgeDefault BadgeVariant = iota
	BadgeSecondary
	BadgeDestructive
	BadgeOutline
	BadgeSuccess
	BadgeWarning
)

// BadgeShape controls the corner geometry of the badge.
type BadgeShape int

const (
	BadgeShapePill BadgeShape = iota
	BadgeShapeRounded
)

// NewBadge constructs a shadcn-styled badge component matching institutional design tokens.
func NewBadge(text string, variant BadgeVariant, shape BadgeShape) fyne.CanvasObject {
	var bg color.Color
	var fg color.Color
	var border color.Color

	switch variant {
	case BadgeDefault:
		// Primary institutional terracotta: #F45A51
		bg = color.NRGBA{R: 244, G: 90, B: 81, A: 255}
		fg = color.NRGBA{R: 255, G: 255, B: 255, A: 255}
		border = color.NRGBA{R: 244, G: 90, B: 81, A: 255}
	case BadgeSecondary:
		// Secondary surface: #282E3E
		bg = color.NRGBA{R: 40, G: 46, B: 62, A: 255}
		fg = color.NRGBA{R: 242, G: 246, B: 248, A: 255}
		border = color.NRGBA{R: 46, G: 56, B: 68, A: 255}
	case BadgeDestructive:
		// Destructive red tint: #451818 bg, #F87171 fg
		bg = color.NRGBA{R: 69, G: 24, B: 24, A: 230}
		fg = color.NRGBA{R: 248, G: 113, B: 113, A: 255}
		border = color.NRGBA{R: 120, G: 35, B: 35, A: 255}
	case BadgeOutline:
		// Transparent with subtle border
		bg = color.NRGBA{R: 16, G: 22, B: 28, A: 0}
		fg = color.NRGBA{R: 242, G: 246, B: 248, A: 255}
		border = color.NRGBA{R: 60, G: 72, B: 92, A: 255}
	case BadgeSuccess:
		// Emerald success: #10402E bg, #34D399 fg
		bg = color.NRGBA{R: 16, G: 64, B: 46, A: 230}
		fg = color.NRGBA{R: 52, G: 211, B: 153, A: 255}
		border = color.NRGBA{R: 20, G: 100, B: 70, A: 255}
	case BadgeWarning:
		// Amber warning: #3D2A0E bg, #FBBF24 fg
		bg = color.NRGBA{R: 61, G: 42, B: 14, A: 230}
		fg = color.NRGBA{R: 251, G: 191, B: 36, A: 255}
		border = color.NRGBA{R: 120, G: 85, B: 20, A: 255}
	}

	bgRect := canvas.NewRectangle(bg)
	bgRect.StrokeColor = border
	bgRect.StrokeWidth = 1

	if shape == BadgeShapePill {
		bgRect.CornerRadius = 12
	} else {
		bgRect.CornerRadius = 5
	}

	badgeText := canvas.NewText("  "+text+"  ", fg)
	badgeText.TextSize = 11
	badgeText.TextStyle = fyne.TextStyle{Bold: true}
	badgeText.Alignment = fyne.TextAlignCenter

	return container.NewStack(
		bgRect,
		container.NewPadded(badgeText),
	)
}

// NewStatusBadge renders an accessible status badge featuring a circular status indicator dot alongside text,
// adhering to Guardrail 15 (Avoid showing raw status as text; combine visual indicators with text for accessibility).
func NewStatusBadge(text string, variant BadgeVariant) fyne.CanvasObject {
	return NewBadge("● "+text, variant, BadgeShapePill)
}

var _ fyne.Widget = (*InteractiveBadge)(nil)
var _ fyne.Disableable = (*InteractiveBadge)(nil)
var _ fyne.Focusable = (*InteractiveBadge)(nil)
var _ desktop.Hoverable = (*InteractiveBadge)(nil)
var _ desktop.Cursorable = (*InteractiveBadge)(nil)
var _ desktop.Mouseable = (*InteractiveBadge)(nil)

// InteractiveBadge represents a clickable tag/chip badge with complete
// Hover, Active/Pressed, Focus ring, and Disabled states.
type InteractiveBadge struct {
	widget.BaseWidget
	Text     string
	Variant  BadgeVariant
	Shape    BadgeShape
	Hovered  bool
	Focused  bool
	Pressed  bool
	disabled bool
	OnTapped func()
}

// NewInteractiveBadge constructs an interactive badge chip.
func NewInteractiveBadge(text string, variant BadgeVariant, shape BadgeShape, onTapped func()) *InteractiveBadge {
	b := &InteractiveBadge{
		Text:     text,
		Variant:  variant,
		Shape:    shape,
		OnTapped: onTapped,
	}
	b.ExtendBaseWidget(b)
	return b
}

// Enable enables the badge interaction.
func (b *InteractiveBadge) Enable() {
	b.disabled = false
	b.Refresh()
}

// Disable disables the badge interaction.
func (b *InteractiveBadge) Disable() {
	b.disabled = true
	b.Focused = false
	b.Pressed = false
	b.Refresh()
}

// Disabled reports whether the badge is disabled.
func (b *InteractiveBadge) Disabled() bool {
	return b.disabled
}

// Tapped handles click events.
func (b *InteractiveBadge) Tapped(_ *fyne.PointEvent) {
	if b.disabled {
		return
	}
	if b.OnTapped != nil {
		b.OnTapped()
	}
}

// MouseIn handles hover entrance.
func (b *InteractiveBadge) MouseIn(_ *desktop.MouseEvent) {
	if b.disabled {
		return
	}
	b.Hovered = true
	b.Refresh()
}

// MouseMoved handles hover move.
func (b *InteractiveBadge) MouseMoved(_ *desktop.MouseEvent) {}

// MouseOut handles hover exit.
func (b *InteractiveBadge) MouseOut() {
	b.Hovered = false
	b.Refresh()
}

// MouseDown handles press state.
func (b *InteractiveBadge) MouseDown(_ *desktop.MouseEvent) {
	if b.disabled {
		return
	}
	b.Pressed = true
	b.Refresh()
}

// MouseUp handles press release.
func (b *InteractiveBadge) MouseUp(_ *desktop.MouseEvent) {
	b.Pressed = false
	b.Refresh()
}

// Cursor returns pointer when active, default when disabled.
func (b *InteractiveBadge) Cursor() desktop.Cursor {
	if b.disabled {
		return desktop.DefaultCursor
	}
	return desktop.PointerCursor
}

// FocusGained activates the keyboard focus ring.
func (b *InteractiveBadge) FocusGained() {
	if b.disabled {
		return
	}
	b.Focused = true
	b.Refresh()
}

// FocusLost deactivates the keyboard focus ring.
func (b *InteractiveBadge) FocusLost() {
	b.Focused = false
	b.Refresh()
}

// TypedRune handles spacebar activation.
func (b *InteractiveBadge) TypedRune(r rune) {
	if b.disabled {
		return
	}
	if r == ' ' && b.OnTapped != nil {
		b.OnTapped()
	}
}

// TypedKey handles Enter / Return key activation.
func (b *InteractiveBadge) TypedKey(k *fyne.KeyEvent) {
	if b.disabled {
		return
	}
	if (k.Name == fyne.KeySpace || k.Name == fyne.KeyReturn || k.Name == fyne.KeyEnter) && b.OnTapped != nil {
		b.OnTapped()
	}
}

func (b *InteractiveBadge) CreateRenderer() fyne.WidgetRenderer {
	focusRing := canvas.NewRectangle(color.Transparent)
	focusRing.StrokeWidth = 2
	focusRing.StrokeColor = theme.Color(theme.ColorNameFocus)
	focusRing.Hide()

	bgRect := canvas.NewRectangle(color.Transparent)
	bgRect.StrokeWidth = 1

	if b.Shape == BadgeShapePill {
		focusRing.CornerRadius = 14
		bgRect.CornerRadius = 12
	} else {
		focusRing.CornerRadius = 7
		bgRect.CornerRadius = 5
	}

	txt := canvas.NewText("  "+b.Text+"  ", color.White)
	txt.TextSize = 11
	txt.TextStyle = fyne.TextStyle{Bold: true}
	txt.Alignment = fyne.TextAlignCenter

	r := &interactiveBadgeRenderer{
		badge:     b,
		focusRing: focusRing,
		bgRect:    bgRect,
		txt:       txt,
		objects:   []fyne.CanvasObject{focusRing, bgRect, txt},
	}
	r.Refresh()
	return r
}

type interactiveBadgeRenderer struct {
	badge     *InteractiveBadge
	focusRing *canvas.Rectangle
	bgRect    *canvas.Rectangle
	txt       *canvas.Text
	objects   []fyne.CanvasObject
}

func (r *interactiveBadgeRenderer) Layout(size fyne.Size) {
	r.focusRing.Resize(size)
	r.focusRing.Move(fyne.NewPos(0, 0))

	pad := float32(2)
	r.bgRect.Resize(fyne.NewSize(size.Width-pad*2, size.Height-pad*2))
	r.bgRect.Move(fyne.NewPos(pad, pad))

	r.txt.Resize(fyne.NewSize(size.Width-pad*2, size.Height-pad*2))
	r.txt.Move(fyne.NewPos(pad, pad+2))
}

func (r *interactiveBadgeRenderer) MinSize() fyne.Size {
	ts := r.txt.MinSize()
	return fyne.NewSize(ts.Width+16, ts.Height+10)
}

func (r *interactiveBadgeRenderer) Refresh() {
	if r.badge.Focused && !r.badge.disabled {
		r.focusRing.StrokeColor = theme.Color(theme.ColorNameFocus)
		r.focusRing.Show()
	} else {
		r.focusRing.Hide()
	}

	var bg, fg, border color.Color

	switch r.badge.Variant {
	case BadgeDefault:
		bg = color.NRGBA{R: 244, G: 90, B: 81, A: 255}
		fg = color.NRGBA{R: 255, G: 255, B: 255, A: 255}
		border = color.NRGBA{R: 244, G: 90, B: 81, A: 255}
	case BadgeSecondary:
		bg = color.NRGBA{R: 40, G: 46, B: 62, A: 255}
		fg = color.NRGBA{R: 242, G: 246, B: 248, A: 255}
		border = color.NRGBA{R: 46, G: 56, B: 68, A: 255}
	case BadgeDestructive:
		bg = color.NRGBA{R: 69, G: 24, B: 24, A: 230}
		fg = color.NRGBA{R: 248, G: 113, B: 113, A: 255}
		border = color.NRGBA{R: 120, G: 35, B: 35, A: 255}
	case BadgeOutline:
		bg = color.NRGBA{R: 16, G: 22, B: 28, A: 0}
		fg = color.NRGBA{R: 242, G: 246, B: 248, A: 255}
		border = color.NRGBA{R: 60, G: 72, B: 92, A: 255}
	case BadgeSuccess:
		bg = color.NRGBA{R: 16, G: 64, B: 46, A: 230}
		fg = color.NRGBA{R: 52, G: 211, B: 153, A: 255}
		border = color.NRGBA{R: 20, G: 100, B: 70, A: 255}
	case BadgeWarning:
		bg = color.NRGBA{R: 61, G: 42, B: 14, A: 230}
		fg = color.NRGBA{R: 251, G: 191, B: 36, A: 255}
		border = color.NRGBA{R: 120, G: 85, B: 20, A: 255}
	}

	if r.badge.disabled {
		bg = color.NRGBA{R: 25, G: 32, B: 42, A: 120}
		fg = color.NRGBA{R: 110, G: 122, B: 138, A: 160}
		border = color.NRGBA{R: 40, G: 48, B: 58, A: 100}
	} else if r.badge.Pressed {
		border = theme.Color(theme.ColorNamePrimary)
	} else if r.badge.Hovered {
		border = color.NRGBA{R: 255, G: 255, B: 255, A: 200}
	}

	r.bgRect.FillColor = bg
	r.bgRect.StrokeColor = border
	r.txt.Color = fg
	r.txt.Text = "  " + r.badge.Text + "  "

	r.Layout(r.badge.Size())
	r.focusRing.Refresh()
	r.bgRect.Refresh()
	r.txt.Refresh()
}

func (r *interactiveBadgeRenderer) Objects() []fyne.CanvasObject {
	return r.objects
}

func (r *interactiveBadgeRenderer) Destroy() {}

