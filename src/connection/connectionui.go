package connection

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"github.com/Pippadi/cIRCle/src/buffer"
)

type UI struct {
	inputFields *fyne.Container
	tabStack    *container.AppTabs
	AddrEntry   *widget.Entry
	PortEntry   *widget.Entry
	NickEntry   *widget.Entry
	PassEntry   *widget.Entry
	ConnectBtn  *widget.Button
	JoinEntry   *widget.Entry
	JoinBtn     *widget.Button
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
	ui.inputFields = container.NewVBox(ui.AddrEntry, ui.PortEntry, ui.NickEntry, ui.PassEntry)

	ui.JoinEntry = widget.NewEntry()
	ui.JoinEntry.SetPlaceHolder("Channel")
	ui.JoinBtn = widget.NewButton("Join", func() {})

	connectPane := container.NewVBox(ui.inputFields, ui.ConnectBtn, layout.NewSpacer(), ui.JoinEntry, ui.JoinBtn)

	ui.tabStack = container.NewAppTabs(container.NewTabItem("Connect", connectPane))

	return &ui
}

func (ui *UI) CanvasObject() fyne.CanvasObject {
	return ui.tabStack
}

func (ui *UI) SetConnectionState(connected bool) {
	var enableWhenConnected = []fyne.Disableable{ui.JoinEntry, ui.JoinBtn}
	var disableWhenConnected = []fyne.Disableable{ui.AddrEntry, ui.PortEntry, ui.NickEntry, ui.PassEntry}
	ui.setWidgetsActive(connected, enableWhenConnected)
	ui.setWidgetsActive(!connected, disableWhenConnected)
}

func (ui *UI) AddBuffer(buf *buffer.Buffer) {
	ui.tabStack.Append(container.NewTabItem(buf.Channel, buf.UI.CanvasObject()))
}

func (ui *UI) setWidgetsActive(active bool, widgets []fyne.Disableable) {
	for _, w := range widgets {
		if active {
			w.Enable()
		} else {
			w.Disable()
		}

	}
}
