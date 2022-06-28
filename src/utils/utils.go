package utils

import "fyne.io/fyne/v2"

func SetWidgetsActive(active bool, widgets []fyne.Disableable) {
	for _, w := range widgets {
		if active {
			w.Enable()
		} else {
			w.Disable()
		}
	}
}
