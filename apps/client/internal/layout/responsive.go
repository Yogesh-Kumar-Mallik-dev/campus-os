package layout

import (
	"sync"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
)

// ResponsiveOptions configures static views, dynamic builders, and lifecycle hooks for responsive layouts.
type ResponsiveOptions struct {
	Breakpoints       Breakpoints
	HeightBreakpoints HeightBreakpoints

	// Static pre-built views:
	Compact  fyne.CanvasObject
	Medium   fyne.CanvasObject
	Expanded fyne.CanvasObject

	// Dynamic lazy builders (invoked once per SizeClass and cached across resizes):
	CompactBuilder  func(vp Viewport) fyne.CanvasObject
	MediumBuilder   func(vp Viewport) fyne.CanvasObject
	ExpandedBuilder func(vp Viewport) fyne.CanvasObject

	// OnChange is an optional notification hook fired when SizeClass transitions occur.
	OnChange func(from, to SizeClass, vp Viewport)
}

// ResponsiveContainer is a smart widget that fluidly resizes within a SizeClass,
// and only executes structural transitions when the SizeClass actually changes.
type ResponsiveContainer struct {
	widget.BaseWidget

	opts             ResponsiveOptions
	mu               sync.Mutex
	lastSizeClass    SizeClass
	hasLastSizeClass bool
	lastViewport     Viewport

	currentChild  fyne.CanvasObject
	compactCache  fyne.CanvasObject
	mediumCache   fyne.CanvasObject
	expandedCache fyne.CanvasObject

	onSizeChanged func(size fyne.Size)
}

// NewResponsive constructs a responsive container that automatically adapts to the current viewport.
func NewResponsive(opts ResponsiveOptions) *ResponsiveContainer {
	if opts.Breakpoints.Medium <= 0 || opts.Breakpoints.Expanded <= 0 {
		opts.Breakpoints = DefaultBreakpoints()
	}
	if opts.HeightBreakpoints.Short <= 0 || opts.HeightBreakpoints.Tall <= 0 {
		opts.HeightBreakpoints = DefaultHeightBreakpoints()
	}

	rc := &ResponsiveContainer{
		opts:          opts,
		compactCache:  opts.Compact,
		mediumCache:   opts.Medium,
		expandedCache: opts.Expanded,
	}
	rc.ExtendBaseWidget(rc)
	return rc
}

// Responsive creates a responsive container from options (alias constructor).
func Responsive(opts ResponsiveOptions) *ResponsiveContainer {
	return NewResponsive(opts)
}

// CurrentSizeClass returns the active SizeClass.
func (r *ResponsiveContainer) CurrentSizeClass() SizeClass {
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.lastSizeClass
}

// CurrentViewport returns the most recent Viewport geometry measured.
func (r *ResponsiveContainer) CurrentViewport() Viewport {
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.lastViewport
}

// SetOnSizeChanged registers an optional callback invoked on every layout resize event.
func (r *ResponsiveContainer) SetOnSizeChanged(fn func(size fyne.Size)) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.onSizeChanged = fn
}

func (r *ResponsiveContainer) resolveChildFor(sc SizeClass, vp Viewport) fyne.CanvasObject {
	switch sc {
	case Compact:
		if r.compactCache != nil {
			return r.compactCache
		}
		if r.opts.CompactBuilder != nil {
			r.compactCache = r.opts.CompactBuilder(vp)
			return r.compactCache
		}
		// Fallbacks:
		if r.mediumCache != nil {
			return r.mediumCache
		}
		return r.expandedCache

	case Medium:
		if r.mediumCache != nil {
			return r.mediumCache
		}
		if r.opts.MediumBuilder != nil {
			r.mediumCache = r.opts.MediumBuilder(vp)
			return r.mediumCache
		}
		// Fallbacks:
		if r.compactCache != nil {
			return r.compactCache
		}
		return r.expandedCache

	case Expanded:
		if r.expandedCache != nil {
			return r.expandedCache
		}
		if r.opts.ExpandedBuilder != nil {
			r.expandedCache = r.opts.ExpandedBuilder(vp)
			return r.expandedCache
		}
		// Fallbacks:
		if r.mediumCache != nil {
			return r.mediumCache
		}
		return r.compactCache

	default:
		return r.compactCache
	}
}

// CreateRenderer creates the widget renderer for the responsive container.
func (r *ResponsiveContainer) CreateRenderer() fyne.WidgetRenderer {
	return &responsiveRenderer{container: r}
}

type responsiveRenderer struct {
	container *ResponsiveContainer
}

func (rr *responsiveRenderer) Destroy() {}

func (rr *responsiveRenderer) Objects() []fyne.CanvasObject {
	rr.container.mu.Lock()
	defer rr.container.mu.Unlock()
	if rr.container.currentChild != nil {
		return []fyne.CanvasObject{rr.container.currentChild}
	}
	return nil
}

func (rr *responsiveRenderer) Refresh() {
	rr.container.mu.Lock()
	child := rr.container.currentChild
	rr.container.mu.Unlock()
	if child != nil {
		child.Refresh()
	}
}

func (rr *responsiveRenderer) MinSize() fyne.Size {
	rr.container.mu.Lock()
	child := rr.container.currentChild
	rr.container.mu.Unlock()
	if child != nil {
		return child.MinSize()
	}
	return fyne.NewSize(240, 160)
}

func (rr *responsiveRenderer) Layout(size fyne.Size) {
	if size.Width <= 0 || size.Height <= 0 {
		return
	}

	rr.container.mu.Lock()
	vp := ViewportFromSize(size, rr.container.opts.Breakpoints, rr.container.opts.HeightBreakpoints)
	newSizeClass := vp.SizeClass()

	sizeClassChanged := !rr.container.hasLastSizeClass || newSizeClass != rr.container.lastSizeClass
	oldSizeClass := rr.container.lastSizeClass

	rr.container.lastViewport = vp
	rr.container.lastSizeClass = newSizeClass
	rr.container.hasLastSizeClass = true

	if sizeClassChanged {
		selected := rr.container.resolveChildFor(newSizeClass, vp)
		rr.container.currentChild = selected

		if rr.container.opts.OnChange != nil {
			onChange := rr.container.opts.OnChange
			rr.container.mu.Unlock()
			onChange(oldSizeClass, newSizeClass, vp)
			rr.container.mu.Lock()
		}
	}

	activeChild := rr.container.currentChild
	onResize := rr.container.onSizeChanged
	rr.container.mu.Unlock()

	if activeChild != nil {
		activeChild.Move(fyne.NewPos(0, 0))
		activeChild.Resize(size)
	}

	if onResize != nil {
		onResize(size)
	}
}
