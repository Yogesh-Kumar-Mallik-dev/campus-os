package layout

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
)

// stackLayout overlays all children in the same bounding area.
type stackLayout struct{}

// NewStackLayout creates a new layout that stacks children on top of each other.
func NewStackLayout() fyne.Layout {
	return &stackLayout{}
}

func (s *stackLayout) Layout(objects []fyne.CanvasObject, size fyne.Size) {
	for _, o := range objects {
		if o.Visible() {
			o.Move(fyne.NewPos(0, 0))
			o.Resize(size)
		}
	}
}

func (s *stackLayout) MinSize(objects []fyne.CanvasObject) fyne.Size {
	maxW := float32(0)
	maxH := float32(0)
	for _, o := range objects {
		if o.Visible() {
			ms := o.MinSize()
			if ms.Width > maxW {
				maxW = ms.Width
			}
			if ms.Height > maxH {
				maxH = ms.Height
			}
		}
	}
	return fyne.NewSize(maxW, maxH)
}

// Stack creates an overlay container layering objects on top of each other.
func Stack(objects ...fyne.CanvasObject) *fyne.Container {
	return container.New(NewStackLayout(), objects...)
}
