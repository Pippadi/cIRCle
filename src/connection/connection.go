package connection

import (
	"log"
	"net"
	"strings"

	"github.com/Pippadi/cIRCle/src/buffer"
	"github.com/Pippadi/cIRCle/src/message"
	"gopkg.in/irc.v3"
)

type Connection struct {
	UI             *UI
	Nick           string
	IrcClient      *irc.Client
	ChannelBuffers map[string](*buffer.Buffer)
}

func New() *Connection {
	c := Connection{}
	c.UI = newUI()
	c.UI.ConnectBtn.OnTapped = c.connect
	c.ChannelBuffers = make(map[string](*buffer.Buffer))
	c.UI.SetConnectionState(false)
	return &c
}

func (conn *Connection) connect() {
	addr := conn.UI.AddrEntry.Text + ":" + conn.UI.PortEntry.Text
	sock, err := net.Dial("tcp", addr)
	if err != nil {
		log.Println(err)
	}
	clientConfig := irc.ClientConfig{
		Nick:    conn.UI.NickEntry.Text,
		Name:    conn.UI.NickEntry.Text,
		User:    conn.UI.NickEntry.Text,
		Pass:    conn.UI.PassEntry.Text,
		Handler: irc.HandlerFunc(conn.handler),
	}
	conn.Nick = conn.UI.NickEntry.Text
	conn.IrcClient = irc.NewClient(sock, clientConfig)

	conn.UI.JoinBtn.OnTapped = func() {
		conn.Join(conn.UI.JoinEntry.Text)
		conn.UI.JoinEntry.SetText("")
	}
	conn.UI.ConnectBtn.OnTapped = conn.disconnect

	go func() {
		conn.UI.SetConnectionState(true)
		err = conn.IrcClient.Run()
		if err != nil {
			log.Println(err)
		}
		conn.UI.SetConnectionState(false)
		conn.UI.ConnectBtn.OnTapped = conn.connect
	}()
}

func (conn *Connection) disconnect() {
	conn.IrcClient.Write("QUIT " + conn.Nick)
}

func (conn *Connection) handler(client *irc.Client, m *irc.Message) {
	switch strings.ToLower(m.Command) {
	case "privmsg":
		var buf *buffer.Buffer
		if conn.IrcClient.FromChannel(m) {
			buf = conn.ChannelBuffers[m.Params[0]] // m.Params[0] is the channel name. Messages can only come from joined channels.
		} else {
			var exists bool
			buf, exists = conn.ChannelBuffers[m.Prefix.Name] // m.Prefix.Name is the sender's name for PMs
			if !exists {
				conn.OpenPM(m.Prefix.Name)
				buf = conn.ChannelBuffers[m.Prefix.Name]
			}
		}
		buf.Incoming <- message.Message{m.Prefix.Name, m.Trailing()} // m.Prefix.Name is the sender's name
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
			conn.IrcClient.Write("PRIVMSG " + channel + " :" + (<-conn.ChannelBuffers[channel].Outgoing).Content)
		}
	}()
}

func (conn *Connection) AddBuffer(channel string) {
	buf := buffer.New(channel, conn.Nick)
	conn.ChannelBuffers[channel] = buf
	conn.UI.AddBuffer(buf)
}
