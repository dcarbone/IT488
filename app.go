package main

import (
	"fmt"
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

	container *fyne.Container
}

func newTaskApp(fyneApp fyne.App, db *gorm.DB) *TaskApp {
	ta := TaskApp{
		db: db,
	}

	logo, err := GetFullSizeLogoPNG()
	if err != nil {
		panic(fmt.Sprintf("error reading logo: %v", err))
	}

	fyneApp.Settings().SetTheme(NewTheme())

	logoImg := canvas.NewImageFromImage(logo)
	logoImg.FillMode = canvas.ImageFillOriginal

	ta.container = container.NewStack(
		canvas.NewRectangle(ThemeBackgroundColor()),
		container.NewBorder(
			nil,
			widget.NewButton("Quit", func() { fyneApp.Quit() }),
			nil,
			nil,
			container.NewCenter(
				container.NewVBox(
					logoImg,
					widget.NewButton("Today's List", func() {

					}),
					widget.NewButton("Create List", func() {

					}),
				),
			),
		),
	)

	return &ta
}

func (ta *TaskApp) Container() *fyne.Container {
	ta.mu.Lock()
	defer ta.mu.Unlock()
	return ta.container
}
