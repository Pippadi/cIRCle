package ircclient

import (
	"net"
	"strings"

	"github.com/Pippadi/cIRCle/src/message"
	"gopkg.in/irc.v3"
)

// Just a simple wrapper
type IRCClient struct {
	Nick           string
	irc            *irc.Client
	addr           string
	config         irc.ClientConfig
	OnJoinable     JoinableHandler
	OnMessage      MessageHandler
	OnPersonJoined PersonJoinedHandler
	OnPersonParted PersonPartedHandler
	OnNames        NamesHandler
	OnPersonQuit   QuitHandler
}

type MessageHandler func(msg message.Message)
type JoinableHandler func()
type PersonJoinedHandler func(person, channel string)
type PersonPartedHandler func(person, channel string)
type NamesHandler func(channel string, names []string)
type QuitHandler func(person string)

func New(addr, password, nick string) *IRCClient {
	c := IRCClient{Nick: nick}
	c.addr = addr

	c.config = irc.ClientConfig{
		Nick:    c.Nick,
		Name:    c.Nick,
		User:    c.Nick,
		Pass:    password,
		Handler: irc.HandlerFunc(c.handler),
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

func (c *IRCClient) handler(cl *irc.Client, m *irc.Message) {
	switch strings.ToLower(m.Command) {
	case "001": // 001 is a welcome event after which channels can be joined
		c.OnJoinable()
	case "join":
		c.OnPersonJoined(m.Prefix.Name, m.Params[0])
	case "part":
		c.OnPersonParted(m.Prefix.Name, m.Params[0])
	case "quit":
		c.OnPersonQuit(m.Prefix.Name)
	case "privmsg":
		var msg message.Message
		msg.From = m.Prefix.Name // Name of sender
		msg.Content = m.Trailing()
		if c.irc.FromChannel(m) {
			msg.To = m.Params[0] // m.Params[0] is the channel name
		} else {
			msg.To = c.Nick
		}
		c.OnMessage(msg)
	case "353": // NAMES (people online)
		c.OnNames(m.Params[len(m.Params)-2], strings.Split(m.Trailing(), " ")) // m.Params[len(m.Params)-2] is the channel
	}
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
