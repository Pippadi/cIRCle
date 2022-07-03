package buffer

import "github.com/Pippadi/cIRCle/src/message"

type Buffer struct {
	UI        *UI
	channel   string
	Incoming  chan message.Message
	Outgoing  chan message.Message
	CommandIn chan message.Command
	nick      string
}

func New(channel, nick string) *Buffer {
	b := Buffer{}
	b.UI = newUI(channel)
	b.channel = channel
	b.Incoming, b.Outgoing = make(chan message.Message), make(chan message.Message)
	b.CommandIn = make(chan message.Command)
	b.nick = nick
	b.UI.SendBtn.OnTapped = b.sendMsg
	go func() {
		for {
			b.UI.AddMessageToChat(<-b.Incoming)
		}
	}()
	go func() {
		for {
			b.handleCommand(<-b.CommandIn)
		}
	}()
	return &b
}

func (b *Buffer) sendMsg() {
	if b.UI.MsgEntry.Text != "" {
		msg := message.Message{From: b.nick, To: b.channel, Content: b.UI.MsgEntry.Text}
		b.Outgoing <- msg
		b.UI.AddMessageToChat(msg)
		b.UI.MsgEntry.SetText("")
	}
}

func (b *Buffer) handleCommand(cmd message.Command) {
	b.UI.HandleCommand(cmd)
}
