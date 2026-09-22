package layout

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

// Alignment specifies perpendicular or baseline alignment for flex layouts.
type Alignment int

const (
	AlignStart Alignment = iota
	AlignCenter
	AlignEnd
	AlignStretch
)

// flexItemWrapper wraps a CanvasObject with a flex expansion weight.
type flexItemWrapper struct {
	widget.BaseWidget
	child  fyne.CanvasObject
	weight float32
}

func (f *flexItemWrapper) CreateRenderer() fyne.WidgetRenderer {
	return widget.NewSimpleRenderer(f.child)
}

func (f *flexItemWrapper) MinSize() fyne.Size {
	return f.child.MinSize()
}

// FlexItem tags a CanvasObject with an expansion weight for Row or Column distribution.
func FlexItem(obj fyne.CanvasObject, weight float32) fyne.CanvasObject {
	if weight <= 0 {
		weight = 1.0
	}
	wrapper := &flexItemWrapper{child: obj, weight: weight}
	wrapper.ExtendBaseWidget(wrapper)
	return wrapper
}

// rowLayout arranges children horizontally with gap, alignment, and flex distribution.
type rowLayout struct {
	gap       float32
	alignment Alignment
}

// NewRowLayout creates a row layout with gap and alignment.
func NewRowLayout(gap float32, alignment Alignment) fyne.Layout {
	if gap < 0 {
		gap = 8
	}
	return &rowLayout{gap: gap, alignment: alignment}
}

func (r *rowLayout) Layout(objects []fyne.CanvasObject, size fyne.Size) {
	visible := make([]fyne.CanvasObject, 0, len(objects))
	for _, o := range objects {
		if o.Visible() {
			visible = append(visible, o)
		}
	}
	if len(visible) == 0 {
		return
	}

	totalGaps := float32(len(visible)-1) * r.gap
	fixedWidth := float32(0)
	totalWeight := float32(0)

	for _, o := range visible {
		if fw, ok := o.(*flexItemWrapper); ok {
			totalWeight += fw.weight
		} else {
			fixedWidth += o.MinSize().Width
		}
	}

	availForFlex := size.Width - totalGaps - fixedWidth
	if availForFlex < 0 {
		availForFlex = 0
	}

	curX := float32(0)
	for _, o := range visible {
		var itemW float32
		if fw, ok := o.(*flexItemWrapper); ok {
			if totalWeight > 0 {
				itemW = availForFlex * (fw.weight / totalWeight)
			} else {
				itemW = fw.child.MinSize().Width
			}
		} else {
			itemW = o.MinSize().Width
		}

		itemH := o.MinSize().Height
		curY := float32(0)

		switch r.alignment {
		case AlignCenter:
			if size.Height > itemH {
				curY = (size.Height - itemH) / 2
			}
		case AlignEnd:
			if size.Height > itemH {
				curY = size.Height - itemH
			}
		case AlignStretch:
			itemH = size.Height
			curY = 0
		default: // AlignStart
			curY = 0
		}

		o.Move(fyne.NewPos(curX, curY))
		o.Resize(fyne.NewSize(itemW, itemH))
		curX += itemW + r.gap
	}
}

func (r *rowLayout) MinSize(objects []fyne.CanvasObject) fyne.Size {
	totalW := float32(0)
	maxH := float32(0)
	visibleCount := 0

	for _, o := range objects {
		if !o.Visible() {
			continue
		}
		visibleCount++
		ms := o.MinSize()
		totalW += ms.Width
		if ms.Height > maxH {
			maxH = ms.Height
		}
	}

	if visibleCount > 1 {
		totalW += float32(visibleCount-1) * r.gap
	}
	return fyne.NewSize(totalW, maxH)
}

// Row creates an aligned horizontal container.
func Row(gap float32, objects ...fyne.CanvasObject) *fyne.Container {
	return container.New(NewRowLayout(gap, AlignCenter), objects...)
}

// RowAligned creates an aligned horizontal container with custom alignment.
func RowAligned(gap float32, alignment Alignment, objects ...fyne.CanvasObject) *fyne.Container {
	return container.New(NewRowLayout(gap, alignment), objects...)
}
