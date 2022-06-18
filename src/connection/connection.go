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
	IrcClient      *irc.Client
	ChannelBuffers map[string](*buffer.Buffer)
}

func New(addr string, clientConfig irc.ClientConfig) *Connection {
	sock, err := net.Dial("tcp", addr)
	log.Println(err)

	c := Connection{}
	c.ChannelBuffers = make(map[string](*buffer.Buffer))
	clientConfig.Handler = irc.HandlerFunc(c.handler)
	c.IrcClient = irc.NewClient(sock, clientConfig)
	go func() {
		e := c.IrcClient.Run()
		log.Println(e)
	}()
	return &c
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
