package layout

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
)

// ContainerOptions configures fluid dimensional constraints and auto-centering rules.
type ContainerOptions struct {
	// MinWidth is the minimum width the child content should occupy (0 = unconstrained).
	MinWidth float32
	// PreferredWidth is the target nominal width when sufficient space is available.
	PreferredWidth float32
	// MaxWidth is the maximum width cap, preventing absurd stretching on ultrawide monitors (0 = unconstrained).
	MaxWidth float32
	// MinHeight is the minimum height the child content should occupy.
	MinHeight float32
	// MaxHeight is the maximum height cap for the container.
	MaxHeight float32
	// Padding is the base inner gutter around the content.
	Padding float32
	// AutoCenter automatically centers the child horizontally when available width exceeds MaxWidth.
	AutoCenter bool
	// ResponsivePad automatically adjusts padding based on SizeClass (Compact: 0.5x, Medium: 1x, Expanded: 1.5x).
	ResponsivePad bool
	// Breakpoints configures the horizontal size thresholds.
	Breakpoints Breakpoints
}

// fluidContainerLayout implements fyne.Layout for fluid, max-width bounded, auto-centering containers.
type fluidContainerLayout struct {
	opts        ContainerOptions
	lastTargetW float32
	container   *fyne.Container
}

// NewFluidContainerLayout constructs a new fluid container layout with the given options.
func NewFluidContainerLayout(opts ContainerOptions) fyne.Layout {
	if opts.Padding <= 0 {
		opts.Padding = 16
	}
	if opts.Breakpoints.Medium <= 0 {
		opts.Breakpoints = DefaultBreakpoints()
	}
	return &fluidContainerLayout{opts: opts}
}

func (f *fluidContainerLayout) effectivePadding(width float32) float32 {
	pad := f.opts.Padding
	if pad <= 0 {
		pad = 16
	}
	if !f.opts.ResponsivePad {
		return pad
	}
	bp := f.opts.Breakpoints
	if bp.Medium <= 0 {
		bp = DefaultBreakpoints()
	}
	if width > 0 {
		if width < bp.Medium {
			return pad * 0.5
		} else if width >= bp.Expanded {
			return pad * 1.5
		}
	}
	return pad
}

func (f *fluidContainerLayout) Layout(objects []fyne.CanvasObject, size fyne.Size) {
	if size.Width <= 0 || size.Height <= 0 {
		return
	}

	pad := f.effectivePadding(size.Width)
	if pad*2 >= size.Width {
		pad = size.Width / 4
		if pad < 0 {
			pad = 0
		}
	}

	for _, child := range objects {
		if !child.Visible() {
			continue
		}

		availW := size.Width - (pad * 2)
		if availW < 0 {
			availW = 0
		}

		targetW := availW
		if f.opts.PreferredWidth > 0 && targetW > f.opts.PreferredWidth {
			targetW = f.opts.PreferredWidth
		}
		if f.opts.MaxWidth > 0 && targetW > f.opts.MaxWidth {
			targetW = f.opts.MaxWidth
		}
		if f.opts.MinWidth > 0 && targetW < f.opts.MinWidth {
			if size.Width >= f.opts.MinWidth {
				targetW = f.opts.MinWidth
			} else {
				targetW = availW
			}
		}

		// Calculate symmetric auto-centering position, ensuring posX >= pad and never negative
		posX := pad
		if f.opts.AutoCenter && size.Width > targetW {
			centerOffset := (size.Width - targetW) / 2
			if centerOffset > pad {
				posX = centerOffset
			}
		}
		if posX < 0 {
			posX = 0
		}

		// Pre-resize child to targetW so text-wrapping widgets (labels, cards)
		// re-evaluate their accurate wrapped height at the allocated width.
		child.Resize(fyne.NewSize(targetW, child.MinSize().Height))
		childMinH := child.MinSize().Height

		availH := size.Height - (pad * 2)
		if availH < 0 {
			availH = 0
		}

		targetH := availH
		if childMinH > targetH {
			targetH = childMinH
		}
		if f.opts.MinHeight > 0 && targetH < f.opts.MinHeight {
			targetH = f.opts.MinHeight
		}
		if f.opts.MaxHeight > 0 && targetH > f.opts.MaxHeight {
			targetH = f.opts.MaxHeight
		}

		f.lastTargetW = targetW

		child.Move(fyne.NewPos(posX, pad))
		child.Resize(fyne.NewSize(targetW, targetH))
	}
}

func (f *fluidContainerLayout) MinSize(objects []fyne.CanvasObject) fyne.Size {
	var minW float32 = 0
	var minH float32 = 0
	for _, child := range objects {
		if !child.Visible() {
			continue
		}
		if f.lastTargetW > 0 {
			child.Resize(fyne.NewSize(f.lastTargetW, child.MinSize().Height))
		}
		ms := child.MinSize()
		if ms.Width > minW {
			minW = ms.Width
		}
		if ms.Height > minH {
			minH = ms.Height
		}
	}
	if f.opts.MinWidth > 0 && minW < f.opts.MinWidth {
		minW = f.opts.MinWidth
	}
	if f.opts.MinHeight > 0 && minH < f.opts.MinHeight {
		minH = f.opts.MinHeight
	}

	padWidth := minW
	if f.lastTargetW > 0 {
		padWidth = f.lastTargetW
	}
	pad := f.effectivePadding(padWidth)

	reportedMinW := minW
	if f.opts.AutoCenter && f.opts.MinWidth > 0 {
		reportedMinW = f.opts.MinWidth
	}

	return fyne.NewSize(reportedMinW+(pad*2), minH+(pad*2)+24)
}

// Container constructs a fluid container that caps max width and centers content on wide screens.
func Container(child fyne.CanvasObject, opts ContainerOptions) *fyne.Container {
	l := NewFluidContainerLayout(opts)
	c := container.New(l, child)
	if fl, ok := l.(*fluidContainerLayout); ok {
		fl.container = c
	}
	return c
}

// FixedWidthLayout restricts a child to a specific width while allowing it to fluidly expand vertically.
type fixedWidthLayout struct {
	width float32
}

// NewFixedWidthLayout creates a layout that forces a fixed width and fluid height.
func NewFixedWidthLayout(width float32) fyne.Layout {
	return &fixedWidthLayout{width: width}
}

func (l *fixedWidthLayout) Layout(objects []fyne.CanvasObject, size fyne.Size) {
	for _, o := range objects {
		o.Move(fyne.NewPos(0, 0))
		o.Resize(fyne.NewSize(l.width, size.Height))
	}
}

func (l *fixedWidthLayout) MinSize(objects []fyne.CanvasObject) fyne.Size {
	var maxH float32 = 0
	for _, o := range objects {
		if o.Visible() && o.MinSize().Height > maxH {
			maxH = o.MinSize().Height
		}
	}
	return fyne.NewSize(l.width, maxH)
}

// FixedWidth wraps a canvas object with a fixed width container that fills 100% vertical height.
func FixedWidth(width float32, child fyne.CanvasObject) *fyne.Container {
	return container.New(NewFixedWidthLayout(width), child)
}
