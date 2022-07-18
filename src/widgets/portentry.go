package widgets

import (
	"errors"
	"strconv"
)

type PortEntry struct {
	NumericEntry
}

func NewPortEntry() *PortEntry {
	pe := PortEntry{}
	pe.ExtendBaseWidget(&pe)
	pe.SetPlaceHolder("Port")
	pe.Validator = validPortString
	return &pe
}

func validPortString(p string) error {
	port, err := strconv.Atoi(p)
	if err != nil {
		return errors.New("Port must be numeric")
	}
	if port > 65535 || port < 0 {
		return errors.New("Port out of range")
	}
	return nil
}
