package main

import (
	"fyne.io/fyne/v2"
)

func logLifecycle(a fyne.App) {
	a.Lifecycle().SetOnStarted(func() {
		log.Debug("Lifecycle: Started")
	})
	a.Lifecycle().SetOnStopped(func() {
		log.Debug("Lifecycle: Stopped")
	})
	a.Lifecycle().SetOnEnteredForeground(func() {
		log.Debug("Lifecycle: Entered Foreground")
	})
	a.Lifecycle().SetOnExitedForeground(func() {
		log.Debug("Lifecycle: Exited Foreground")
	})
}
