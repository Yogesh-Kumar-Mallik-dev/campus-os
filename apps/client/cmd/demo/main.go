package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"github.com/Yogesh-Kumar-Mallik-dev/campus-os/apps/client/internal/layout"
)

func main() {
	a := app.NewWithID("com.campusos.layout.demo")
	w := a.NewWindow("Campus OS — Fully Adaptive Responsive Layout Engine")

	w.SetContent(layout.BuildDemoScreen(w))
	w.Resize(fyne.NewSize(1024, 720))
	w.CenterOnScreen()

	w.ShowAndRun()
}
