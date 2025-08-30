package main

import (
	"sync"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
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
	window  fyne.Window
	db      *gorm.DB

	container      *fyne.Container
	body           *fyne.Container
	contentWrapper *fyne.Container
	activeView     View
	previousView   View

	showNavBtn *widget.Button
}

func newTaskApp(fyneApp fyne.App, window fyne.Window, db *gorm.DB) *TaskApp {
	ta := TaskApp{
		fyneApp: fyneApp,
		window:  window,
		db:      db,
	}

	fyneApp.Settings().SetTheme(NewTheme())

	ta.contentWrapper = container.NewStack()

	ta.showNavBtn = widget.NewButtonWithIcon("", theme.Icon(theme.IconNameList), func() {
		ta.RenderNavigation()
	})

	ta.body = container.NewBorder(
		container.NewHBox(ta.showNavBtn),
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
		ta.previousView = ta.activeView
	}
	if ta.showNavBtn.Hidden {
		ta.showNavBtn.Show()
	}
	ta.contentWrapper.RemoveAll()
	ta.activeView = view
	ta.contentWrapper.Add(ta.activeView.Foreground())
}

func (ta *TaskApp) RenderPreviousView() {
	ta.mu.Lock()
	defer ta.mu.Unlock()
	if ta.previousView != nil {
		ta.renderView(ta.previousView)
	} else {
		ta.renderView(NewHomeView(ta))
	}
}

func (ta *TaskApp) RenderNavigation() {
	ta.mu.Lock()
	defer ta.mu.Unlock()
	ta.renderView(NewNavigationView(ta))
	ta.showNavBtn.Hide()
}

func (ta *TaskApp) RenderHomeView() {
	ta.mu.Lock()
	defer ta.mu.Unlock()
	ta.renderView(NewHomeView(ta))
}

func (ta *TaskApp) RenderMutateTaskListView(taskList *TaskList) {
	ta.mu.Lock()
	defer ta.mu.Unlock()
	ta.renderView(NewMutateTaskListView(ta, taskList))
}

func (ta *TaskApp) RenderTaskListsView() {
	ta.mu.Lock()
	defer ta.mu.Unlock()
	ta.renderView(NewTaskListsView(ta))
}

func (ta *TaskApp) RenderTaskListView(title string, opts ...ModelQueryOpt) {
	ta.mu.Lock()
	defer ta.mu.Unlock()
	ta.renderView(NewTaskListView(ta, title, opts...))
}

func (ta *TaskApp) RenderMutateTaskView(task *Task, taskList *TaskList) {
	ta.mu.Lock()
	defer ta.mu.Unlock()
	ta.renderView(NewMutateTaskView(ta, task, taskList))
}

func (ta *TaskApp) Container() *fyne.Container {
	ta.mu.Lock()
	defer ta.mu.Unlock()
	return ta.container
}

func (ta *TaskApp) DB() *gorm.DB {
	return ta.db
}
