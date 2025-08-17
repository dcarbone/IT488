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

	fyneApp fyne.App
	db      *gorm.DB

	container      *fyne.Container
	body           *fyne.Container
	contentWrapper *fyne.Container
	activeView     View
}

func newTaskApp(fyneApp fyne.App, db *gorm.DB) *TaskApp {
	ta := TaskApp{
		fyneApp: fyneApp,
		db:      db,
	}

	fyneApp.Settings().SetTheme(NewTheme())

	ta.contentWrapper = container.NewStack()

	ta.body = container.NewBorder(
		nil,
		widget.NewButton("Quit", func() { fyneApp.Quit() }),
		nil,
		nil,
		ta.contentWrapper,
	)

	ta.container = container.NewStack(
		canvas.NewRectangle(ThemeBackgroundColor()),
		ta.body,
	)

	return &ta
}

func (ta *TaskApp) renderView(view View) {
	if ta.activeView != nil {
		ta.activeView.Background()
	}
	ta.contentWrapper.RemoveAll()
	ta.activeView = view
	ta.contentWrapper.Add(ta.activeView.Foreground())
}

func (ta *TaskApp) RenderHomeView() {
	ta.mu.Lock()
	defer ta.mu.Unlock()

	if _, ok := ta.activeView.(*HomeView); ok {
		return
	}

	ta.renderView(NewHomeView(ta))
}

func (ta *TaskApp) RenderCreateListView() {
	ta.mu.Lock()
	defer ta.mu.Unlock()

	if _, ok := ta.activeView.(*CreateTaskListView); ok {
		return
	}

	ta.renderView(NewCreateTaskListView(ta))
}

func (ta *TaskApp) RenderTaskListView(taskList TaskList) {
	ta.mu.Lock()
	defer ta.mu.Unlock()

	//if tlv, ok := ta.activeView.(*TaskListView); ok && tlv.taskList.ID != taskList.ID {
	//	return
	//}
	//
	//ta.renderView(NewTaskListView(ta, taskList))
}

func (ta *TaskApp) Container() *fyne.Container {
	ta.mu.Lock()
	defer ta.mu.Unlock()
	return ta.container
}

func (ta *TaskApp) DB() *gorm.DB {
	return ta.db
}
