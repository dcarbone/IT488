package main

import (
	"sync"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
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

	hdr     fyne.CanvasObject
	ftr     fyne.CanvasObject
	fly     fyne.CanvasObject
	content []fyne.CanvasObject

	root *fyne.Container
	body *fyne.Container
}

func newTaskApp(fyneApp fyne.App, db *gorm.DB) *TaskApp {
	ta := TaskApp{
		db:  db,
		ftr: widget.NewButton("Quit", func() { fyneApp.Quit() }),
	}

	hv := NewHomeView()
	hv.Foreground()

	ta.body = container.NewBorder(
		nil,
		ta.ftr,
		nil,
		nil,
		hv.Content(),
	)
	ta.root = container.NewStack(
		canvas.NewRectangle(ThemeBackgroundColor()),
		ta.body,
	)

	return &ta
}

func (ta *TaskApp) Root() *fyne.Container {
	return ta.root
}
