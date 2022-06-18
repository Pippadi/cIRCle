package main

import (
	"net"
	"strings"

	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"github.com/Pippadi/cIRCle/src/buffer"
	"github.com/Pippadi/cIRCle/src/message"
	"gopkg.in/irc.v3"
)

func main() {
	a := app.New()
	w := a.NewWindow("cIRCle")
	incoming, outgoing := make(chan message.Message), make(chan message.Message)
	buffer := buffer.New(incoming, outgoing)

	msgHandler := irc.HandlerFunc(func(c *irc.Client, m *irc.Message) {
		switch strings.ToLower(m.Command) {
		case "001":
			c.Write("JOIN #op")

		case "privmsg":
			buffer.Incoming <- message.Message{m.Prefix.Name, m.Trailing()}
		}
	})

	conn, _ := net.Dial("tcp", "localhost:6667")
	ircConf := irc.ClientConfig{
		Nick:    "cIRCle",
		Name:    "cIRCle",
		User:    "cIRCle",
		Pass:    "password",
		Handler: msgHandler,
	}

	client := irc.NewClient(conn, ircConf)
	go client.Run()

	go func() {
		for {
			client.Write("PRIVMSG #op :" + (<-buffer.Outgoing).Content)
		}
	}()

	bufTab := container.NewTabItem("Buffer", buffer.UI.CanvasObject())

	w.SetContent(container.NewAppTabs(bufTab))
	w.ShowAndRun()
}
