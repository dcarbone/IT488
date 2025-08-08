package main

import (
	"context"
	"sync"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
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

	hdr  fyne.CanvasObject
	ftr  fyne.CanvasObject
	fly  fyne.CanvasObject
	main fyne.CanvasObject

	root *fyne.Container
}

func newTaskApp(ctx context.Context, db *gorm.DB) *TaskApp {
	_, cancel := context.WithCancel(ctx)
	ta := TaskApp{
		db:   db,
		ftr:  widget.NewButton("Quit", func() { cancel() }),
		main: NewHomeView().Content(),
	}

	ta.root = container.NewBorder(
		nil,
		ta.ftr,
		nil,
		nil,
	)
	return &ta
}

func (ta *TaskApp) Root() *fyne.Container {
	return ta.root
}
