package persistence

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"path"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/storage"
	"github.com/Pippadi/cIRCle/src/connection"
)

type Persistor struct {
	uri fyne.URI
}

func New(a fyne.App) *Persistor {
	return &Persistor{storage.NewFileURI(path.Join(a.Storage().RootURI().Path(), "prefs.json"))}
}

func (p *Persistor) LoadConnConfig() *connection.Config {
	c := &connection.Config{Port: 6667}
	configs := make([]connection.Config, 0) // To allow for multiple connections in the future
	ex, err := storage.Exists(p.uri)
	checkError(err)
	cr, err := storage.CanRead(p.uri)
	checkError(err)
	if ex && cr {
		reader, err := storage.Reader(p.uri)
		checkError(err)
		bytes, err := ioutil.ReadAll(reader)
		checkError(err)
		checkError(json.Unmarshal(bytes, &configs))
		checkError(reader.Close())
	}
	if len(configs) > 0 {
		c = &(configs[0])
	}
	return c
}

func (p *Persistor) DumpConnConfig(conf *connection.Config) {
	configs := make([]connection.Config, 1)
	configs[0] = *conf
	if cw, err := storage.CanWrite(p.uri); cw {
		checkError(err)
		bytes, err := json.MarshalIndent(&configs, "", "\t")
		checkError(err)
		writer, err := storage.Writer(p.uri)
		checkError(err)
		writer.Write(bytes)
		writer.Close()
	}
}

func checkError(err error) {
	if err != nil {
		log.Println(err)
	}
}
