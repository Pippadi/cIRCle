package connection

func (conn *Connection) chat() {
	conn.Join(conn.UI.JoinEntry.Text)
	conn.UI.JoinEntry.SetText("")
}
