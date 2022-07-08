package connection

import "strings"

func (conn *Connection) loadConfig() {
	prefs := conn.app.Preferences()
	conn.UI.AddrEntry.SetText(prefs.String("Host"))
	conn.UI.PortEntry.SetText(prefs.String("Port"))
	conn.UI.PassEntry.SetText(prefs.String("Pass"))
	conn.UI.NickEntry.SetText(prefs.String("Nick"))
}

func (conn *Connection) dumpConfig() {
	prefs := conn.app.Preferences()
	prefs.SetString("Host", conn.UI.AddrEntry.Text)
	prefs.SetString("Port", conn.UI.PortEntry.Text)
	prefs.SetString("Nick", conn.UI.NickEntry.Text)
	prefs.SetString("Pass", conn.UI.PassEntry.Text)
}

func (conn *Connection) autojoin() {
	chanListStr := conn.app.Preferences().String("Channels")
	if chanListStr != "" {
		channels := strings.Split(chanListStr, " ")
		for _, channel := range channels {
			conn.join(channel)
		}
	}
}

func (conn *Connection) dumpAutojoin() {
	channels := ""
	for c, _ := range conn.Buffers {
		if c[0] == '#' {
			channels = c + " "
		}
	}
	if channels != "" {
		conn.app.Preferences().SetString("Channels", channels[:len(channels)-1])
	}
}

func (conn *Connection) PersistSettings() {
	conn.dumpAutojoin()
	conn.dumpConfig()
}
