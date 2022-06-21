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
	conn.UI.SetConnParamsActive(false)
	conn.UI.JoinBtn.OnTapped = func() {
		conn.Join(conn.UI.JoinEntry.Text)
		conn.UI.JoinEntry.SetText("")
	}
	go func() {
		err = conn.IrcClient.Run()
		if err != nil {
			log.Println(err)
		}
		conn.UI.SetConnParamsActive(true)
	}()
}

func (conn *Connection) handler(client *irc.Client, m *irc.Message) {
	channel := m.Params[0]
	switch strings.ToLower(m.Command) {
	case "privmsg":
		buf, exists := conn.ChannelBuffers[channel] // m.Params[0] is the channel
		if !exists {
			conn.AddBuffer(channel)
		}
		buf.Incoming <- message.Message{m.Prefix.Name, m.Trailing()}
	}
}

func (conn *Connection) Join(channel string) {
	conn.IrcClient.Write("JOIN " + channel)
	conn.AddBuffer(channel)
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
