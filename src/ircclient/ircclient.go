package ircclient

import (
	"net"

	"github.com/Pippadi/cIRCle/src/message"
	"gopkg.in/irc.v3"
)

// Just a simple wrapper
type IRCClient struct {
	Nick   string
	irc    *irc.Client
	addr   string
	config irc.ClientConfig
}

type HandlerFunc irc.HandlerFunc

func New(addr, password, nick string, handler HandlerFunc) *IRCClient {
	c := IRCClient{Nick: nick}
	c.addr = addr

	c.config = irc.ClientConfig{
		Nick:    c.Nick,
		Name:    c.Nick,
		User:    c.Nick,
		Pass:    password,
		Handler: irc.HandlerFunc(handler),
	}
	return &c
}

func (c *IRCClient) Run() error {
	sock, err := net.Dial("tcp", c.addr)
	if err != nil {
		return err
	}

	c.irc = irc.NewClient(sock, c.config)
	c.irc.Run()
	return nil
}

func (c *IRCClient) Join(channel string) {
	c.irc.Write("JOIN " + channel)
}

func (c *IRCClient) Quit() {
	c.irc.Write("QUIT cIRCle")
}

func (c *IRCClient) WriteMessage(msg message.Message) {
	c.irc.Write("PRIVMSG " + msg.To + " :" + msg.Content)
}

func (c *IRCClient) ListenAndWriteMessages(outgoing chan message.Message) {
	for {
		c.WriteMessage(<-outgoing)
	}
}
