package connection

import (
	"errors"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"github.com/Pippadi/cIRCle/src/buffer"
	"github.com/Pippadi/cIRCle/src/utils"
	"github.com/Pippadi/cIRCle/src/widgets"
)

type UI struct {
	inputFields *fyne.Container
	tabStack    *container.AppTabs
	AddrEntry   *widgets.HostEntry
	PortEntry   *widgets.PortEntry
	NickEntry   *widget.Entry
	PassEntry   *widget.Entry
	ConnectBtn  *widget.Button
	JoinEntry   *widgets.EnterCatchingEntry
	JoinBtn     *widget.Button
	window      fyne.Window
}

func newUI(w fyne.Window) *UI {
	ui := UI{}

	ui.AddrEntry = widgets.NewHostEntry()

	ui.PortEntry = widgets.NewPortEntry()

	ui.NickEntry = widget.NewEntry()
	ui.NickEntry.SetPlaceHolder("Nick")
	ui.NickEntry.Validator = validNickString

	ui.PassEntry = widget.NewEntry()
	ui.PassEntry.SetPlaceHolder("Password (optional)")
	ui.PassEntry.Password = true

	ui.ConnectBtn = widget.NewButton("Connect", func() {})
	ui.inputFields = container.NewVBox(ui.AddrEntry, ui.PortEntry, ui.NickEntry, ui.PassEntry)

	ui.JoinEntry = widgets.NewEnterCatchingEntry()
	ui.JoinEntry.SetPlaceHolder("Channel or Nick")
	ui.JoinBtn = widget.NewButton("Chat", func() {})

	joinctrls := widgets.NewEntryButtonContainer(ui.JoinEntry, ui.JoinBtn)

	connectPane := container.NewVBox(ui.inputFields, ui.ConnectBtn, layout.NewSpacer(), joinctrls)

	ui.tabStack = container.NewAppTabs(container.NewTabItem("Connect", connectPane))

	ui.window = w

	return &ui
}

func (ui *UI) CanvasObject() fyne.CanvasObject {
	return ui.tabStack
}

func (ui *UI) SetConnectionState(connected bool) {
	var disableWhenConnected = []fyne.Disableable{ui.AddrEntry, ui.PortEntry, ui.NickEntry, ui.PassEntry}
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
	ui.tabStack.Select(ui.tabStack.Items[0])
	ui.tabStack.Remove(buf.UI.TabItem())
}

func (ui *UI) ShowError(err error) {
	dialog.ShowError(err, ui.window)
}

func (ui *UI) ValidateConnParams() error {
	if err := ui.AddrEntry.Validate(); err != nil {
		return err
	}
	if err := ui.NickEntry.Validate(); err != nil {
		return err
	}
	if err := ui.PortEntry.Validate(); err != nil {
		return err
	}
	return nil
}

func validNickString(nick string) error {
	if nick == "" {
		return errors.New("Nick cannot be empty")
	}
	return nil
}
