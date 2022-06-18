package buffer

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/Pippadi/cIRCle/src/message"
)

type UI struct {
	MsgEntry *widget.Entry
	ChatArea *widget.RichText
	SendBtn  *widget.Button
	box      *fyne.Container
}

func newUI() *UI {
	b := UI{}
	b.ChatArea = widget.NewRichText()
	b.MsgEntry = widget.NewEntry()
	b.MsgEntry.SetPlaceHolder("Message")
	b.SendBtn = widget.NewButton("Send", func() {})
	b.box = container.NewVBox(b.ChatArea, b.MsgEntry, b.SendBtn)
	return &b
}

func (b *UI) AddMessageToChat(msg message.Message) {
	b.ChatArea.Segments = append(b.ChatArea.Segments, &widget.TextSegment{Text: msg.Person + ": " + msg.Content})
	b.ChatArea.Refresh()
}

func (b *UI) CanvasObject() fyne.CanvasObject {
	return b.box
}
