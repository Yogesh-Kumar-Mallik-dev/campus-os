package layout

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
)

// columnLayout arranges children vertically with gap, alignment, and flex distribution.
type columnLayout struct {
	gap       float32
	alignment Alignment
}

// NewColumnLayout creates a column layout with gap and alignment.
func NewColumnLayout(gap float32, alignment Alignment) fyne.Layout {
	if gap < 0 {
		gap = 8
	}
	return &columnLayout{gap: gap, alignment: alignment}
}

func (c *columnLayout) Layout(objects []fyne.CanvasObject, size fyne.Size) {
	visible := make([]fyne.CanvasObject, 0, len(objects))
	for _, o := range objects {
		if o.Visible() {
			visible = append(visible, o)
		}
	}
	if len(visible) == 0 {
		return
	}

	totalGaps := float32(len(visible)-1) * c.gap
	fixedHeight := float32(0)
	totalWeight := float32(0)

	for _, o := range visible {
		if fw, ok := o.(*flexItemWrapper); ok {
			totalWeight += fw.weight
		} else {
			fixedHeight += o.MinSize().Height
		}
	}

	availForFlex := size.Height - totalGaps - fixedHeight
	if availForFlex < 0 {
		availForFlex = 0
	}

	curY := float32(0)
	for _, o := range visible {
		var itemH float32
		if fw, ok := o.(*flexItemWrapper); ok {
			if totalWeight > 0 {
				itemH = availForFlex * (fw.weight / totalWeight)
			} else {
				itemH = fw.child.MinSize().Height
			}
		} else {
			itemH = o.MinSize().Height
		}

		itemW := o.MinSize().Width
		curX := float32(0)

		switch c.alignment {
		case AlignCenter:
			if size.Width > itemW {
				curX = (size.Width - itemW) / 2
			}
		case AlignEnd:
			if size.Width > itemW {
				curX = size.Width - itemW
			}
		case AlignStretch:
			itemW = size.Width
			curX = 0
		default: // AlignStart
			curX = 0
		}

		o.Move(fyne.NewPos(curX, curY))
		o.Resize(fyne.NewSize(itemW, itemH))
		curY += itemH + c.gap
	}
}

func (c *columnLayout) MinSize(objects []fyne.CanvasObject) fyne.Size {
	totalH := float32(0)
	maxW := float32(0)
	visibleCount := 0

	for _, o := range objects {
		if !o.Visible() {
			continue
		}
		visibleCount++
		ms := o.MinSize()
		totalH += ms.Height
		if ms.Width > maxW {
			maxW = ms.Width
		}
	}

	if visibleCount > 1 {
		totalH += float32(visibleCount-1) * c.gap
	}
	return fyne.NewSize(maxW, totalH)
}

// Column creates a vertical container with start alignment.
func Column(gap float32, objects ...fyne.CanvasObject) *fyne.Container {
	return container.New(NewColumnLayout(gap, AlignStart), objects...)
}

// ColumnAligned creates a vertical container with custom alignment.
func ColumnAligned(gap float32, alignment Alignment, objects ...fyne.CanvasObject) *fyne.Container {
	return container.New(NewColumnLayout(gap, alignment), objects...)
}
