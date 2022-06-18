package main

import (
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/Pippadi/cIRCle/src/connection"
	"gopkg.in/irc.v3"
)

func main() {
	a := app.New()
	w := a.NewWindow("cIRCle")

	ircConf := irc.ClientConfig{
		Nick: "cIRCle",
		Name: "cIRCle",
		User: "cIRCle",
		Pass: "password",
	}
	conn := connection.New("localhost:6667", ircConf)

	tabs := container.NewAppTabs()

	joinEntry := widget.NewEntry()
	joinEntry.SetPlaceHolder("Channel")
	joinBtn := widget.NewButton("Join", func() {
		buf := conn.Join(joinEntry.Text)
		bufTab := container.NewTabItem(buf.Channel, buf.UI.CanvasObject())
		tabs.Append(bufTab)
		joinEntry.Text = ""
		joinEntry.Refresh()
	})
	joinVBox := container.NewVBox(joinEntry, joinBtn)
	tabs.Append(container.NewTabItem("Join", joinVBox))

	w.SetContent(tabs)
	w.ShowAndRun()
}
