package layout

// Breakpoints configures the horizontal thresholds separating Compact, Medium, and Expanded viewports.
type Breakpoints struct {
	// Medium is the horizontal threshold separating Compact from Medium (default: 600dp).
	Medium float32
	// Expanded is the horizontal threshold separating Medium from Expanded (default: 900dp).
	Expanded float32
}

// DefaultBreakpoints returns industry standard horizontal responsive breakpoints.
func DefaultBreakpoints() Breakpoints {
	return Breakpoints{
		Medium:   600,
		Expanded: 900,
	}
}

// CampusOSBreakpoints returns institutional responsive breakpoints matching Design System 2026:
// - Compact (< 750dp): Mobile view with off-canvas drawer and single-column cards
// - Medium (750dp - 1050dp): Tablet view with collapsed 64dp icon rail and 2-column cards
// - Expanded (>= 1050dp): Desktop view with full 240dp sidebar and 4-column cards
func CampusOSBreakpoints() Breakpoints {
	return Breakpoints{
		Medium:   750,
		Expanded: 1050,
	}
}

// HeightBreakpoints configures vertical thresholds separating Short, Normal, and Tall viewports.
type HeightBreakpoints struct {
	// Short is the vertical threshold separating Short from Normal (default: 450dp).
	Short float32
	// Tall is the vertical threshold separating Normal from Tall (default: 900dp).
	Tall float32
}

// DefaultHeightBreakpoints returns industry standard vertical responsive breakpoints.
func DefaultHeightBreakpoints() HeightBreakpoints {
	return HeightBreakpoints{
		Short: 450,
		Tall:  900,
	}
}
