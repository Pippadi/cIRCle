package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"github.com/Pippadi/cIRCle/src/connection"
)

func main() {
	a := app.New()
	w := a.NewWindow("cIRCle")
	conn := connection.New()

	w.SetContent(conn.UI.CanvasObject())
	w.Resize(fyne.NewSize(400, 450))
	w.ShowAndRun()
}
