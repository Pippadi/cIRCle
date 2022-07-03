package connection

import (
	"errors"
	"strings"

	"fyne.io/fyne/v2"
	"github.com/Pippadi/cIRCle/src/buffer"
	"github.com/Pippadi/cIRCle/src/ircclient"
	"github.com/Pippadi/cIRCle/src/message"
	"gopkg.in/irc.v3"
)

type Connection struct {
	UI      *UI
	Nick    string
	client  *ircclient.IRCClient
	Buffers map[string](*buffer.Buffer)
}

func New(w fyne.Window) *Connection {
	c := Connection{}
	c.UI = newUI(w)
	c.UI.ConnectBtn.OnTapped = c.connect
	c.Buffers = make(map[string](*buffer.Buffer))
	c.UI.SetConnectionState(false)
	return &c
}

func (conn *Connection) connect() {
	if !conn.UI.ConnParamsValid() {
		conn.UI.ShowError(errors.New("Invalid connection parameters"))
		return
	}

	addr := conn.UI.AddrEntry.Text + ":" + conn.UI.PortEntry.Text
	conn.Nick = conn.UI.NickEntry.Text
	conn.client = ircclient.New(addr, conn.UI.PassEntry.Text, conn.Nick, conn.handler)

	conn.UI.JoinBtn.OnTapped = conn.chat
	conn.UI.ConnectBtn.OnTapped = conn.disconnect

	conn.UI.SetConnectionState(true)

	go func() {
		if err := conn.client.Run(); err != nil {
			conn.UI.ShowError(err)
		}
		conn.UI.SetConnectionState(false)
		conn.UI.ConnectBtn.OnTapped = conn.connect
		for c, _ := range conn.Buffers {
			conn.RemoveBuffer(c)
		}
	}()
}

func (conn *Connection) disconnect() {
	conn.client.Quit()
}

func (conn *Connection) handler(client *irc.Client, m *irc.Message) {
	switch strings.ToLower(m.Command) {
	case "001": // 001 is a welcome event after which channels can be joined
		conn.UI.SetJoinable(true)
	case "privmsg":
		var buf *buffer.Buffer
		if client.FromChannel(m) {
			buf = conn.Buffers[m.Params[0]] // m.Params[0] is the channel name. Messages can only come from joined channels.
		} else {
			var exists bool
			buf, exists = conn.Buffers[m.Prefix.Name] // m.Prefix.Name is the sender's name for PMs
			if !exists {
				conn.OpenPM(m.Prefix.Name)
				buf = conn.Buffers[m.Prefix.Name]
			}
		}
		buf.Incoming <- message.Message{From: m.Prefix.Name, To: conn.Nick, Content: m.Trailing()} // m.Prefix.Name is the sender's name
	case "join":
		buf := conn.Buffers[m.Params[0]]
		buf.CommandIn <- message.Command{m.Prefix.Name, "join"}
	case "part":
		buf := conn.Buffers[m.Params[0]]
		buf.CommandIn <- message.Command{m.Prefix.Name, "part"}
	}
}

func (conn *Connection) Join(channel string) {
	conn.client.Join(channel)
	conn.AddBuffer(channel)
	go conn.client.ListenAndWriteMessages(conn.Buffers[channel].Outgoing)
}

func (conn *Connection) OpenPM(who string) {
	conn.AddBuffer(who)
	go conn.client.ListenAndWriteMessages(conn.Buffers[who].Outgoing)
}

func (conn *Connection) AddBuffer(channel string) {
	buf := buffer.New(channel, conn.Nick)
	conn.Buffers[channel] = buf
	conn.UI.AddBuffer(buf)
}

func (conn *Connection) RemoveBuffer(channel string) {
	conn.UI.RemoveBuffer(conn.Buffers[channel])
	delete(conn.Buffers, channel)
}
