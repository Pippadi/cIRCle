package widgets

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
)

type EnterCatchingEntry struct {
	widget.Entry
	OnEnter EnterHandler
}

type EnterHandler func()

func NewEnterCatchingEntry() *EnterCatchingEntry {
	e := &EnterCatchingEntry{}
	e.ExtendBaseWidget(e)
	return e
}

func (e *EnterCatchingEntry) TypedKey(key *fyne.KeyEvent) {
	if key.Name == fyne.KeyReturn {
		e.OnEnter()
	} else {
		e.Entry.TypedKey(key)
	}
}
