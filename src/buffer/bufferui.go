package buffer

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
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
	b.ChatArea.Wrapping = fyne.TextWrapBreak
	scrollBox := container.NewVScroll(container.NewVBox(b.ChatArea))
	b.MsgEntry = widget.NewEntry()
	b.MsgEntry.SetPlaceHolder("Message")
	b.SendBtn = widget.NewButton("Send", func() {})
	controls := container.NewVBox(b.MsgEntry, b.SendBtn)
	b.box = container.New(layout.NewBorderLayout(nil, controls, nil, nil), controls, scrollBox)
	return &b
}

func (b *UI) AddMessageToChat(msg message.Message) {
	b.ChatArea.Segments = append(b.ChatArea.Segments, &widget.TextSegment{Text: msg.Person + ": " + msg.Content})
	b.ChatArea.Refresh()
}

func (b *UI) CanvasObject() fyne.CanvasObject {
	return b.box
}
