package buffer

import (
	"github.com/Pippadi/cIRCle/src/message"
	"github.com/Pippadi/cIRCle/src/utils"
)

type Buffer struct {
	UI         *UI
	channel    string
	Incoming   chan message.Message
	Outgoing   chan message.Message
	CommandIn  chan message.Command
	CommandOut chan message.Command
	nicklist   []string
	nick       string
}

func New(channel, nick string) *Buffer {
	b := Buffer{}
	b.channel = channel
	b.Incoming, b.Outgoing = make(chan message.Message), make(chan message.Message)
	b.CommandIn, b.CommandOut = make(chan message.Command), make(chan message.Command)
	b.nicklist = make([]string, 0)
	//b.nicklist = append(b.nicklist, nick)
	b.nick = nick

	b.UI = newUI(channel, &b.nicklist)
	b.UI.SendBtn.OnTapped = b.sendMsg
	b.UI.MsgEntry.OnEnter = b.sendMsg
	b.UI.CloseBtn.OnTapped = func() {
		b.CommandOut <- message.Command{Action: "part", Person: nick}
	}

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
	switch cmd.Action {
	case "names":
		b.nicklist = cmd.Args
		b.UI.HandleCommand(cmd)
	case "join":
		b.nicklist = append(b.nicklist, cmd.Person)
		b.UI.HandleCommand(cmd)
	case "part", "quit":
		if utils.In(cmd.Person, b.nicklist) {
			b.nicklist = utils.RemoveAtIndex(utils.IndexOf(cmd.Person, b.nicklist), b.nicklist)
			b.UI.HandleCommand(cmd)
		}
	default:
		b.UI.HandleCommand(cmd)
	}
}
