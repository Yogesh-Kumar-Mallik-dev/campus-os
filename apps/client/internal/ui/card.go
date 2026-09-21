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

// CardParts defines structural sections of a compound shadcn Card.
type CardParts struct {
	Title       string
	Description string
	Badge       fyne.CanvasObject
	Content     fyne.CanvasObject
	Footer      fyne.CanvasObject
}

// NewShadcnCard constructs a compound card widget mirroring shadcn-svelte's Card compound architecture.
func NewShadcnCard(parts CardParts) fyne.CanvasObject {
	cardBg := canvas.NewRectangle(theme.Color(theme.ColorNameMenuBackground))
	cardBg.CornerRadius = 10
	cardBg.StrokeColor = theme.Color(theme.ColorNameInputBorder)
	cardBg.StrokeWidth = 1

	elements := make([]fyne.CanvasObject, 0, 4)

	// Header section
	if parts.Title != "" || parts.Description != "" || parts.Badge != nil {
		headerItems := make([]fyne.CanvasObject, 0, 3)

		if parts.Badge != nil {
			headerItems = append(headerItems, container.NewHBox(parts.Badge))
		}

		if parts.Title != "" {
			titleLbl := widget.NewLabelWithStyle(parts.Title, fyne.TextAlignLeading, fyne.TextStyle{Bold: true})
			titleLbl.Wrapping = fyne.TextWrapWord
			headerItems = append(headerItems, titleLbl)
		}

		if parts.Description != "" {
			descLbl := widget.NewLabel(parts.Description)
			descLbl.Wrapping = fyne.TextWrapWord
			headerItems = append(headerItems, descLbl)
		}

		headerItems = append(headerItems, widget.NewSeparator())
		elements = append(elements, container.NewVBox(headerItems...))
	}

	// Content section
	if parts.Content != nil {
		elements = append(elements, parts.Content)
	}

	// Footer section
	if parts.Footer != nil {
		if parts.Content != nil || (parts.Title == "" && parts.Description == "" && parts.Badge == nil) {
			elements = append(elements, widget.NewSeparator(), parts.Footer)
		} else {
			elements = append(elements, parts.Footer)
		}
	}

	return container.NewStack(
		cardBg,
		container.NewPadded(container.NewVBox(elements...)),
	)
}

var _ fyne.Widget = (*SelectableCard)(nil)
var _ fyne.Disableable = (*SelectableCard)(nil)
var _ fyne.Focusable = (*SelectableCard)(nil)
var _ desktop.Hoverable = (*SelectableCard)(nil)
var _ desktop.Cursorable = (*SelectableCard)(nil)
var _ desktop.Mouseable = (*SelectableCard)(nil)

// SelectableCard implements an interactive card container supporting
// Hover, Active/Pressed, Focus ring, Selected, and Disabled states.
type SelectableCard struct {
	widget.BaseWidget
	Selected  bool
	Hovered   bool
	Focused   bool
	Pressed   bool
	disabled  bool
	Content   fyne.CanvasObject
	OnTapped  func()
}

// NewSelectableCard creates an interactive selectable card.
func NewSelectableCard(content fyne.CanvasObject, selected bool, onTapped func()) *SelectableCard {
	c := &SelectableCard{
		Selected: selected,
		Content:  content,
		OnTapped: onTapped,
	}
	c.ExtendBaseWidget(c)
	return c
}

// SetSelected updates the selected state and refreshes visuals.
func (c *SelectableCard) SetSelected(selected bool) {
	c.Selected = selected
	c.Refresh()
}

// Enable enables interaction on this card.
func (c *SelectableCard) Enable() {
	c.disabled = false
	c.Refresh()
}

// Disable disables interaction on this card.
func (c *SelectableCard) Disable() {
	c.disabled = true
	c.Focused = false
	c.Pressed = false
	c.Refresh()
}

// Disabled reports whether the card is disabled.
func (c *SelectableCard) Disabled() bool {
	return c.disabled
}

// Tapped handles click selection.
func (c *SelectableCard) Tapped(_ *fyne.PointEvent) {
	if c.disabled {
		return
	}
	if c.OnTapped != nil {
		c.OnTapped()
	}
}

// MouseIn handles hover entrance.
func (c *SelectableCard) MouseIn(_ *desktop.MouseEvent) {
	if c.disabled {
		return
	}
	c.Hovered = true
	c.Refresh()
}

// MouseMoved handles hover movement.
func (c *SelectableCard) MouseMoved(_ *desktop.MouseEvent) {}

// MouseOut handles hover exit.
func (c *SelectableCard) MouseOut() {
	c.Hovered = false
	c.Refresh()
}

// MouseDown handles active press state.
func (c *SelectableCard) MouseDown(_ *desktop.MouseEvent) {
	if c.disabled {
		return
	}
	c.Pressed = true
	c.Refresh()
}

// MouseUp handles active press release.
func (c *SelectableCard) MouseUp(_ *desktop.MouseEvent) {
	c.Pressed = false
	c.Refresh()
}

