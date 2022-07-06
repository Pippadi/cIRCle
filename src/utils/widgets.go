package utils

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

// NewEntryButtonContainer creates a container with the passed button sitting to the right of the entry,
// with the entry taking up as much space as possible
func NewEntryButtonContainer(entry fyne.CanvasObject, button fyne.CanvasObject) *fyne.Container {
	return container.New(layout.NewBorderLayout(nil, nil, nil, button), button, entry)
}

type NumericEntry struct {
	widget.Entry
}

func NewNumericEntry() *NumericEntry {
	ne := &NumericEntry{}
	ne.ExtendBaseWidget(ne)
	return ne
}

func (ne *NumericEntry) TypedRune(r rune) {
	switch r {
	case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
		ne.Entry.TypedRune(r)
	}
}

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
