package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"github.com/Pippadi/cIRCle/src/connection"
	"github.com/Pippadi/cIRCle/src/persistence"
)

func main() {
	a := app.NewWithID("com.plootarg.circle")
	w := a.NewWindow("cIRCle")
	p := persistence.New(a)
	conn := connection.New(w)
	conn.LoadConfig(p.LoadConnConfig())

	w.SetContent(conn.UI.CanvasObject())
	w.Resize(fyne.NewSize(400, 450))
	w.ShowAndRun()
	p.DumpConnConfig(conn.GetConfig())
}
