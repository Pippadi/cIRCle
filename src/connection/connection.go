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
	IrcClient      *irc.Client
	ChannelBuffers map[string](*buffer.Buffer)
}

func New() *Connection {

	c := Connection{}
	c.UI = newUI()
	c.UI.ConnectBtn.OnTapped = c.connect
	c.ChannelBuffers = make(map[string](*buffer.Buffer))
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
	conn.IrcClient = irc.NewClient(sock, clientConfig)
	go func() {
		err = conn.IrcClient.Run()
		if err != nil {
			log.Println(err)
		}
	}()
}

func (conn *Connection) handler(client *irc.Client, m *irc.Message) {
	switch strings.ToLower(m.Command) {
	case "privmsg":
		conn.ChannelBuffers[m.Params[0]].Incoming <- message.Message{m.Prefix.Name, m.Trailing()}
	}
}

func (conn *Connection) Join(channel string) *buffer.Buffer {
	conn.IrcClient.Write("JOIN " + channel)
	buf := buffer.New(channel)
	conn.ChannelBuffers[channel] = buf
	go func() {
		for {
			conn.IrcClient.Write("PRIVMSG " + channel + " :" + (<-conn.ChannelBuffers[channel].Outgoing).Content)
		}
	}()
	return buf
}
