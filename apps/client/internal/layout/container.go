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
	opts ContainerOptions
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

func (f *fluidContainerLayout) Layout(objects []fyne.CanvasObject, size fyne.Size) {
	pad := f.opts.Padding
	if f.opts.ResponsivePad {
		vp := ViewportFromSize(size, f.opts.Breakpoints, DefaultHeightBreakpoints())
		switch vp.SizeClass() {
		case Compact:
			pad = f.opts.Padding * 0.5
		case Expanded:
			pad = f.opts.Padding * 1.5
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
			targetW = f.opts.MinWidth
		}

		posX := pad
		if f.opts.AutoCenter && size.Width > targetW+(pad*2) {
			posX = (size.Width - targetW) / 2
		}

		availH := size.Height - (pad * 2)
		if availH < 0 {
			availH = 0
		}

		targetH := availH
		childMinH := child.MinSize().Height
		if childMinH > targetH {
			targetH = childMinH
		}
		if f.opts.MinHeight > 0 && targetH < f.opts.MinHeight {
			targetH = f.opts.MinHeight
		}
		if f.opts.MaxHeight > 0 && targetH > f.opts.MaxHeight {
			targetH = f.opts.MaxHeight
		}

		child.Move(fyne.NewPos(posX, pad))
		child.Resize(fyne.NewSize(targetW, targetH))
	}
}

func (f *fluidContainerLayout) MinSize(objects []fyne.CanvasObject) fyne.Size {
	pad := f.opts.Padding
	var minW float32 = 0
	var minH float32 = 0
	for _, child := range objects {
		if !child.Visible() {
			continue
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
	return fyne.NewSize(minW+(pad*2), minH+(pad*2))
}

// Container constructs a fluid container that caps max width and centers content on wide screens.
func Container(child fyne.CanvasObject, opts ContainerOptions) *fyne.Container {
	return container.New(NewFluidContainerLayout(opts), child)
}
