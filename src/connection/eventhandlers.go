package connection

import (
	"errors"

	"github.com/Pippadi/cIRCle/src/buffer"
	"github.com/Pippadi/cIRCle/src/message"
)

func (conn *Connection) onJoinBtnTapped() {
	toJoin := conn.UI.JoinEntry.Text
	if toJoin == "" {
		conn.UI.ShowError(errors.New("Invalid channel/nick"))
		return
	}
	if toJoin[0] == '#' {
		conn.join(conn.UI.JoinEntry.Text)
	} else {
		conn.openPM(toJoin)
	}
	conn.UI.JoinEntry.SetText("")
}

func (conn *Connection) onJoinable() {
	conn.UI.SetJoinable(true)
	for _, channel := range conn.autojoin {
		conn.join(channel)
	}
}

func (conn *Connection) onIncomingMessage(msg message.Message) {
	var buf *buffer.Buffer
	if msg.To == conn.Nick { // Private message
		var exists bool
		buf, exists = conn.Buffers[msg.From]
		if !exists {
			conn.openPM(msg.From)
			buf = conn.Buffers[msg.From]
		}
	} else { // Message on channel
		buf = conn.Buffers[msg.To]
	}
	buf.Incoming <- msg
}

func (conn *Connection) onPersonJoined(person, channel string) {
	conn.Buffers[channel].CommandIn <- message.Command{Person: person, Action: "join"}
}

func (conn *Connection) onPersonParted(person, channel string) {
	if person != conn.Nick {
		conn.Buffers[channel].CommandIn <- message.Command{Person: person, Action: "part"}
	}
}

func (conn *Connection) onPersonQuit(person string) {
	for _, buf := range conn.Buffers {
		buf.CommandIn <- message.Command{Person: person, Action: "quit"}
	}
}

func (conn *Connection) onDisconnected() {
	conn.UI.SetConnectionState(false)
	conn.UI.ConnectBtn.OnTapped = conn.connect
	for c, _ := range conn.Buffers {
		conn.removeBuffer(c)
	}
}

func (conn *Connection) onNames(channel string, names []string) {
	conn.Buffers[channel].CommandIn <- message.Command{Action: "names", Args: names}
}
