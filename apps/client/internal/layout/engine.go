package layout

import "fyne.io/fyne/v2"

// Engine coordinates global layout configuration, breakpoint policies, and responsive factory builders.
type Engine struct {
	Breakpoints       Breakpoints
	HeightBreakpoints HeightBreakpoints
	Spacing           Spacing
}

// NewEngine constructs a layout engine with industry standard defaults.
func NewEngine() *Engine {
	return &Engine{
		Breakpoints:       DefaultBreakpoints(),
		HeightBreakpoints: DefaultHeightBreakpoints(),
		Spacing:           DefaultSpacing(),
	}
}

// NewEngineWithConfig constructs a layout engine with custom rules.
func NewEngineWithConfig(bp Breakpoints, hbp HeightBreakpoints, sp Spacing) *Engine {
	if bp.Medium <= 0 || bp.Expanded <= 0 {
		bp = DefaultBreakpoints()
	}
	if hbp.Short <= 0 || hbp.Tall <= 0 {
		hbp = DefaultHeightBreakpoints()
	}
	return &Engine{
		Breakpoints:       bp,
		HeightBreakpoints: hbp,
		Spacing:           sp,
	}
}

// DefaultEngine is the global default responsive layout engine instance.
var DefaultEngine = NewEngine()

// Viewport constructs a Viewport abstraction from a fyne.Size using this engine's breakpoints.
func (e *Engine) Viewport(size fyne.Size) Viewport {
	return ViewportFromSize(size, e.Breakpoints, e.HeightBreakpoints)
}

// Container constructs a fluid, max-width capped, auto-centering container.
func (e *Engine) Container(child fyne.CanvasObject, opts ContainerOptions) *fyne.Container {
	if opts.Breakpoints.Medium <= 0 {
		opts.Breakpoints = e.Breakpoints
	}
	return Container(child, opts)
}

// Grid creates a responsive grid container that dynamically calculates column count from width.
func (e *Engine) Grid(opts GridOptions, objects ...fyne.CanvasObject) *fyne.Container {
	return Grid(opts, objects...)
}

// Flow creates a horizontal flow container wrapping items to next rows when space is exceeded.
func (e *Engine) Flow(gap float32, objects ...fyne.CanvasObject) *fyne.Container {
	return Flow(gap, objects...)
}

// Row creates an aligned horizontal flex container.
func (e *Engine) Row(gap float32, objects ...fyne.CanvasObject) *fyne.Container {
	return Row(gap, objects...)
}

// Column creates an aligned vertical flex container.
func (e *Engine) Column(gap float32, objects ...fyne.CanvasObject) *fyne.Container {
	return Column(gap, objects...)
}

// Stack creates an overlay container layering objects on top of each other.
func (e *Engine) Stack(objects ...fyne.CanvasObject) *fyne.Container {
	return Stack(objects...)
}

// Responsive creates a responsive container that switches views across SizeClasses.
func (e *Engine) Responsive(opts ResponsiveOptions) *ResponsiveContainer {
	if opts.Breakpoints.Medium <= 0 {
		opts.Breakpoints = e.Breakpoints
	}
	if opts.HeightBreakpoints.Short <= 0 {
		opts.HeightBreakpoints = e.HeightBreakpoints
	}
	return NewResponsive(opts)
}

// AppShell creates an adaptive application shell supporting Desktop, Tablet, and Mobile layouts.
func (e *Engine) AppShell(opts AppShellOptions) *ResponsiveAppShell {
	if opts.Breakpoints.Medium <= 0 {
		opts.Breakpoints = e.Breakpoints
	}
	return NewResponsiveAppShell(opts)
}
