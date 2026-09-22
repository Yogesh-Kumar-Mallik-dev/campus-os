package layout

import (
	"fmt"

	"fyne.io/fyne/v2"
)

// SizeClass categorizes the horizontal dimension of the available viewport.
type SizeClass int

const (
	// Compact represents small mobile screens or narrow window widths (< 600dp by default).
	Compact SizeClass = iota
	// Medium represents tablets, split-screen desktop windows, or intermediate widths (600dp..900dp by default).
	Medium
	// Expanded represents desktop monitors, full-screen windows, and ultrawide displays (>= 900dp by default).
	Expanded
)

// String returns human-readable representation of SizeClass.
func (s SizeClass) String() string {
	switch s {
	case Compact:
		return "Compact"
	case Medium:
		return "Medium"
	case Expanded:
		return "Expanded"
	default:
		return fmt.Sprintf("SizeClass(%d)", s)
	}
}

// HeightClass categorizes the vertical dimension of the available viewport,
// allowing the engine to adapt to short or extreme aspect-ratio windows.
type HeightClass int

const (
	// Short represents shallow windows or mobile landscape mode (< 450dp by default).
	Short HeightClass = iota
	// Normal represents typical vertical screen real estate (450dp..900dp by default).
	Normal
	// Tall represents portrait screens or large desktop monitors (>= 900dp by default).
	Tall
)

// String returns human-readable representation of HeightClass.
func (h HeightClass) String() string {
	switch h {
	case Short:
		return "Short"
	case Normal:
		return "Normal"
	case Tall:
		return "Tall"
	default:
		return fmt.Sprintf("HeightClass(%d)", h)
	}
}

// Orientation represents the geometric aspect ratio profile of the viewport.
type Orientation int

const (
	// Portrait represents a viewport where height > width.
	Portrait Orientation = iota
	// Landscape represents a viewport where width >= height.
	Landscape
)

// String returns human-readable representation of Orientation.
func (o Orientation) String() string {
	switch o {
	case Portrait:
		return "Portrait"
	case Landscape:
		return "Landscape"
	default:
		return fmt.Sprintf("Orientation(%d)", o)
	}
}

// Viewport provides a device-agnostic, geometry-driven abstraction
// of the available screen, canvas, or window dimensions.
type Viewport struct {
	Width             float32
	Height            float32
	Breakpoints       Breakpoints
	HeightBreakpoints HeightBreakpoints
}

// NewViewport constructs a Viewport with specified dimensions and breakpoint rules.
func NewViewport(width, height float32, bp Breakpoints, hbp HeightBreakpoints) Viewport {
	if bp.Medium <= 0 || bp.Expanded <= 0 {
		bp = DefaultBreakpoints()
	}
	if hbp.Short <= 0 || hbp.Tall <= 0 {
		hbp = DefaultHeightBreakpoints()
	}
	return Viewport{
		Width:             width,
		Height:            height,
		Breakpoints:       bp,
		HeightBreakpoints: hbp,
	}
}

// ViewportFromSize constructs a Viewport from a fyne.Size.
func ViewportFromSize(size fyne.Size, bp Breakpoints, hbp HeightBreakpoints) Viewport {
	return NewViewport(size.Width, size.Height, bp, hbp)
}

// AspectRatio computes the width-to-height ratio (width / height).
func (v Viewport) AspectRatio() float32 {
	if v.Height <= 0 {
		return 1.0
	}
	return v.Width / v.Height
}

// Orientation evaluates whether the viewport is Portrait or Landscape.
func (v Viewport) Orientation() Orientation {
	if v.Height > v.Width {
		return Portrait
	}
	return Landscape
}

// SizeClass computes the horizontal sizing category according to configured Breakpoints.
func (v Viewport) SizeClass() SizeClass {
	if v.Width < v.Breakpoints.Medium {
		return Compact
	}
	if v.Width < v.Breakpoints.Expanded {
		return Medium
	}
	return Expanded
}

// HeightClass computes the vertical sizing category according to configured HeightBreakpoints.
func (v Viewport) HeightClass() HeightClass {
	if v.Height < v.HeightBreakpoints.Short {
		return Short
	}
	if v.Height < v.HeightBreakpoints.Tall {
		return Normal
	}
	return Tall
}

// Size returns the fyne.Size equivalent of this viewport.
func (v Viewport) Size() fyne.Size {
	return fyne.NewSize(v.Width, v.Height)
}

// IsUltrawide returns true if aspect ratio exceeds standard widescreen (typically > 2.0).
func (v Viewport) IsUltrawide() bool {
	return v.AspectRatio() >= 2.1
}

// String provides formatted diagnostic summary of the current viewport state.
func (v Viewport) String() string {
	return fmt.Sprintf("%.0fx%.0f (%s, %s, %s, Aspect: %.2f)",
		v.Width, v.Height, v.SizeClass(), v.HeightClass(), v.Orientation(), v.AspectRatio())
}
