package buffer

import "github.com/Pippadi/cIRCle/src/message"

type Buffer struct {
	UI       *UI
	Channel  string
	Incoming chan message.Message
	Outgoing chan message.Message
}

func New(channel, nick string) *Buffer {
	b := Buffer{}
	b.UI = newUI()
	b.Channel = channel
	b.Incoming, b.Outgoing = make(chan message.Message), make(chan message.Message)
	b.UI.SendBtn.OnTapped = func() {
		if b.UI.MsgEntry.Text != "" {
			msg := message.Message{nick, b.UI.MsgEntry.Text}
			b.Outgoing <- msg
			b.UI.AddMessageToChat(msg)
			b.UI.MsgEntry.SetText("")
		}
	}
	go func() {
		for {
			b.UI.AddMessageToChat(<-b.Incoming)
		}
	}()
	return &b
}
