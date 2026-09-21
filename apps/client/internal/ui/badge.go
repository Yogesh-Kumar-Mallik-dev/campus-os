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

// badgeColors returns the modern 2026 luminous color palette for badges.
func badgeColors(variant BadgeVariant) (bg, fg, border color.Color) {
	switch variant {
	case BadgeDefault:
		// Refined Primary Terracotta: solid crisp terracotta with luminous border
		bg = color.NRGBA{R: 244, G: 90, B: 81, A: 255}
		fg = color.NRGBA{R: 255, G: 255, B: 255, A: 255}
		border = color.NRGBA{R: 255, G: 114, B: 104, A: 255}
	case BadgeSecondary:
		// Glass Surface: translucent luminous white overlay on dark card
		bg = color.NRGBA{R: 255, G: 255, B: 255, A: 18}
		fg = color.NRGBA{R: 226, G: 232, B: 240, A: 240}
		border = color.NRGBA{R: 255, G: 255, B: 255, A: 32}
	case BadgeDestructive:
		// Luminous Rose tint
		bg = color.NRGBA{R: 239, G: 68, B: 68, A: 35}
		fg = color.NRGBA{R: 248, G: 113, B: 113, A: 255}
		border = color.NRGBA{R: 239, G: 68, B: 68, A: 85}
	case BadgeOutline:
		// Elegant Minimalist Hairline Pill
		bg = color.NRGBA{R: 255, G: 255, B: 255, A: 8}
		fg = color.NRGBA{R: 160, G: 174, B: 192, A: 230}
		border = color.NRGBA{R: 255, G: 255, B: 255, A: 38}
	case BadgeSuccess:
		// Luminous Emerald tint
		bg = color.NRGBA{R: 16, G: 185, B: 129, A: 35}
		fg = color.NRGBA{R: 52, G: 211, B: 153, A: 255}
		border = color.NRGBA{R: 16, G: 185, B: 129, A: 85}
	case BadgeWarning:
		// Luminous Amber tint
		bg = color.NRGBA{R: 245, G: 158, B: 11, A: 35}
		fg = color.NRGBA{R: 251, G: 191, B: 36, A: 255}
		border = color.NRGBA{R: 245, G: 158, B: 11, A: 85}
	default:
		bg = color.NRGBA{R: 255, G: 255, B: 255, A: 14}
		fg = color.NRGBA{R: 226, G: 232, B: 240, A: 240}
		border = color.NRGBA{R: 255, G: 255, B: 255, A: 28}
	}
	return
}

type badgeLayout struct {
	padX  float32
	padY  float32
	shape BadgeShape
}

func (l *badgeLayout) Layout(objects []fyne.CanvasObject, size fyne.Size) {
	if len(objects) == 0 {
		return
	}

	bg := objects[0]
	bg.Resize(size)
	bg.Move(fyne.NewPos(0, 0))

	if rect, ok := bg.(*canvas.Rectangle); ok {
		if l.shape == BadgeShapePill {
			rect.CornerRadius = size.Height / 2
		} else {
			rect.CornerRadius = 4
		}
	}

	if len(objects) > 1 {
		inner := objects[1]
		innerW := size.Width - (l.padX * 2)
		innerH := size.Height - (l.padY * 2)
		if innerW < 0 {
			innerW = 0
		}
		if innerH < 0 {
			innerH = 0
		}
		inner.Resize(fyne.NewSize(innerW, innerH))
		inner.Move(fyne.NewPos(l.padX, l.padY))
	}
}

func (l *badgeLayout) MinSize(objects []fyne.CanvasObject) fyne.Size {
	if len(objects) < 2 {
		return fyne.NewSize(24, 20)
	}
	innerMin := objects[1].MinSize()
	w := innerMin.Width + (l.padX * 2)
	h := innerMin.Height + (l.padY * 2)
	if h < 20 {
		h = 20
	}
	return fyne.NewSize(w, h)
}

// NewBadge constructs a modern 2026 badge component with pixel-perfect pill curvature and luminous palette.
func NewBadge(text string, variant BadgeVariant, shape BadgeShape) fyne.CanvasObject {
	bg, fg, border := badgeColors(variant)

	bgRect := canvas.NewRectangle(bg)
	bgRect.StrokeColor = border
	bgRect.StrokeWidth = 1

	badgeText := canvas.NewText(text, fg)
	badgeText.TextSize = 10.5
	badgeText.TextStyle = fyne.TextStyle{Bold: true}
	badgeText.Alignment = fyne.TextAlignCenter

	return container.New(&badgeLayout{padX: 9, padY: 2.5, shape: shape}, bgRect, badgeText)
}

