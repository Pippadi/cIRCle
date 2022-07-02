package connection

import (
	"errors"
	"net"
	"strings"

	"fyne.io/fyne/v2"
	"github.com/Pippadi/cIRCle/src/buffer"
	"github.com/Pippadi/cIRCle/src/message"
	"gopkg.in/irc.v3"
)

type Connection struct {
	UI        *UI
	Nick      string
	IrcClient *irc.Client
	Buffers   map[string](*buffer.Buffer)
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
	sock, err := net.Dial("tcp", addr)
	if err != nil {
		conn.UI.ShowError(err)
		return
	}

	conn.Nick = conn.UI.NickEntry.Text
	clientConfig := irc.ClientConfig{
		Nick:    conn.Nick,
		Name:    conn.Nick,
		User:    conn.Nick,
		Pass:    conn.UI.PassEntry.Text,
		Handler: irc.HandlerFunc(conn.handler),
	}
	conn.IrcClient = irc.NewClient(sock, clientConfig)

	conn.UI.JoinBtn.OnTapped = func() {
		conn.Join(conn.UI.JoinEntry.Text)
		conn.UI.JoinEntry.SetText("")
	}
	conn.UI.ConnectBtn.OnTapped = conn.disconnect

	go func() {
		conn.UI.SetConnectionState(true)
		err = conn.IrcClient.Run()
		conn.UI.SetConnectionState(false)
		conn.UI.ConnectBtn.OnTapped = conn.connect
		for c, _ := range conn.Buffers {
			conn.RemoveBuffer(c)
		}
	}()
}

func (conn *Connection) disconnect() {
	conn.IrcClient.Write("QUIT " + conn.Nick)
}

func (conn *Connection) handler(client *irc.Client, m *irc.Message) {
	switch strings.ToLower(m.Command) {
	case "001": // 001 is a welcome event after which channels can be joined
		conn.UI.SetJoinable(true)
	case "privmsg":
		var buf *buffer.Buffer
		if conn.IrcClient.FromChannel(m) {
			buf = conn.Buffers[m.Params[0]] // m.Params[0] is the channel name. Messages can only come from joined channels.
		} else {
			var exists bool
			buf, exists = conn.Buffers[m.Prefix.Name] // m.Prefix.Name is the sender's name for PMs
			if !exists {
				conn.OpenPM(m.Prefix.Name)
				buf = conn.Buffers[m.Prefix.Name]
			}
		}
		buf.Incoming <- message.Message{m.Prefix.Name, m.Trailing()} // m.Prefix.Name is the sender's name
	case "join":
		buf := conn.Buffers[m.Params[0]]
		buf.CommandIn <- message.Command{m.Prefix.Name, "join"}
	case "part":
		buf := conn.Buffers[m.Params[0]]
		buf.CommandIn <- message.Command{m.Prefix.Name, "part"}
	}
}

func (conn *Connection) Join(channel string) {
	conn.IrcClient.Write("JOIN " + channel)
	conn.AddBuffer(channel)
	conn.ListenAndWriteMessages(channel)
}

func (conn *Connection) OpenPM(who string) {
	conn.AddBuffer(who)
	conn.ListenAndWriteMessages(who)
}

func (conn *Connection) ListenAndWriteMessages(channel string) {
	go func() {
		for {
			conn.IrcClient.Write("PRIVMSG " + channel + " :" + (<-conn.Buffers[channel].Outgoing).Content)
		}
	}()
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
