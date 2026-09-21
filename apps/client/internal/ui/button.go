package ui

import (
	"image/color"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

// ButtonVariant enumerates shadcn button style variants.
type ButtonVariant int

const (
	ButtonDefault ButtonVariant = iota
	ButtonSecondary
	ButtonDestructive
	ButtonOutline
	ButtonGhost
)

// ButtonSize controls button vertical padding and font size.
type ButtonSize int

const (
	ButtonSizeDefault ButtonSize = iota
	ButtonSizeSm
	ButtonSizeLg
	ButtonSizeIcon
)

// ShadcnButton extends Fyne's Button with shadcn variant tokens, accessible hit targets,
// and full interactive states (Hover, Active, Focus, and Disabled).
type ShadcnButton struct {
	widget.Button
	Variant    ButtonVariant
	ButtonSize ButtonSize
}

// NewShadcnButton constructs a button adhering to shadcn variants and sizing.
func NewShadcnButton(label string, variant ButtonVariant, size ButtonSize, icon fyne.Resource, onClick func()) *ShadcnButton {
	btn := &ShadcnButton{
		Variant:    variant,
		ButtonSize: size,
	}
	btn.Text = label
	btn.Icon = icon
	btn.OnTapped = onClick

	switch variant {
	case ButtonDefault:
		btn.Importance = widget.HighImportance
	case ButtonSecondary:
		btn.Importance = widget.MediumImportance
	case ButtonDestructive:
		btn.Importance = widget.DangerImportance
	case ButtonOutline, ButtonGhost:
		btn.Importance = widget.LowImportance
	}

	btn.ExtendBaseWidget(btn)
	return btn
}

// MinSize enforces Guardrail 26: comfortable minimum touch and click targets (>= 36-38px height).
func (b *ShadcnButton) MinSize() fyne.Size {
	base := b.Button.MinSize()
	targetHeight := float32(38)
	targetWidth := base.Width

	switch b.ButtonSize {
	case ButtonSizeSm:
		targetHeight = 32
	case ButtonSizeLg:
		targetHeight = 44
	case ButtonSizeIcon:
		targetHeight = 38
		if targetWidth < 38 {
			targetWidth = 38
		}
	default:
		targetHeight = 38
	}

	if base.Height < targetHeight {
		base.Height = targetHeight
	}
	if base.Width < targetWidth {
		base.Width = targetWidth
	}
	return base
}

var _ fyne.Widget = (*AnimatedArrowButton)(nil)
var _ fyne.Disableable = (*AnimatedArrowButton)(nil)
var _ fyne.Focusable = (*AnimatedArrowButton)(nil)
var _ desktop.Hoverable = (*AnimatedArrowButton)(nil)
var _ desktop.Cursorable = (*AnimatedArrowButton)(nil)
var _ desktop.Mouseable = (*AnimatedArrowButton)(nil)
var _ fyne.Tappable = (*AnimatedArrowButton)(nil)

// AnimatedArrowButton represents a primary action button with an integrated trailing arrow
// that smoothly translates horizontally on hover to provide forward progression affordance.
type AnimatedArrowButton struct {
	widget.BaseWidget
	Text        string
	Variant     ButtonVariant
	ButtonSize  ButtonSize
	LeadIcon    fyne.Resource
	arrowOffset float32
	Hovered     bool
	Focused     bool
	Pressed     bool
	disabled    bool
	OnTapped    func()
	anim        *fyne.Animation
}

// NewAnimatedArrowButton constructs a button with an integrated trailing arrow that animates on hover.
func NewAnimatedArrowButton(label string, variant ButtonVariant, size ButtonSize, leadIcon fyne.Resource, onClick func()) *AnimatedArrowButton {
	b := &AnimatedArrowButton{
		Text:       label,
		Variant:    variant,
		ButtonSize: size,
		LeadIcon:   leadIcon,
		OnTapped:   onClick,
	}
	b.ExtendBaseWidget(b)
	return b
}

func (b *AnimatedArrowButton) startHoverAnimation(target float32) {
	if b.anim != nil {
		b.anim.Stop()
	}
	start := b.arrowOffset
	diff := target - start
	b.anim = fyne.NewAnimation(150*time.Millisecond, func(progress float32) {
		b.arrowOffset = start + diff*progress
		b.Refresh()
	})
	b.anim.Start()
}

// SetText updates the button label text.
func (b *AnimatedArrowButton) SetText(text string) {
	b.Text = text
	b.Refresh()
}

// Enable enables interaction on the button.
func (b *AnimatedArrowButton) Enable() {
	b.disabled = false
	b.Refresh()
}

// Disable disables interaction on the button.
func (b *AnimatedArrowButton) Disable() {
	b.disabled = true
	b.Focused = false
	b.Pressed = false
	b.startHoverAnimation(0)
	b.Refresh()
}

// Disabled reports whether the button is disabled.
func (b *AnimatedArrowButton) Disabled() bool {
	return b.disabled
}

// Tapped executes the click action when enabled.
func (b *AnimatedArrowButton) Tapped(_ *fyne.PointEvent) {
	if b.disabled {
		return
	}
	if b.OnTapped != nil {
		b.OnTapped()
	}
}

// MouseIn handles hover entrance and animates the arrow sliding forward.
func (b *AnimatedArrowButton) MouseIn(_ *desktop.MouseEvent) {
	if b.disabled {
		return
	}
	b.Hovered = true
	b.startHoverAnimation(5.0)
}

// MouseMoved handles hover move.
func (b *AnimatedArrowButton) MouseMoved(_ *desktop.MouseEvent) {}

// MouseOut handles hover exit and slides the arrow back to resting position.
func (b *AnimatedArrowButton) MouseOut() {
	b.Hovered = false
	b.startHoverAnimation(0.0)
}

// MouseDown handles press feedback.
func (b *AnimatedArrowButton) MouseDown(_ *desktop.MouseEvent) {
	if b.disabled {
		return
	}
	b.Pressed = true
	b.Refresh()
}

// MouseUp handles press release.
func (b *AnimatedArrowButton) MouseUp(_ *desktop.MouseEvent) {
	b.Pressed = false
	b.Refresh()
}

// Cursor returns pointer when enabled, default when disabled.
func (b *AnimatedArrowButton) Cursor() desktop.Cursor {
	if b.disabled {
		return desktop.DefaultCursor
	}
	return desktop.PointerCursor
}

// FocusGained activates the keyboard focus ring.
func (b *AnimatedArrowButton) FocusGained() {
	if b.disabled {
		return
	}
	b.Focused = true
	b.Refresh()
}

// FocusLost deactivates the keyboard focus ring.
func (b *AnimatedArrowButton) FocusLost() {
	b.Focused = false
	b.Refresh()
}

// TypedRune handles spacebar activation.
func (b *AnimatedArrowButton) TypedRune(r rune) {
	if b.disabled {
		return
	}
	if r == ' ' && b.OnTapped != nil {
		b.OnTapped()
	}
}

// TypedKey handles Enter / Return key activation.
func (b *AnimatedArrowButton) TypedKey(k *fyne.KeyEvent) {
	if b.disabled {
		return
	}
	if (k.Name == fyne.KeySpace || k.Name == fyne.KeyReturn || k.Name == fyne.KeyEnter) && b.OnTapped != nil {
		b.OnTapped()
	}
}

func (b *AnimatedArrowButton) CreateRenderer() fyne.WidgetRenderer {
	focusRing := canvas.NewRectangle(color.Transparent)
	focusRing.CornerRadius = 10
	focusRing.StrokeWidth = 2
	focusRing.StrokeColor = theme.Color(theme.ColorNameFocus)
	focusRing.Hide()

	bg := canvas.NewRectangle(color.NRGBA{R: 244, G: 90, B: 81, A: 255})
	bg.CornerRadius = 8

	label := canvas.NewText(b.Text, color.White)
	label.TextSize = 13.5
	label.TextStyle = fyne.TextStyle{Bold: true}
	label.Alignment = fyne.TextAlignCenter

	arrowImg := canvas.NewImageFromResource(WhiteResourceFromSVG("arrow_right.svg", LucideArrowRight))
	arrowImg.FillMode = canvas.ImageFillContain

	objects := []fyne.CanvasObject{focusRing, bg, label, arrowImg}

	var leadImg *canvas.Image
	if b.LeadIcon != nil {
		leadImg = canvas.NewImageFromResource(b.LeadIcon)
		leadImg.FillMode = canvas.ImageFillContain
		objects = append(objects, leadImg)
	}

	r := &animatedArrowButtonRenderer{
		btn:       b,
		focusRing: focusRing,
		bg:        bg,
		label:     label,
		arrowImg:  arrowImg,
		leadImg:   leadImg,
		objects:   objects,
	}
	r.Refresh()
	return r
}

type animatedArrowButtonRenderer struct {
	btn       *AnimatedArrowButton
	focusRing *canvas.Rectangle
	bg        *canvas.Rectangle
	label     *canvas.Text
	arrowImg  *canvas.Image
	leadImg   *canvas.Image
	objects   []fyne.CanvasObject
}

func (r *animatedArrowButtonRenderer) Layout(size fyne.Size) {
	r.focusRing.Resize(size)
	r.focusRing.Move(fyne.NewPos(0, 0))

	pad := float32(2)
	r.bg.Resize(fyne.NewSize(size.Width-pad*2, size.Height-pad*2))
	r.bg.Move(fyne.NewPos(pad, pad))

	labelMin := r.label.MinSize()
	iconSize := float32(18)
	gap := float32(8)

	totalWidth := labelMin.Width + gap + iconSize
	if r.leadImg != nil {
		totalWidth += iconSize + gap
	}

	startX := (size.Width - totalWidth) / 2
	centerY := size.Height / 2

	curX := startX
	if r.leadImg != nil {
		r.leadImg.Resize(fyne.NewSize(iconSize, iconSize))
		r.leadImg.Move(fyne.NewPos(curX, centerY-iconSize/2))
		curX += iconSize + gap
	}

	r.label.Resize(fyne.NewSize(labelMin.Width, labelMin.Height))
	r.label.Move(fyne.NewPos(curX, centerY-labelMin.Height/2))
	curX += labelMin.Width + gap

	// The arrow's X position dynamically reflects r.btn.arrowOffset (animated translation)
	r.arrowImg.Resize(fyne.NewSize(iconSize, iconSize))
	r.arrowImg.Move(fyne.NewPos(curX+r.btn.arrowOffset, centerY-iconSize/2))
}

func (r *animatedArrowButtonRenderer) MinSize() fyne.Size {
	labelMin := r.label.MinSize()
	iconSize := float32(18)
	gap := float32(8)

	totalWidth := labelMin.Width + gap + iconSize + 36
	if r.leadImg != nil {
		totalWidth += iconSize + gap
	}

	h := float32(38)
	if r.btn.ButtonSize == ButtonSizeLg {
		h = 44
	} else if r.btn.ButtonSize == ButtonSizeSm {
		h = 32
	}
	return fyne.NewSize(totalWidth, h)
}

func (r *animatedArrowButtonRenderer) Refresh() {
	if r.btn.Focused && !r.btn.disabled {
		r.focusRing.StrokeColor = theme.Color(theme.ColorNameFocus)
		r.focusRing.Show()
	} else {
		r.focusRing.Hide()
	}

	if r.btn.disabled {
		r.bg.FillColor = color.NRGBA{R: 24, G: 30, B: 40, A: 255}
		r.label.Color = color.NRGBA{R: 110, G: 122, B: 138, A: 255}
		r.arrowImg.Resource = ResourceFromSVG("arrow_dis.svg", LucideArrowRight)
	} else {
		if r.btn.Pressed {
			r.bg.FillColor = color.NRGBA{R: 215, G: 75, B: 68, A: 255}
		} else if r.btn.Hovered {
			r.bg.FillColor = color.NRGBA{R: 255, G: 108, B: 100, A: 255}
		} else {
			r.bg.FillColor = color.NRGBA{R: 244, G: 90, B: 81, A: 255}
		}
		r.label.Color = color.White
		r.arrowImg.Resource = WhiteResourceFromSVG("arrow_right.svg", LucideArrowRight)
	}

	r.label.Text = r.btn.Text
	r.Layout(r.btn.Size())
	r.focusRing.Refresh()
	r.bg.Refresh()
	r.label.Refresh()
	r.arrowImg.Refresh()
	if r.leadImg != nil {
		r.leadImg.Refresh()
	}
}

func (r *animatedArrowButtonRenderer) Objects() []fyne.CanvasObject {
	return r.objects
}

func (r *animatedArrowButtonRenderer) Destroy() {
	if r.btn.anim != nil {
		r.btn.anim.Stop()
	}
}

