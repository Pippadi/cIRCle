package widgets

import "fyne.io/fyne/v2/widget"

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
