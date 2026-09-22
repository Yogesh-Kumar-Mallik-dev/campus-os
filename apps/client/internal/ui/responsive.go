package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"github.com/Yogesh-Kumar-Mallik-dev/campus-os/apps/client/internal/layout"
)

// NewResponsiveLayout dynamically centers content on wide desktop screens
// and fluidly expands to 100% width with touch-friendly margins on mobile viewports.
func NewResponsiveLayout(maxWidth, margin float32) fyne.Layout {
	if maxWidth <= 0 {
		maxWidth = 580
	}
	if margin <= 0 {
		margin = 12
	}
	return layout.NewFluidContainerLayout(layout.ContainerOptions{
		MinWidth:   240,
		MaxWidth:   maxWidth,
		Padding:    margin,
		AutoCenter: true,
	})
}

// NewResponsiveCardContainer wraps a CanvasObject in a fluidly adapting container.
func NewResponsiveCardContainer(content fyne.CanvasObject, maxWidth float32) fyne.CanvasObject {
	if maxWidth <= 0 {
		maxWidth = 580
	}
	inner := layout.Container(content, layout.ContainerOptions{
		MinWidth:      240,
		MaxWidth:      maxWidth,
		Padding:       12,
		AutoCenter:    true,
		ResponsivePad: true,
	})
	return container.NewVScroll(inner)
}

// NewAdaptiveGridLayout dynamically computes the optimal number of columns
// given the container's real-time width, ensuring children maintain at least minItemWidth.
func NewAdaptiveGridLayout(minItemWidth, gap float32) fyne.Layout {
	return layout.NewResponsiveGridLayout(layout.GridOptions{
		MinItemWidth:  minItemWidth,
		Gap:           gap,
		UniformHeight: true,
	})
}

// NewFlowLayout arranges children horizontally and wraps to subsequent rows when space is exceeded.
func NewFlowLayout(gap float32) fyne.Layout {
	return layout.NewFlowLayout(gap)
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
