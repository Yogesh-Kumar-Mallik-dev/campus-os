package layout

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
)

// GridOptions configures fluid column derivation and item gap spacing.
type GridOptions struct {
	// MinItemWidth is the minimum width a single item requires.
	MinItemWidth float32
	// MaxItemWidth optionally caps the maximum width an item can expand to (0 = unconstrained).
	MaxItemWidth float32
	// MaxCols optionally sets an upper limit on the number of generated columns (0 = unconstrained).
	MaxCols int
	// Gap sets horizontal spacing between columns.
	Gap float32
	// RowGap sets vertical spacing between rows (defaults to Gap if <= 0).
	RowGap float32
	// UniformHeight forces all items in a row to match the tallest item's height.
	UniformHeight bool
}

// responsiveGridLayout dynamically derives column count from available width.
type responsiveGridLayout struct {
	opts     GridOptions
	lastCols int
}

// NewResponsiveGridLayout constructs a fyne.Layout for adaptive grid arrangements.
func NewResponsiveGridLayout(opts GridOptions) fyne.Layout {
	if opts.MinItemWidth <= 0 {
		opts.MinItemWidth = 240
	}
	if opts.Gap <= 0 {
		opts.Gap = 16
	}
	if opts.RowGap <= 0 {
		opts.RowGap = opts.Gap
	}
	return &responsiveGridLayout{opts: opts, lastCols: 1}
}

func (g *responsiveGridLayout) Layout(objects []fyne.CanvasObject, size fyne.Size) {
	visible := make([]fyne.CanvasObject, 0, len(objects))
	for _, o := range objects {
		if o.Visible() {
			visible = append(visible, o)
		}
	}
	if len(visible) == 0 {
		return
	}

	cols := Columns(size.Width, g.opts.MinItemWidth, g.opts.Gap)
	if g.opts.MaxCols > 0 && cols > g.opts.MaxCols {
		cols = g.opts.MaxCols
	}
	if cols > len(visible) {
		cols = len(visible)
	}
	if cols < 1 {
		cols = 1
	}
	g.lastCols = cols

	colW := (size.Width - float32(cols-1)*g.opts.Gap) / float32(cols)
	if g.opts.MaxItemWidth > 0 && colW > g.opts.MaxItemWidth {
		colW = g.opts.MaxItemWidth
	}

	// Calculate row count and heights
	rows := (len(visible) + cols - 1) / cols
	rowHeights := make([]float32, rows)
	maxRowH := float32(0)

	for i, o := range visible {
		r := i / cols
		h := o.MinSize().Height
		if h > rowHeights[r] {
			rowHeights[r] = h
		}
		if h > maxRowH {
			maxRowH = h
		}
	}

	// Layout items
	rowY := float32(0)
	for r := 0; r < rows; r++ {
		currentH := rowHeights[r]
		if g.opts.UniformHeight {
			currentH = maxRowH
		}
		if currentH < 30 {
			currentH = 30
		}
		for c := 0; c < cols; c++ {
			idx := r*cols + c
			if idx >= len(visible) {
				break
			}
			o := visible[idx]
			x := float32(c) * (colW + g.opts.Gap)
			o.Move(fyne.NewPos(x, rowY))
			o.Resize(fyne.NewSize(colW, currentH))
		}
		rowY += currentH + g.opts.RowGap
	}
}

func (g *responsiveGridLayout) MinSize(objects []fyne.CanvasObject) fyne.Size {
	visible := 0
	maxChildH := float32(0)
	for _, o := range objects {
		if o.Visible() {
			visible++
			if h := o.MinSize().Height; h > maxChildH {
				maxChildH = h
			}
		}
	}
	if visible == 0 {
		return fyne.NewSize(0, 0)
	}
	cols := g.lastCols
	if cols < 1 {
		cols = 1
	}
	rows := (visible + cols - 1) / cols
	totalH := float32(rows)*maxChildH + float32(rows-1)*g.opts.RowGap
	return fyne.NewSize(g.opts.MinItemWidth, totalH)
}

// Grid creates a responsive grid container that dynamically calculates column count from available width.
func Grid(opts GridOptions, objects ...fyne.CanvasObject) *fyne.Container {
	return container.New(NewResponsiveGridLayout(opts), objects...)
}
