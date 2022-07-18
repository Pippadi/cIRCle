package widgets

import (
	"errors"

	"fyne.io/fyne/v2/widget"
)

// NickEntry is an entry meant for nicknames
type NickEntry struct {
	widget.Entry
}

func NewNickEntry() *NickEntry {
	ne := NickEntry{}
	ne.ExtendBaseWidget(&ne)
	ne.SetPlaceHolder("Nickname")
	ne.Validator = validNickString
	return &ne
}

func validNickString(nick string) error {
	if len(nick) < 3 {
		return errors.New("Nickname must be three characters or more")
	}
	return nil
}
