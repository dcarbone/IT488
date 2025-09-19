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
	viewTitle  *canvas.Text
}

func newTaskApp(fyneApp fyne.App, window fyne.Window, db *gorm.DB) *TaskApp {
	ta := TaskApp{
		fyneApp: fyneApp,
		window:  window,
		db:      db,
	}

	fyneApp.Settings().SetTheme(NewTheme())

	ta.showNavBtn = widget.NewButtonWithIcon("", theme.ListIcon(), func() {
		ta.RenderNavigation()
	})
	ta.viewTitle = HeaderCanvas("")

	ta.contentWrapper = container.NewStack()

	ta.body = container.NewBorder(
		container.NewHBox(ta.showNavBtn, ta.viewTitle),
		widget.NewButton("Quit", func() { fyneApp.Quit() }),
		nil,
		nil,
		ta.contentWrapper,
	)

	ta.container = container.NewStack(
		canvas.NewRectangle(ColorBackground),
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
	ta.viewTitle.Text = view.Title()
	ta.contentWrapper.Add(ta.activeView.Foreground())
}

func (ta *TaskApp) PreviousView() View {
	ta.mu.Lock()
	defer ta.mu.Unlock()
	return ta.previousView
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

func (ta *TaskApp) RenderListOfTasksView(title string, taskList *TaskList, opts ...ModelQueryOpt) {
	ta.mu.Lock()
	defer ta.mu.Unlock()
	ta.renderView(NewListOfTasksView(ta, title, taskList, opts...))
}

func (ta *TaskApp) RenderMutateTaskView(task *Task, taskList *TaskList, onDelete func()) {
	ta.mu.Lock()
	defer ta.mu.Unlock()
	ta.renderView(NewMutateTaskView(ta, task, taskList, onDelete))
}

func (ta *TaskApp) RenderTaskView(task Task, onDelete func()) {
	ta.mu.Lock()
	defer ta.mu.Unlock()
	ta.renderView(NewTaskView(ta, task, onDelete))
}

//func (ta *TaskApp) RenderTaskListView(taskList TaskList, onDelete func())  {
//	ta.mu.Lock()
//	defer ta.mu.Unlock()
//	ta.renderView()
//}

func (ta *TaskApp) Container() *fyne.Container {
	ta.mu.Lock()
	defer ta.mu.Unlock()
	return ta.container
}

func (ta *TaskApp) DB() *gorm.DB {
	return ta.db
}
