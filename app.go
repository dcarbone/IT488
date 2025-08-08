package main

import (
	"sync"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"gorm.io/gorm"
)

func logAppLifecycle(a fyne.App) {
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

type TaskApp struct {
	mu sync.RWMutex

	db *gorm.DB

	hdr  fyne.Canvas
	ftr  fyne.Canvas
	fly  fyne.Canvas
	main fyne.Canvas

	root *fyne.Container
}

func newTaskApp(db *gorm.DB) *TaskApp {
	ta := TaskApp{
		db: db,
	}

	ta.root = container.NewBorder(
		nil,
		nil,
		nil,
		nil,
	)
	return &ta
}

func (ta *TaskApp) Root() *fyne.Container {
	return ta.root
}