// Cursor returns pointer when enabled, default when disabled.
func (c *SelectableCard) Cursor() desktop.Cursor {
	if c.disabled {
		return desktop.DefaultCursor
	}
	return desktop.PointerCursor
}

// FocusGained activates the keyboard focus ring.
func (c *SelectableCard) FocusGained() {
	if c.disabled {
		return
	}
	c.Focused = true
	c.Refresh()
}

// FocusLost deactivates the keyboard focus ring.
func (c *SelectableCard) FocusLost() {
	c.Focused = false
	c.Refresh()
}

// TypedRune handles spacebar key activation.
func (c *SelectableCard) TypedRune(r rune) {
	if c.disabled {
		return
	}
	if r == ' ' && c.OnTapped != nil {
		c.OnTapped()
	}
}

// TypedKey handles Enter / Return key activation.
func (c *SelectableCard) TypedKey(k *fyne.KeyEvent) {
	if c.disabled {
		return
	}
	if (k.Name == fyne.KeySpace || k.Name == fyne.KeyReturn || k.Name == fyne.KeyEnter) && c.OnTapped != nil {
		c.OnTapped()
	}
}

func (c *SelectableCard) CreateRenderer() fyne.WidgetRenderer {
	focusRing := canvas.NewRectangle(color.Transparent)
	focusRing.CornerRadius = 12
	focusRing.StrokeWidth = 2
	focusRing.StrokeColor = theme.Color(theme.ColorNameFocus)
	focusRing.Hide()

	cardBg := canvas.NewRectangle(theme.Color(theme.ColorNameMenuBackground))
	cardBg.CornerRadius = 10
	cardBg.StrokeColor = theme.Color(theme.ColorNameInputBorder)
	cardBg.StrokeWidth = 1

	objects := []fyne.CanvasObject{focusRing, cardBg}
	if c.Content != nil {
		objects = append(objects, c.Content)
	}

	r := &selectableCardRenderer{
		card:      c,
		focusRing: focusRing,
		cardBg:    cardBg,
		objects:   objects,
	}
	r.Refresh()
	return r
}

type selectableCardRenderer struct {
	card      *SelectableCard
	focusRing *canvas.Rectangle
	cardBg    *canvas.Rectangle
	objects   []fyne.CanvasObject
}

func (r *selectableCardRenderer) Layout(size fyne.Size) {
	r.focusRing.Resize(size)
	r.focusRing.Move(fyne.NewPos(0, 0))

	pad := float32(2)
	r.cardBg.Resize(fyne.NewSize(size.Width-pad*2, size.Height-pad*2))
	r.cardBg.Move(fyne.NewPos(pad, pad))

	if r.card.Content != nil {
		contentPad := float32(10)
		r.card.Content.Resize(fyne.NewSize(size.Width-contentPad*2, size.Height-contentPad*2))
		r.card.Content.Move(fyne.NewPos(contentPad, contentPad))
	}
}

func (r *selectableCardRenderer) MinSize() fyne.Size {
	if r.card.Content != nil {
		min := r.card.Content.MinSize()
		return fyne.NewSize(min.Width+24, min.Height+24)
	}
	return fyne.NewSize(120, 48)
}

func (r *selectableCardRenderer) Refresh() {
	if r.card.Focused && !r.card.disabled {
		r.focusRing.StrokeColor = theme.Color(theme.ColorNameFocus)
		r.focusRing.Show()
	} else {
		r.focusRing.Hide()
	}

	if r.card.disabled {
		r.cardBg.FillColor = color.NRGBA{R: 20, G: 26, B: 34, A: 140}
		r.cardBg.StrokeColor = color.NRGBA{R: 40, G: 48, B: 58, A: 100}
		r.cardBg.StrokeWidth = 1
	} else if r.card.Selected {
		r.cardBg.FillColor = color.NRGBA{R: 28, G: 34, B: 44, A: 255}
		r.cardBg.StrokeColor = theme.Color(theme.ColorNamePrimary)
		r.cardBg.StrokeWidth = 2
	} else if r.card.Pressed {
		r.cardBg.FillColor = color.NRGBA{R: 32, G: 42, B: 56, A: 255}
		r.cardBg.StrokeColor = color.NRGBA{R: 70, G: 84, B: 104, A: 255}
		r.cardBg.StrokeWidth = 1
	} else if r.card.Hovered {
		r.cardBg.FillColor = color.NRGBA{R: 28, G: 36, B: 46, A: 255}
		r.cardBg.StrokeColor = color.NRGBA{R: 64, G: 76, B: 96, A: 255}
		r.cardBg.StrokeWidth = 1
	} else {
		r.cardBg.FillColor = theme.Color(theme.ColorNameMenuBackground)
		r.cardBg.StrokeColor = theme.Color(theme.ColorNameInputBorder)
		r.cardBg.StrokeWidth = 1
	}

	r.Layout(r.card.Size())
	r.focusRing.Refresh()
	r.cardBg.Refresh()
	if r.card.Content != nil {
		r.card.Content.Refresh()
	}
}

func (r *selectableCardRenderer) Objects() []fyne.CanvasObject {
	return r.objects
}

func (r *selectableCardRenderer) Destroy() {}

