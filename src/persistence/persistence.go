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
	if cr, err := storage.CanRead(p.uri); cr {
		log.Println(cr)
		checkError(err)
		reader, err := storage.Reader(p.uri)
		checkError(err)
		bytes, err := ioutil.ReadAll(reader)
		checkError(err)
		json.Unmarshal(bytes, c)
		reader.Close()
	}
	return c
}

func (p *Persistor) DumpConnConfig(conf *connection.Config) {
	if cw, err := storage.CanWrite(p.uri); cw {
		log.Println(cw)
		checkError(err)
		bytes, err := json.MarshalIndent(conf, "", "\t")
		checkError(err)
		writer, err := storage.Writer(p.uri)
		checkError(err)
		log.Println(writer.Write(bytes))
		writer.Close()
	}
}

func checkError(err error) {
	if err != nil {
		log.Println(err)
	}
}
