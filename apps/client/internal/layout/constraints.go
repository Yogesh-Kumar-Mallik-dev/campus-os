package layout

import "fyne.io/fyne/v2"

// Constraints specifies dimension bounding rules, padding, and gaps for components.
type Constraints struct {
	MinWidth  float32
	MaxWidth  float32
	MinHeight float32
	MaxHeight float32

	Padding float32
	Gap     float32
}

// ClampWidth bounds a proposed width within MinWidth and MaxWidth if set (> 0).
func (c Constraints) ClampWidth(w float32) float32 {
	if c.MinWidth > 0 && w < c.MinWidth {
		w = c.MinWidth
	}
	if c.MaxWidth > 0 && w > c.MaxWidth {
		w = c.MaxWidth
	}
	return w
}

// ClampHeight bounds a proposed height within MinHeight and MaxHeight if set (> 0).
func (c Constraints) ClampHeight(h float32) float32 {
	if c.MinHeight > 0 && h < c.MinHeight {
		h = c.MinHeight
	}
	if c.MaxHeight > 0 && h > c.MaxHeight {
		h = c.MaxHeight
	}
	return h
}

// ClampSize bounds both width and height of a fyne.Size against these constraints.
func (c Constraints) ClampSize(s fyne.Size) fyne.Size {
	return fyne.NewSize(c.ClampWidth(s.Width), c.ClampHeight(s.Height))
}

// EffectiveWidth returns the usable width after subtracting horizontal padding.
func (c Constraints) EffectiveWidth(availW float32) float32 {
	w := availW - (c.Padding * 2)
	if w < 0 {
		w = 0
	}
	return c.ClampWidth(w)
}

// EffectiveHeight returns the usable height after subtracting vertical padding.
func (c Constraints) EffectiveHeight(availH float32) float32 {
	h := availH - (c.Padding * 2)
	if h < 0 {
		h = 0
	}
	return c.ClampHeight(h)
}
