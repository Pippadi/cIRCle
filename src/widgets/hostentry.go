package widgets

import (
	"errors"

	"fyne.io/fyne/v2/widget"
)

// HostEntry is an entry meant for hostnames/addresses
type HostEntry struct {
	widget.Entry
}

func NewHostEntry() *HostEntry {
	he := HostEntry{}
	he.ExtendBaseWidget(&he)
	he.SetPlaceHolder("Address")
	he.Validator = validAddrString
	return &he
}

func validAddrString(addr string) error {
	if addr == "" {
		return errors.New("Address cannot be empty")
	}
	return nil
}
