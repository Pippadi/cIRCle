package connection

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

type UI struct {
	box        *fyne.Container
	AddrEntry  *widget.Entry
	PortEntry  *widget.Entry
	NickEntry  *widget.Entry
	PassEntry  *widget.Entry
	ConnectBtn *widget.Button
}

func newUI() *UI {
	ui := UI{}
	ui.AddrEntry = widget.NewEntry()
	ui.AddrEntry.SetPlaceHolder("Address")
	ui.PortEntry = widget.NewEntry()
	ui.PortEntry.SetPlaceHolder("Port")
	ui.PortEntry.SetText("6667")
	ui.NickEntry = widget.NewEntry()
	ui.NickEntry.SetPlaceHolder("Nick")
	ui.PassEntry = widget.NewEntry()
	ui.PassEntry.SetPlaceHolder("Password")
	ui.ConnectBtn = widget.NewButton("Connect", func() {})
	inputs := container.NewVBox(ui.AddrEntry, ui.PortEntry, ui.NickEntry, ui.PassEntry)
	ui.box = container.NewVBox(inputs, ui.ConnectBtn)
	return &ui
}

func (ui *UI) CanvasObject() fyne.CanvasObject {
	return ui.box
}
