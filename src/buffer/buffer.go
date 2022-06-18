package buffer

import "github.com/Pippadi/cIRCle/src/message"

type Buffer struct {
	UI       *UI
	Incoming chan message.Message
	Outgoing chan message.Message
}

func New(incoming, outgoing chan message.Message) *Buffer {
	b := Buffer{}
	b.UI = newUI()
	b.Incoming = incoming
	b.Outgoing = outgoing
	b.UI.SendBtn.OnTapped = func() {
		msg := message.Message{"cIRCle", b.UI.MsgEntry.Text}
		b.Outgoing <- msg
		b.UI.AddMessageToChat(msg)
		b.UI.MsgEntry.Text = ""
		b.UI.MsgEntry.Refresh()
	}
	go func() {
		for {
			b.UI.AddMessageToChat(<-incoming)
		}
	}()
	return &b
}
