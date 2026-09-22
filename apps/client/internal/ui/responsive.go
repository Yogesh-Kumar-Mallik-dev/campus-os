package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
)

// responsiveCardLayout dynamically centers content on wide desktop screens
// and fluidly expands to 100% width with touch-friendly margins on mobile viewports.
type responsiveCardLayout struct {
	maxWidth float32
	margin   float32
}

func NewResponsiveLayout(maxWidth, margin float32) fyne.Layout {
	if maxWidth <= 0 {
		maxWidth = 580
	}
	if margin <= 0 {
		margin = 12
	}
	return &responsiveCardLayout{
		maxWidth: maxWidth,
		margin:   margin,
	}
}

func (r *responsiveCardLayout) Layout(objects []fyne.CanvasObject, size fyne.Size) {
	for _, child := range objects {
		var targetW float32
		if size.Width <= 0 {
			targetW = r.maxWidth
		} else {
			targetW = size.Width - (r.margin * 2)
			if targetW > r.maxWidth {
				targetW = r.maxWidth
			}
			if targetW < 240 {
				targetW = 240
			}
		}

		var posX float32 = r.margin
		if size.Width > targetW {
			posX = (size.Width - targetW) / 2
		}

		// Ensure child is given targetW before computing wrapped height
		child.Resize(fyne.NewSize(targetW, child.MinSize().Height))
		contentH := child.MinSize().Height + 28 // 28px bottom breathing margin
		if contentH < size.Height {
			contentH = size.Height
		}

		child.Move(fyne.NewPos(posX, 8))
		child.Resize(fyne.NewSize(targetW, contentH))
	}
}

func (r *responsiveCardLayout) MinSize(objects []fyne.CanvasObject) fyne.Size {
	var minW float32 = 280
	var maxMinH float32 = 0

	for _, child := range objects {
		ms := child.MinSize()
		if ms.Height > maxMinH {
			maxMinH = ms.Height
		}
	}

	return fyne.NewSize(minW, maxMinH+28)
}

// NewResponsiveCardContainer wraps a CanvasObject in a fluidly adapting container.
func NewResponsiveCardContainer(content fyne.CanvasObject, maxWidth float32) fyne.CanvasObject {
	inner := container.New(NewResponsiveLayout(maxWidth, 12), content)
	return container.NewVScroll(inner)
}

// AdaptiveGridLayout dynamically computes the optimal number of columns
// given the container's real-time width, ensuring children maintain at least minItemWidth.
type adaptiveGridLayout struct {
	minItemWidth float32
	gap          float32
	lastCols     int
}

func NewAdaptiveGridLayout(minItemWidth, gap float32) fyne.Layout {
	if minItemWidth <= 0 {
		minItemWidth = 220
	}
	if gap <= 0 {
		gap = 8
	}
	return &adaptiveGridLayout{
		minItemWidth: minItemWidth,
		gap:          gap,
		lastCols:     1,
	}
}

func (a *adaptiveGridLayout) Layout(objects []fyne.CanvasObject, size fyne.Size) {
	visible := make([]fyne.CanvasObject, 0, len(objects))
	for _, o := range objects {
		if o.Visible() {
			visible = append(visible, o)
		}
	}
	if len(visible) == 0 {
		return
	}

	cols := int((size.Width + a.gap) / (a.minItemWidth + a.gap))
	if cols < 1 {
		cols = 1
	}
	if cols > len(visible) {
		cols = len(visible)
	}
	a.lastCols = cols

	colW := (size.Width - float32(cols-1)*a.gap) / float32(cols)

	// Measure max row height for visual consistency
	rowH := float32(0)
	for _, o := range visible {
		if h := o.MinSize().Height; h > rowH {
			rowH = h
		}
	}
	if rowH < 50 {
		rowH = 50
	}

	for i, o := range visible {
		c := i % cols
		r := i / cols
		x := float32(c) * (colW + a.gap)
		y := float32(r) * (rowH + a.gap)
		o.Move(fyne.NewPos(x, y))
		o.Resize(fyne.NewSize(colW, rowH))
	}
}

func (a *adaptiveGridLayout) MinSize(objects []fyne.CanvasObject) fyne.Size {
	visible := 0
	maxH := float32(0)
	for _, o := range objects {
		if o.Visible() {
			visible++
			if h := o.MinSize().Height; h > maxH {
				maxH = h
			}
		}
	}
	if visible == 0 {
		return fyne.NewSize(0, 0)
	}
	cols := a.lastCols
	if cols < 1 {
		cols = 1
	}
	rows := (visible + cols - 1) / cols
	totalH := float32(rows)*maxH + float32(rows-1)*a.gap
	return fyne.NewSize(a.minItemWidth, totalH)
}

// FlowLayout arranges children horizontally and wraps to subsequent rows when space is exceeded.
type flowLayout struct {
	gap float32
}

func NewFlowLayout(gap float32) fyne.Layout {
	if gap <= 0 {
		gap = 8
	}
	return &flowLayout{gap: gap}
}

func (f *flowLayout) Layout(objects []fyne.CanvasObject, size fyne.Size) {
	x := float32(0)
	y := float32(0)
	rowH := float32(0)

	for _, o := range objects {
		if !o.Visible() {
			continue
		}
		ms := o.MinSize()
		if x+ms.Width > size.Width && x > 0 {
			// Wrap to next row
			x = 0
			y += rowH + f.gap
			rowH = 0
		}
		o.Move(fyne.NewPos(x, y))
		o.Resize(ms)
		x += ms.Width + f.gap
		if ms.Height > rowH {
			rowH = ms.Height
		}
	}
}

func (f *flowLayout) MinSize(objects []fyne.CanvasObject) fyne.Size {
	minW := float32(0)
	minH := float32(0)
	for _, o := range objects {
		if !o.Visible() {
			continue
		}
		ms := o.MinSize()
		if ms.Width > minW {
			minW = ms.Width
		}
		if ms.Height > minH {
			minH = ms.Height
		}
	}
	return fyne.NewSize(minW, minH)
}

// ResizeNotifierLayout wraps canvas objects and invokes onResize whenever container dimensions change.
type ResizeNotifierLayout struct {
	onResize func(size fyne.Size)
	lastSize fyne.Size
}

func NewResizeNotifierLayout(onResize func(size fyne.Size)) fyne.Layout {
	return &ResizeNotifierLayout{onResize: onResize}
}

func (l *ResizeNotifierLayout) Layout(objects []fyne.CanvasObject, size fyne.Size) {
	for _, o := range objects {
		o.Move(fyne.NewPos(0, 0))
		o.Resize(size)
	}
	if size.Width > 0 && int(size.Width) != int(l.lastSize.Width) {
		l.lastSize = size
		if l.onResize != nil {
			fyne.Do(func() {
				l.onResize(size)
			})
		}
	}
}

func (l *ResizeNotifierLayout) MinSize(objects []fyne.CanvasObject) fyne.Size {
	if len(objects) > 0 {
		return objects[0].MinSize()
	}
	return fyne.NewSize(0, 0)
}

