package connection

import (
	"fyne.io/fyne/v2"
	"github.com/Pippadi/cIRCle/src/buffer"
	"github.com/Pippadi/cIRCle/src/ircclient"
	"github.com/Pippadi/cIRCle/src/message"
	"github.com/Pippadi/cIRCle/src/utils"
)

type Connection struct {
	UI       *UI
	Nick     string
	client   *ircclient.IRCClient
	Buffers  map[string](*buffer.Buffer)
	autojoin []string
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
	if err := conn.UI.ValidateConnParams(); err != nil {
		conn.UI.ShowError(err)
		return
	}

	addr := conn.UI.AddrEntry.Text + ":" + conn.UI.PortEntry.Text
	conn.Nick = conn.UI.NickEntry.Text

	conn.client = ircclient.New(addr, conn.UI.PassEntry.Text, conn.Nick)
	conn.client.OnJoinable = conn.onJoinable
	conn.client.OnMessage = conn.onIncomingMessage
	conn.client.OnPersonJoined = conn.onPersonJoined
	conn.client.OnPersonParted = conn.onPersonParted
	conn.client.OnPersonQuit = conn.onPersonQuit
	conn.client.OnNames = conn.onNames

	conn.UI.JoinBtn.OnTapped = conn.onJoinBtnTapped
	conn.UI.JoinEntry.OnEnter = conn.onJoinBtnTapped
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
	conn.addBufferIfNotExists(channel)
	if !utils.In(channel, conn.autojoin) {
		conn.autojoin = append(conn.autojoin, channel)
	}
	go conn.handleCommandFromBuffer(channel)
}

func (conn *Connection) openPM(who string) {
	conn.addBufferIfNotExists(who)
	go conn.handleCommandFromBuffer(who)
}

func (conn *Connection) addBufferIfNotExists(channel string) {
	buf, exists := conn.Buffers[channel]
	if !exists {
		buf = buffer.New(channel, conn.Nick)
		conn.Buffers[channel] = buf
		conn.UI.AddBuffer(buf)
		go conn.client.ListenAndWriteMessages(conn.Buffers[channel].Outgoing)
	}
	conn.UI.tabStack.Select(buf.UI.TabItem())
}

func (conn *Connection) removeBuffer(channel string) {
	conn.UI.RemoveBuffer(conn.Buffers[channel])
	delete(conn.Buffers, channel)
}

func (conn *Connection) handleCommandFromBuffer(channel string) {
	var cmd message.Command
	shouldContinue := true
	for shouldContinue {
		cmd = <-conn.Buffers[channel].CommandOut
		switch cmd.Action {
		case "part":
			if channel[0] == '#' {
				conn.client.Part(channel)
				conn.autojoin = utils.RemoveAtIndex(utils.IndexOf(channel, conn.autojoin), conn.autojoin)
			}
			conn.removeBuffer(channel)
			shouldContinue = false
		}
	}
}
