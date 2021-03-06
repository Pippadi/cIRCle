package buffer

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"github.com/Pippadi/cIRCle/src/message"
	"github.com/Pippadi/cIRCle/src/utils"
	"github.com/Pippadi/cIRCle/src/widgets"
)

type UI struct {
	MsgEntry       *widgets.EnterCatchingEntry
	ChatArea       *widget.RichText
	SendBtn        *widget.Button
	CloseBtn       *widget.Button
	tabItem        *container.TabItem
	chatAreaScroll *container.Scroll
	nickListWidget *widget.List
	nickList       binding.ExternalStringList
}

func newUI(channel string, nicklist *[]string) *UI {
	b := UI{}

	b.ChatArea = widget.NewRichText()
	b.ChatArea.Wrapping = fyne.TextWrapBreak
	b.chatAreaScroll = container.NewVScroll(container.NewVBox(b.ChatArea))

	b.MsgEntry = widgets.NewEnterCatchingEntry()
	b.MsgEntry.SetPlaceHolder("Message")

	b.SendBtn = widget.NewButton("Send", func() {})
	b.CloseBtn = widget.NewButton("Close", func() {})

	b.nickList = binding.BindStringList(nicklist)
	b.nickListWidget = widget.NewListWithData(
		b.nickList,
		func() fyne.CanvasObject { return widget.NewLabel("template") },
		func(i binding.DataItem, o fyne.CanvasObject) {
			o.(*widget.Label).Bind(i.(binding.String))
		})
	b.nickList.Reload()

	controls := container.New(layout.NewBorderLayout(nil, nil, b.CloseBtn, b.SendBtn), b.CloseBtn, b.SendBtn, b.MsgEntry)

	if channel[0] == '#' {
		splitView := container.NewVSplit(b.chatAreaScroll, b.nickListWidget)
		splitView.SetOffset(0.8) // Make the nick list smaller than the chat area
		b.tabItem = container.NewTabItem(channel, container.New(layout.NewBorderLayout(nil, controls, nil, nil), controls, splitView))
	} else {
		b.tabItem = container.NewTabItem(channel, container.New(layout.NewBorderLayout(nil, controls, nil, nil), controls, b.chatAreaScroll))
	}

	return &b
}

func (b *UI) addTextToChatArea(text string) {
	b.ChatArea.Segments = append(b.ChatArea.Segments, &widget.TextSegment{Text: text})
	b.ChatArea.Refresh()
	b.chatAreaScroll.ScrollToBottom()
}

func (b *UI) AddMessageToChat(msg message.Message) {
	b.addTextToChatArea(msg.From + ": " + msg.Content)
}

func (b *UI) TabItem() *container.TabItem {
	return b.tabItem
}

func (b *UI) SetActive(active bool) {
	var toset = []fyne.Disableable{b.MsgEntry, b.SendBtn}
	utils.SetWidgetsActive(active, toset)
}

func (b *UI) HandleCommand(cmd message.Command) {
	switch cmd.Action {
	case "quit":
		b.addTextToChatArea(cmd.Person + " has quit")
	case "join":
		b.addTextToChatArea(cmd.Person + " has joined")
	case "part":
		b.addTextToChatArea(cmd.Person + " has left")
	}
	b.nickList.Reload()
}
