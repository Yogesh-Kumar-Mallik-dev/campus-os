package ui

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

// ModalOptions controls dismissibility and alignment behavior for overlays and dialogs.
type ModalOptions struct {
	CloseOnEsc          bool // Dismiss on ESC key (default: true)
	CloseOnClickOutside bool // Dismiss when clicking the backdrop outside content (default: true)
	AlignLeft           bool // If true, docks content to the left edge full-height (drawer mode)
	OnDismiss           func()
}

// DefaultModalOptions returns the standard accessible overlay defaults.
func DefaultModalOptions() ModalOptions {
	return ModalOptions{
		CloseOnEsc:          true,
		CloseOnClickOutside: true,
		AlignLeft:           false,
	}
}

// ModalBackdrop provides a full-viewport clickable backdrop overlay.
type ModalBackdrop struct {
	widget.BaseWidget
	OnTapped func()
}

func NewModalBackdrop(onTapped func()) *ModalBackdrop {
	b := &ModalBackdrop{OnTapped: onTapped}
	b.ExtendBaseWidget(b)
	return b
}

func (b *ModalBackdrop) Tapped(_ *fyne.PointEvent) {
	if b.OnTapped != nil {
		b.OnTapped()
	}
}

func (b *ModalBackdrop) CreateRenderer() fyne.WidgetRenderer {
	bg := canvas.NewRectangle(color.NRGBA{R: 0, G: 0, B: 0, A: 175})
	return &modalBackdropRenderer{bg: bg}
}

type modalBackdropRenderer struct {
	bg *canvas.Rectangle
}

func (r *modalBackdropRenderer) Layout(size fyne.Size) {
	r.bg.Resize(size)
}

func (r *modalBackdropRenderer) MinSize() fyne.Size {
	return fyne.NewSize(100, 100)
}

func (r *modalBackdropRenderer) Refresh() {
	r.bg.Refresh()
}

func (r *modalBackdropRenderer) Objects() []fyne.CanvasObject {
	return []fyne.CanvasObject{r.bg}
}

func (r *modalBackdropRenderer) Destroy() {}

// ModalPlate wraps modal content and absorbs tap events so clicks inside the card never reach the backdrop.
type ModalPlate struct {
	widget.BaseWidget
	content fyne.CanvasObject
}

func NewModalPlate(content fyne.CanvasObject) *ModalPlate {
	p := &ModalPlate{content: content}
	p.ExtendBaseWidget(p)
	return p
}

func (p *ModalPlate) Tapped(_ *fyne.PointEvent) {
	// Absorb clicks inside the modal card/drawer
}

func (p *ModalPlate) CreateRenderer() fyne.WidgetRenderer {
	return &modalPlateRenderer{plate: p}
}

type modalPlateRenderer struct {
	plate *ModalPlate
}

func (r *modalPlateRenderer) Layout(size fyne.Size) {
	if r.plate.content != nil {
		r.plate.content.Resize(size)
	}
}

func (r *modalPlateRenderer) MinSize() fyne.Size {
	if r.plate.content != nil {
		return r.plate.content.MinSize()
	}
	return fyne.NewSize(0, 0)
}

func (r *modalPlateRenderer) Refresh() {
	if r.plate.content != nil {
		r.plate.content.Refresh()
	}
}

func (r *modalPlateRenderer) Objects() []fyne.CanvasObject {
	if r.plate.content != nil {
		return []fyne.CanvasObject{r.plate.content}
	}
	return nil
}

func (r *modalPlateRenderer) Destroy() {}

