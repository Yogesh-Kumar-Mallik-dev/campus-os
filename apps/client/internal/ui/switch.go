package ui

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/widget"
)

// Switch implements a modern toggle switch component mirroring shadcn switch.svelte.
type Switch struct {
	widget.BaseWidget
	Checked   bool
	OnChanged func(bool)
}

// NewSwitch constructs a new toggle switch.
func NewSwitch(checked bool, onChanged func(bool)) *Switch {
	s := &Switch{
		Checked:   checked,
		OnChanged: onChanged,
	}
	s.ExtendBaseWidget(s)
	return s
}

// SetChecked programmatically updates the switch state.
func (s *Switch) SetChecked(checked bool) {
	s.Checked = checked
	s.Refresh()
	if s.OnChanged != nil {
		s.OnChanged(s.Checked)
	}
}

// Tapped toggles the checked state on click or touch tap.
func (s *Switch) Tapped(_ *fyne.PointEvent) {
	s.SetChecked(!s.Checked)
}

func (s *Switch) CreateRenderer() fyne.WidgetRenderer {
	track := canvas.NewRectangle(color.NRGBA{R: 40, G: 46, B: 62, A: 255})
	track.CornerRadius = 11

	knob := canvas.NewCircle(color.NRGBA{R: 255, G: 255, B: 255, A: 255})
	knob.Resize(fyne.NewSize(18, 18))

	r := &switchRenderer{
		sw:      s,
		track:   track,
		knob:    knob,
		objects: []fyne.CanvasObject{track, knob},
	}
	r.Refresh()
	return r
}

type switchRenderer struct {
	sw      *Switch
	track   *canvas.Rectangle
	knob    *canvas.Circle
	objects []fyne.CanvasObject
}

func (r *switchRenderer) Layout(size fyne.Size) {
	r.track.Resize(fyne.NewSize(42, 22))
	r.track.Move(fyne.NewPos(0, 0))

	if r.sw.Checked {
		r.knob.Move(fyne.NewPos(22, 2))
	} else {
		r.knob.Move(fyne.NewPos(2, 2))
	}
}

func (r *switchRenderer) MinSize() fyne.Size {
	return fyne.NewSize(44, 24)
}

func (r *switchRenderer) Refresh() {
	if r.sw.Checked {
		// Active state: primary institutional terracotta
		r.track.FillColor = color.NRGBA{R: 244, G: 90, B: 81, A: 255}
		r.knob.Move(fyne.NewPos(22, 2))
	} else {
		// Inactive state: secondary slate
		r.track.FillColor = color.NRGBA{R: 46, G: 56, B: 68, A: 255}
		r.knob.Move(fyne.NewPos(2, 2))
	}
	r.track.Refresh()
	r.knob.Refresh()
}

func (r *switchRenderer) Objects() []fyne.CanvasObject {
	return r.objects
}

func (r *switchRenderer) Destroy() {}
