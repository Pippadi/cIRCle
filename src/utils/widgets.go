package utils

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
)

// NewEntryButtonContainer creates a container with the passed button sitting to the right of the entry,
// with the entry taking up as much space as possible
func NewEntryButtonContainer(entry fyne.CanvasObject, button fyne.CanvasObject) *fyne.Container {
	return container.New(layout.NewBorderLayout(nil, nil, nil, button), button, entry)
}
