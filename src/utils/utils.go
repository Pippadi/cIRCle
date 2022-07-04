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

func RemoveAtIndex(index int, slice []string) []string {
	slice[index] = slice[len(slice)-1]
	return slice[:len(slice)-1]
}

func IndexOf(element string, slice []string) int {
	for index, el := range slice {
		if el == element {
			return index
		}
	}
	return -1
}

func In(element string, slice []string) bool {
	for _, el := range slice {
		if el == element {
			return true
		}
	}
	return false
}