// ShowModal displays any CanvasObject within a responsive overlay with configurable
// ESC key and backdrop click dismissal.
func ShowModal(w fyne.Window, content fyne.CanvasObject, opts *ModalOptions) *widget.PopUp {
	if w == nil || w.Canvas() == nil {
		return nil
	}

	options := DefaultModalOptions()
	if opts != nil {
		options = *opts
	}

	var popup *widget.PopUp
	dismissed := false

	dismiss := func() {
		if dismissed {
			return
		}
		dismissed = true
		if popup != nil {
			popup.Hide()
		}
		if options.OnDismiss != nil {
			options.OnDismiss()
		}
	}

	// ESC Key listener
	if options.CloseOnEsc {
		prevKeyHandler := w.Canvas().OnTypedKey()
		w.Canvas().SetOnTypedKey(func(ev *fyne.KeyEvent) {
			if ev.Name == fyne.KeyEscape {
				w.Canvas().SetOnTypedKey(prevKeyHandler)
				dismiss()
				return
			}
			if prevKeyHandler != nil {
				prevKeyHandler(ev)
			}
		})

		originalDismiss := options.OnDismiss
		options.OnDismiss = func() {
			w.Canvas().SetOnTypedKey(prevKeyHandler)
			if originalDismiss != nil {
				originalDismiss()
			}
		}
	}

	// Clickable backdrop
	backdrop := NewModalBackdrop(func() {
		if options.CloseOnClickOutside {
			dismiss()
		}
	})

	plate := NewModalPlate(content)

	var foreground fyne.CanvasObject
	if options.AlignLeft {
		// Docked cleanly to the left edge for off-canvas drawer mode
		foreground = container.NewBorder(nil, nil, plate, nil, nil)
	} else {
		// Centered dialog
		foreground = container.NewCenter(plate)
	}

	overlay := container.NewStack(
		backdrop,
		foreground,
	)

	popup = widget.NewPopUp(overlay, w.Canvas())

	canvasSize := w.Canvas().Size()
	if canvasSize.Width <= 0 {
		canvasSize.Width = 800
	}
	if canvasSize.Height <= 0 {
		canvasSize.Height = 600
	}
	popup.Resize(canvasSize)
	popup.Move(fyne.NewPos(0, 0))
	popup.Show()

	return popup
}

// ShowAlertDialogWithOptions renders a modal confirmation dialog overlay with configurable accessibility options.
func ShowAlertDialogWithOptions(w fyne.Window, title, description, confirmText string, isDestructive bool, onConfirm func(), opts ModalOptions) *widget.PopUp {
	var popup *widget.PopUp

	cancelBtn := NewShadcnButton("Cancel", ButtonOutline, ButtonSizeDefault, nil, func() {
		if popup != nil {
			popup.Hide()
		}
	})

	btnVariant := ButtonDefault
	if isDestructive {
		btnVariant = ButtonDestructive
	}

	confirmBtn := NewShadcnButton(confirmText, btnVariant, ButtonSizeDefault, nil, func() {
		if popup != nil {
			popup.Hide()
		}
		if onConfirm != nil {
			onConfirm()
		}
	})

	buttons := container.NewHBox(
		layout.NewSpacer(),
		cancelBtn,
		confirmBtn,
	)

	card := NewShadcnCard(CardParts{
		Title:       title,
		Description: description,
		Footer:      buttons,
	})

	cardMin := card.MinSize()
	targetWidth := float32(440)
	if w != nil && w.Canvas() != nil && w.Canvas().Size().Width > 0 {
		avail := w.Canvas().Size().Width - 32
		if avail < targetWidth {
			targetWidth = avail
		}
	}
	targetHeight := cardMin.Height
	if targetHeight < 180 {
		targetHeight = 180
	}

	wrapped := container.NewGridWrap(fyne.NewSize(targetWidth, targetHeight), card)
	popup = ShowModal(w, wrapped, &opts)
	return popup
}

// ShowAlertDialog renders a standard modal confirmation dialog overlay matching shadcn alert-dialog.
// Defaults to CloseOnEsc=true and CloseOnClickOutside=true.
func ShowAlertDialog(w fyne.Window, title, description, confirmText string, isDestructive bool, onConfirm func()) *widget.PopUp {
	return ShowAlertDialogWithOptions(w, title, description, confirmText, isDestructive, onConfirm, DefaultModalOptions())
}

