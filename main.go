package main

import (
	"net"
	"strings"

	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"gopkg.in/irc.v3"
)

func main() {
	a := app.New()
	w := a.NewWindow("cIRCle")

	msgEntry := widget.NewEntry()
	msgEntry.SetPlaceHolder("Message")
	sendBtn := widget.NewButton("Send", func() {})
	chatBox := widget.NewRichText()
	vBox := container.NewVBox(chatBox, msgEntry, sendBtn)

	msgHandler := irc.HandlerFunc(func(c *irc.Client, m *irc.Message) {
		switch strings.ToLower(m.Command) {
		case "001":
			c.Write("JOIN #op")

		case "privmsg":
			chatBox.Segments = append(chatBox.Segments, &widget.TextSegment{Text: m.Trailing()})
			chatBox.Refresh()
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

	sendBtn.OnTapped = func() { client.Write("PRIVMSG #op :" + msgEntry.Text) }

	w.SetContent(vBox)
	w.ShowAndRun()
}
