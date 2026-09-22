package layout

import "fyne.io/fyne/v2"

// Clamp bounds a float32 value within [min, max].
func Clamp(val, min, max float32) float32 {
	if val < min {
		return min
	}
	if max >= min && val > max {
		return max
	}
	return val
}

// FluidWidth calculates usable content width within bounds [minW, maxW] after deducting margins.
func FluidWidth(availW, minW, maxW, margin float32) float32 {
	target := availW - (margin * 2)
	if target < 0 {
		target = 0
	}
	if minW > 0 && target < minW {
		target = minW
	}
	if maxW > 0 && target > maxW {
		target = maxW
	}
	return target
}

// Columns calculates the optimal number of grid columns that can fit in available width
// while guaranteeing each column is at least minItemW and spaced by gap. Always returns >= 1.
func Columns(availW, minItemW, gap float32) int {
	if minItemW <= 0 {
		return 1
	}
	if availW <= 0 {
		return 1
	}
	cols := int((availW + gap) / (minItemW + gap))
	if cols < 1 {
		return 1
	}
	return cols
}

// Lerp performs linear interpolation between a and b for t in [0, 1].
func Lerp(a, b, t float32) float32 {
	return a + (b-a)*Clamp(t, 0, 1)
}

// AspectScale scales srcSize to fit within target bounding box while maintaining aspect ratio.
func AspectScale(srcSize fyne.Size, maxW, maxH float32) fyne.Size {
	if srcSize.Width <= 0 || srcSize.Height <= 0 || maxW <= 0 || maxH <= 0 {
		return fyne.NewSize(0, 0)
	}
	scaleW := maxW / srcSize.Width
	scaleH := maxH / srcSize.Height
	scale := scaleW
	if scaleH < scale {
		scale = scaleH
	}
	return fyne.NewSize(srcSize.Width*scale, srcSize.Height*scale)
}
