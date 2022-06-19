package main

import (
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"github.com/Pippadi/cIRCle/src/connection"
)

func main() {
	a := app.New()
	w := a.NewWindow("cIRCle")
	conn := connection.New()

	tabs := container.NewAppTabs()

	joinEntry := widget.NewEntry()
	joinEntry.SetPlaceHolder("Channel")
	joinBtn := widget.NewButton("Join", func() {
		buf := conn.Join(joinEntry.Text)
		bufTab := container.NewTabItem(buf.Channel, buf.UI.CanvasObject())
		tabs.Append(bufTab)
		joinEntry.SetText("")
	})
	joinVBox := container.NewVBox(joinEntry, joinBtn)

	connectVBox := container.NewVBox(conn.UI.CanvasObject(), layout.NewSpacer(), joinVBox)
	tabs.Append(container.NewTabItem("Connect", connectVBox))

	w.SetContent(tabs)
	w.ShowAndRun()
}
