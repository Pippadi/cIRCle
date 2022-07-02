package connection

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"github.com/Pippadi/cIRCle/src/buffer"
	"github.com/Pippadi/cIRCle/src/utils"
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
	window      fyne.Window
}

func newUI(w fyne.Window) *UI {
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
	ui.JoinEntry.SetPlaceHolder("Channel or Nick")
	ui.JoinBtn = widget.NewButton("Chat", func() {})

	connectPane := container.NewVBox(ui.inputFields, ui.ConnectBtn, layout.NewSpacer(), ui.JoinEntry, ui.JoinBtn)

	ui.tabStack = container.NewAppTabs(container.NewTabItem("Connect", connectPane))

	ui.window = w

	return &ui
}

func (ui *UI) CanvasObject() fyne.CanvasObject {
	return ui.tabStack
}

func (ui *UI) SetConnectionState(connected bool) {
	//var enableWhenConnected = []fyne.Disableable{}
	var disableWhenConnected = []fyne.Disableable{ui.AddrEntry, ui.PortEntry, ui.NickEntry, ui.PassEntry}
	//utils.SetWidgetsActive(connected, enableWhenConnected)
	utils.SetWidgetsActive(!connected, disableWhenConnected)
	if connected {
		ui.ConnectBtn.SetText("Disconnect")
	} else {
		ui.ConnectBtn.SetText("Connect")
		ui.SetJoinable(false)
	}
}

func (ui *UI) SetJoinable(joinable bool) {
	var enableWhenJoinable = []fyne.Disableable{ui.JoinEntry, ui.JoinBtn}
	utils.SetWidgetsActive(joinable, enableWhenJoinable)
}

func (ui *UI) AddBuffer(buf *buffer.Buffer) {
	ui.tabStack.Append(buf.UI.TabItem())
}

func (ui *UI) RemoveBuffer(buf *buffer.Buffer) {
	ui.tabStack.Remove(buf.UI.TabItem())
}

func (ui *UI) ShowError(err error) {
	dialog.ShowError(err, ui.window)
}

func (ui *UI) ConnParamsValid() bool {
	return ui.AddrEntry.Text != "" && ui.NickEntry.Text != "" && ui.PortEntry.Text != ""
}
