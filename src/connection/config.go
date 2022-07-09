package connection

import (
	"strconv"
)

type Config struct {
	Host     string
	Port     int
	Nick     string
	Pass     string
	Channels []string
}

func (conn *Connection) LoadConfig(conf *Config) {
	conn.UI.AddrEntry.SetText(conf.Host)
	conn.UI.PortEntry.SetText(strconv.Itoa(conf.Port))
	conn.UI.NickEntry.SetText(conf.Nick)
	conn.UI.PassEntry.SetText(conf.Pass)
	conn.autojoin = conf.Channels
}

func (conn *Connection) GetConfig() *Config {
	conf := Config{}
	conf.Host = conn.UI.AddrEntry.Text
	conf.Port, _ = strconv.Atoi(conn.UI.PortEntry.Text)
	conf.Nick = conn.UI.NickEntry.Text
	conf.Pass = conn.UI.PassEntry.Text
	conf.Channels = conn.autojoin
	return &conf
}
