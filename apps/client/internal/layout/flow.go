package layout

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
)

// FlowOptions configures horizontal and vertical wrapping gaps.
type FlowOptions struct {
	Gap    float32
	RowGap float32
}

// responsiveFlowLayout arranges children horizontally and wraps to subsequent rows when space is exceeded.
type responsiveFlowLayout struct {
	opts      FlowOptions
	lastWidth float32
	container *fyne.Container
}

// NewFlowLayout creates a flow layout with equal horizontal and vertical gap.
func NewFlowLayout(gap float32) fyne.Layout {
	if gap <= 0 {
		gap = 8
	}
	return &responsiveFlowLayout{opts: FlowOptions{Gap: gap, RowGap: gap}}
}

// NewFlowLayoutWithOptions creates a flow layout with independent gaps.
func NewFlowLayoutWithOptions(opts FlowOptions) fyne.Layout {
	if opts.Gap <= 0 {
		opts.Gap = 8
	}
	if opts.RowGap <= 0 {
		opts.RowGap = opts.Gap
	}
	return &responsiveFlowLayout{opts: opts}
}

func (f *responsiveFlowLayout) Layout(objects []fyne.CanvasObject, size fyne.Size) {
	x := float32(0)
	y := float32(0)
	rowH := float32(0)
	rows := 1

	for _, o := range objects {
		if !o.Visible() {
			continue
		}
		ms := o.MinSize()
		if x+ms.Width > size.Width && x > 0 {
			// Wrap to next line
			x = 0
			y += rowH + f.opts.RowGap
			rowH = 0
			rows++
		}
		o.Move(fyne.NewPos(x, y))
		o.Resize(ms)
		x += ms.Width + f.opts.Gap
		if ms.Height > rowH {
			rowH = ms.Height
		}
	}

	totalH := y + rowH
	widthChanged := f.lastWidth != size.Width
	f.lastWidth = size.Width

	if widthChanged && f.container != nil && totalH > 0 {
		fyne.Do(func() {
			f.container.Refresh()
		})
	}
}

func (f *responsiveFlowLayout) MinSize(objects []fyne.CanvasObject) fyne.Size {
	minW := float32(0)
	maxChildH := float32(0)
	for _, o := range objects {
		if !o.Visible() {
			continue
		}
		ms := o.MinSize()
		if ms.Width > minW {
			minW = ms.Width
		}
		if ms.Height > maxChildH {
			maxChildH = ms.Height
		}
	}

	if f.lastWidth <= 0 {
		return fyne.NewSize(minW, maxChildH)
	}

	// Calculate wrapped height based on last known container width
	x := float32(0)
	y := float32(0)
	rowH := float32(0)
	for _, o := range objects {
		if !o.Visible() {
			continue
		}
		ms := o.MinSize()
		if x+ms.Width > f.lastWidth && x > 0 {
			x = 0
			y += rowH + f.opts.RowGap
			rowH = 0
		}
		x += ms.Width + f.opts.Gap
		if ms.Height > rowH {
			rowH = ms.Height
		}
	}
	totalH := y + rowH
	if totalH < maxChildH {
		totalH = maxChildH
	}
	return fyne.NewSize(minW, totalH)
}

// Flow creates a flow container wrapping items when row width is exceeded.
func Flow(gap float32, objects ...fyne.CanvasObject) *fyne.Container {
	l := NewFlowLayout(gap)
	c := container.New(l, objects...)
	if fl, ok := l.(*responsiveFlowLayout); ok {
		fl.container = c
	}
	return c
}
