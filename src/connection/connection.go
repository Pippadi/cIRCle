package connection

import (
	"errors"

	"fyne.io/fyne/v2"
	"github.com/Pippadi/cIRCle/src/buffer"
	"github.com/Pippadi/cIRCle/src/ircclient"
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

	conn.client = ircclient.New(addr, conn.UI.PassEntry.Text, conn.Nick)
	conn.client.OnJoinable = conn.onJoinable
	conn.client.OnMessage = conn.onIncomingMessage
	conn.client.OnPersonJoined = conn.onPersonJoined
	conn.client.OnPersonParted = conn.onPersonParted

	conn.UI.JoinBtn.OnTapped = conn.onJoinBtnTapped
	conn.UI.ConnectBtn.OnTapped = conn.disconnect

	conn.UI.SetConnectionState(true)

	go func() {
		if err := conn.client.Run(); err != nil {
			conn.UI.ShowError(err)
		}
		conn.onDisconnected()
	}()
}

func (conn *Connection) disconnect() {
	conn.client.Quit()
}

func (conn *Connection) join(channel string) {
	conn.client.Join(channel)
	conn.AddBuffer(channel)
}

func (conn *Connection) openPM(who string) {
	conn.AddBuffer(who)
}

func (conn *Connection) AddBuffer(channel string) {
	buf := buffer.New(channel, conn.Nick)
	conn.Buffers[channel] = buf
	conn.UI.AddBuffer(buf)
	go conn.client.ListenAndWriteMessages(conn.Buffers[channel].Outgoing)
}

func (conn *Connection) RemoveBuffer(channel string) {
	conn.UI.RemoveBuffer(conn.Buffers[channel])
	delete(conn.Buffers, channel)
}
