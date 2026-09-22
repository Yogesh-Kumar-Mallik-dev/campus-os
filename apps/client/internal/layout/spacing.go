package layout

// Spacing defines a centralized token hierarchy for consistent padding and gaps across the application.
type Spacing struct {
	XS float32 // Extra small (default: 4dp)
	SM float32 // Small (default: 8dp)
	MD float32 // Medium (default: 16dp)
	LG float32 // Large (default: 24dp)
	XL float32 // Extra large (default: 32dp)
}

// DefaultSpacing returns standard spacing tokens.
func DefaultSpacing() Spacing {
	return Spacing{
		XS: 4,
		SM: 8,
		MD: 16,
		LG: 24,
		XL: 32,
	}
}

// ForSizeClass returns adapted spacing tokens scaled appropriately for the active viewport SizeClass.
func (s Spacing) ForSizeClass(sc SizeClass) Spacing {
	switch sc {
	case Compact:
		return Spacing{
			XS: s.XS * 0.75,
			SM: s.SM * 0.75,
			MD: s.MD * 0.75,
			LG: s.LG * 0.75,
			XL: s.XL * 0.75,
		}
	case Medium:
		return s
	case Expanded:
		return Spacing{
			XS: s.XS,
			SM: s.SM,
			MD: s.MD * 1.25,
			LG: s.LG * 1.25,
			XL: s.XL * 1.25,
		}
	default:
		return s
	}
}
