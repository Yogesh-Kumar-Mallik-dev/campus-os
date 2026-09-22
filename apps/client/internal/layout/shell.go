package layout

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
)

// AppShellOptions configures components for the multi-mode adaptive application shell.
type AppShellOptions struct {
	Header      fyne.CanvasObject
	Sidebar     fyne.CanvasObject
	Content     fyne.CanvasObject
	BottomNav   fyne.CanvasObject
	Drawer      fyne.CanvasObject
	Breakpoints Breakpoints

	// SidebarWidth sets the fixed desktop sidebar width (default: 220dp).
	SidebarWidth float32
}

// ResponsiveAppShell manages seamless structural adaptation between Desktop, Tablet, and Mobile shells.
type ResponsiveAppShell struct {
	opts         AppShellOptions
	container    *ResponsiveContainer
	drawerActive bool
	rootStack    *fyne.Container
}

// NewResponsiveAppShell constructs an adaptive application shell adhering to Design System 2026.
func NewResponsiveAppShell(opts AppShellOptions) *ResponsiveAppShell {
	if opts.SidebarWidth <= 0 {
		opts.SidebarWidth = 220
	}
	if opts.Breakpoints.Medium <= 0 {
		opts.Breakpoints = DefaultBreakpoints()
	}

	shell := &ResponsiveAppShell{
		opts: opts,
	}

	// 1. Build Expanded (Desktop) View:
	// Header at top, Fixed Sidebar at left, Content in center.
	expandedView := container.NewBorder(
		opts.Header,
		nil,
		opts.Sidebar,
		nil,
		opts.Content,
	)

	// 2. Build Medium (Tablet / Split-screen) View:
	// Header at top, Content in center (Sidebar collapses into toggleable drawer).
	mediumView := container.NewBorder(
		opts.Header,
		nil,
		nil,
		nil,
		opts.Content,
	)

	// 3. Build Compact (Mobile) View:
	// Header at top, BottomNav at bottom, Content in center.
	compactView := container.NewBorder(
		opts.Header,
		opts.BottomNav,
		nil,
		nil,
		opts.Content,
	)

	responsiveContainer := NewResponsive(ResponsiveOptions{
		Breakpoints: opts.Breakpoints,
		Expanded:    expandedView,
		Medium:      mediumView,
		Compact:     compactView,
	})

	shell.container = responsiveContainer

	// If a slide-in drawer is provided, wrap in an overlay Stack
	if opts.Drawer != nil {
		opts.Drawer.Hide()
		shell.rootStack = Stack(responsiveContainer, opts.Drawer)
	} else {
		shell.rootStack = Stack(responsiveContainer)
	}

	return shell
}

// Root returns the root canvas object for window embedding.
func (s *ResponsiveAppShell) Root() fyne.CanvasObject {
	return s.rootStack
}

// ToggleDrawer toggles visibility of the mobile slide-in drawer.
func (s *ResponsiveAppShell) ToggleDrawer() {
	if s.opts.Drawer == nil {
		return
	}
	if s.drawerActive {
		s.opts.Drawer.Hide()
		s.drawerActive = false
	} else {
		s.opts.Drawer.Show()
		s.drawerActive = true
	}
	s.rootStack.Refresh()
}

// CloseDrawer closes the slide-in drawer if currently open.
func (s *ResponsiveAppShell) CloseDrawer() {
	if s.opts.Drawer != nil && s.drawerActive {
		s.opts.Drawer.Hide()
		s.drawerActive = false
		s.rootStack.Refresh()
	}
}

// CurrentSizeClass returns the active shell SizeClass.
func (s *ResponsiveAppShell) CurrentSizeClass() SizeClass {
	return s.container.CurrentSizeClass()
}

// CurrentViewport returns the measured active Viewport.
func (s *ResponsiveAppShell) CurrentViewport() Viewport {
	return s.container.CurrentViewport()
}
