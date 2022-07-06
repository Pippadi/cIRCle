package buffer

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"github.com/Pippadi/cIRCle/src/message"
	"github.com/Pippadi/cIRCle/src/utils"
)

type UI struct {
	MsgEntry       *widget.Entry
	ChatArea       *widget.RichText
	SendBtn        *widget.Button
	tabItem        *container.TabItem
	chatAreaScroll *container.Scroll
}

func newUI(channel string) *UI {
	b := UI{}
	b.ChatArea = widget.NewRichText()
	b.ChatArea.Wrapping = fyne.TextWrapBreak
	b.chatAreaScroll = container.NewVScroll(container.NewVBox(b.ChatArea))
	b.MsgEntry = widget.NewEntry()
	b.MsgEntry.SetPlaceHolder("Message")
	b.SendBtn = widget.NewButton("Send", func() {})
	controls := utils.NewEntryButtonContainer(b.MsgEntry, b.SendBtn)
	b.tabItem = container.NewTabItem(channel, container.New(layout.NewBorderLayout(nil, controls, nil, nil), controls, b.chatAreaScroll))
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
}
