package main

import (
	"context"
	"errors"
	"fmt"
	"image/color"
	"log/slog"
	"sync"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

type ViewState int

const (
	ViewStateForeground ViewState = iota
	ViewStateBackground
)

type View interface {
	Name() string
	State() ViewState
	Foreground() fyne.CanvasObject
	Background()
}

type baseView struct {
	mu          sync.Mutex
	log         *slog.Logger
	state       ViewState
	name        string
	app         *TaskApp
	deactivated chan struct{}
	children    []View
}

func newBaseView(name string, app *TaskApp) *baseView {
	v := baseView{
		log:      log.With("view", name),
		state:    ViewStateBackground,
		name:     name,
		app:      app,
		children: make([]View, 0),
	}
	return &v
}

func (v *baseView) Name() string {
	return v.name
}

func (v *baseView) State() ViewState {
	v.mu.Lock()
	defer v.mu.Unlock()
	return v.state
}

func (v *baseView) foreground() bool {
	if v.state == ViewStateForeground {
		v.log.Debug("View already in foreground")
		return false
	}
	v.log.Debug("Setting view state to foreground")
	v.state = ViewStateForeground
	v.deactivated = make(chan struct{})
	return true
}

func (v *baseView) background() bool {
	if v.state == ViewStateBackground {
		v.log.Debug("View already in background")
		return false
	}
	v.log.Debug("Setting view state to background")
	v.state = ViewStateBackground
	for _, child := range v.children {
		v.log.Debug("Backgrounding view child...", "child", child.Name())
		child.Background()
	}
	close(v.deactivated)
	v.children = make([]View, 0)
	return true
}

var _ View = (*HomeView)(nil)

type HomeView struct {
	*baseView

	logoImg *canvas.Image
}

func NewHomeView(app *TaskApp) *HomeView {
	v := HomeView{
		baseView: newBaseView("home", app),
	}

	logo, err := GetFullSizeLogoPNG()
	if err != nil {
		panic(fmt.Sprintf("error reading logo: %v", err))
	}
	v.logoImg = canvas.NewImageFromImage(logo)
	v.logoImg.FillMode = canvas.ImageFillOriginal

	return &v
}

func (v *HomeView) Foreground() fyne.CanvasObject {
	v.mu.Lock()
	defer v.mu.Unlock()
	if v.foreground() {
		return container.NewCenter(
			container.NewVBox(
				v.logoImg,
				widget.NewButton("Today's List", func() {

				}),
				widget.NewButton("Create List", v.app.RenderCreateListView),
			),
		)
	}
	return nil
}

func (v *HomeView) Background() {
	v.mu.Lock()
	defer v.mu.Unlock()
	v.background()
}

var _ View = (*CreateTaskListView)(nil)

type CreateTaskListView struct {
	*baseView
}

func NewCreateTaskListView(app *TaskApp) *CreateTaskListView {
	v := CreateTaskListView{
		baseView: newBaseView("Create Task List", app),
	}
	return &v
}

func (v *CreateTaskListView) Foreground() fyne.CanvasObject {
	v.mu.Lock()
	defer v.mu.Unlock()
	if v.foreground() {
		return v.render(nil)
	}
	return nil
}

func (v *CreateTaskListView) Background() {
	v.mu.Lock()
	defer v.mu.Unlock()
	v.background()
}

func (v *CreateTaskListView) render(err error) fyne.CanvasObject {
	content := container.NewVBox()

	hdr := canvas.NewText("Create New List", color.Black)
	hdr.Alignment = fyne.TextAlignCenter
	hdr.TextStyle = fyne.TextStyle{Bold: true}
	hdr.TextSize = 32

	content.Add(hdr)

	if err != nil {
		content.Add(canvas.NewText("Error:", ColorRed))
		content.Add(canvas.NewText(err.Error(), ColorRed))
	}

	content.Add(canvas.NewText("Name:", color.Black))

	nameInput := widget.NewEntry()
	nameInput.OnChanged = func(s string) {
		if len(s) > 50 {
			nameInput.SetText(s[:50])
		}
	}

	content.Add(nameInput)

	content.Add(canvas.NewText("Description:", color.Black))

	descInput := widget.NewMultiLineEntry()
	descInput.PlaceHolder = "Enter Markdown formatted text."
	descInput.OnChanged = func(s string) {
		if len(s) > 500 {
			descInput.SetText(s[:500])
		}
	}

	content.Add(descInput)

	createBtn := widget.NewButtonWithIcon(
		"Save",
		theme.DocumentSaveIcon(),
		func() {
			v.mu.Lock()
			defer v.mu.Unlock()
			if v.state != ViewStateForeground {
				return
			}
			tl, err := v.createList(nameInput.Text, descInput.Text)
			if err != nil {
				v.render(err)
				return
			}
			v.app.RenderTaskListView(tl)
		},
	)

	content.Add(createBtn)

	return content
}

func (v *CreateTaskListView) createList(name, description string) (TaskList, error) {
	tl := TaskList{
		Label:       name,
		Date:        time.Now(),
		Description: description,
	}
	res := v.app.DB().Create(&tl)
	if res.Error != nil {
		log.Error("Error creating list", "err", res.Error)
	} else {
		log.Info("New task list created", "id", tl.ID, "name", tl.Label)
	}
	return tl, res.Error
}

var _ View = (*TaskListView)(nil)

type TaskListView struct {
	*baseView
	taskList TaskList
}

func NewTaskListView(app *TaskApp, taskList TaskList) *TaskListView {
	v := TaskListView{
		baseView: newBaseView(fmt.Sprintf("Task List %s", taskList.Label), app),
		taskList: taskList,
	}
	return &v
}

func (v *TaskListView) Foreground() fyne.CanvasObject {
	v.mu.Lock()
	defer v.mu.Unlock()

	if !v.foreground() {
		return nil
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		<-v.deactivated
		cancel()
	}()

	hdr := HeaderCanvas(v.taskList.Label)

	taskCount, err := CountAssociation[TaskList](ctx, v.app.DB(), "Tasks")
	if err != nil {
		if errors.Is(err, context.Canceled) {
			return nil
		}
		log.Error("Error counting tasks in list", "list", v.taskList.Label, "err", err)
		panic(fmt.Sprintf("Error counting tasks in list %q: %v", v.taskList.Label, err))
	}

	ftr := canvas.NewText(fmt.Sprintf("Total tasks: %d", taskCount), color.Black)

	tasks, err := FindAssociation[TaskList, Task](ctx, v.app.DB(), "Tasks")
	if err != nil {
		if errors.Is(err, context.Canceled) {
			return nil
		}
		log.Error("Error finding tasks in list", "list", v.taskList.Label, "err", err)
		panic(fmt.Sprintf("Error finding tasks in list %q: %v", v.taskList.Label, err))
	}

	taskViews := make([]fyne.CanvasObject, taskCount)
	for i := range tasks {
		taskViews[i] = container.NewBorder(
			nil,
			canvas.NewText(fmt.Sprintf("Created: %s", tasks[i].CreatedAt), color.Black),
			nil,
			widget.NewButtonWithIcon("", theme.Icon(theme.IconNameDelete), func() {

			}),
			canvas.NewText(tasks[i].Label, color.Black),
		)
	}

	body := container.NewHScroll(container.NewHBox(taskViews...))

	return container.NewBorder(
		hdr,
		ftr,
		nil,
		nil,
		body,
	)
}

func (v *TaskListView) Background() {
	v.mu.Lock()
	defer v.mu.Unlock()
	v.background()
}