// NewStatusBadge renders an accessible status badge featuring a circular status indicator dot alongside text.
func NewStatusBadge(text string, variant BadgeVariant) fyne.CanvasObject {
	var dotColor color.Color
	switch variant {
	case BadgeSuccess:
		dotColor = color.NRGBA{R: 52, G: 211, B: 153, A: 255}
	case BadgeWarning:
		dotColor = color.NRGBA{R: 251, G: 191, B: 36, A: 255}
	case BadgeDestructive:
		dotColor = color.NRGBA{R: 248, G: 113, B: 113, A: 255}
	case BadgeDefault:
		dotColor = color.NRGBA{R: 255, G: 255, B: 255, A: 255}
	default:
		dotColor = color.NRGBA{R: 160, G: 174, B: 192, A: 255}
	}

	dot := canvas.NewCircle(dotColor)
	dot.Resize(fyne.NewSize(5, 5))

	bg, fg, border := badgeColors(variant)

	label := canvas.NewText(text, fg)
	label.TextSize = 10.5
	label.TextStyle = fyne.TextStyle{Bold: true}

	content := container.NewHBox(container.NewCenter(dot), label)

	bgRect := canvas.NewRectangle(bg)
	bgRect.StrokeColor = border
	bgRect.StrokeWidth = 1

	return container.New(&badgeLayout{padX: 8, padY: 2.5, shape: BadgeShapePill}, bgRect, content)
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

	txt := canvas.NewText(b.Text, color.White)
	txt.TextSize = 10.5
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

	pad := float32(1.5)
	innerW := size.Width - pad*2
	innerH := size.Height - pad*2
	if innerW < 0 {
		innerW = 0
	}
	if innerH < 0 {
		innerH = 0
	}
	r.bgRect.Resize(fyne.NewSize(innerW, innerH))
	r.bgRect.Move(fyne.NewPos(pad, pad))

	if r.badge.Shape == BadgeShapePill {
		r.bgRect.CornerRadius = innerH / 2
		r.focusRing.CornerRadius = size.Height / 2
	} else {
		r.bgRect.CornerRadius = 4
		r.focusRing.CornerRadius = 6
	}

	txtW := innerW - 14
	if txtW < 0 {
		txtW = 0
	}
	r.txt.Resize(fyne.NewSize(txtW, innerH))
	r.txt.Move(fyne.NewPos(pad+7, pad))
}

func (r *interactiveBadgeRenderer) MinSize() fyne.Size {
	ts := r.txt.MinSize()
	h := ts.Height + 5
	if h < 20 {
		h = 20
	}
	return fyne.NewSize(ts.Width+18, h)
}

func (r *interactiveBadgeRenderer) Refresh() {
	if r.badge.Focused && !r.badge.disabled {
		r.focusRing.StrokeColor = theme.Color(theme.ColorNameFocus)
		r.focusRing.Show()
	} else {
		r.focusRing.Hide()
	}

	bg, fg, border := badgeColors(r.badge.Variant)

	if r.badge.disabled {
		bg = color.NRGBA{R: 255, G: 255, B: 255, A: 6}
		fg = color.NRGBA{R: 110, G: 122, B: 138, A: 160}
		border = color.NRGBA{R: 255, G: 255, B: 255, A: 16}
	} else if r.badge.Pressed {
		border = theme.Color(theme.ColorNamePrimary)
	} else if r.badge.Hovered {
		border = color.NRGBA{R: 255, G: 255, B: 255, A: 120}
	}

	r.bgRect.FillColor = bg
	r.bgRect.StrokeColor = border
	r.txt.Color = fg
	r.txt.Text = r.badge.Text

	r.Layout(r.badge.Size())
	r.focusRing.Refresh()
	r.bgRect.Refresh()
	r.txt.Refresh()
}

func (r *interactiveBadgeRenderer) Objects() []fyne.CanvasObject {
	return r.objects
}

func (r *interactiveBadgeRenderer) Destroy() {}

