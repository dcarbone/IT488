package main

import (
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
	Content() fyne.CanvasObject
	Foreground()
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
	content fyne.CanvasObject
}

func NewHomeView(app *TaskApp) *HomeView {
	v := HomeView{
		baseView: newBaseView("home", app),
		content:  widget.NewLabel("Loading..."),
	}

	logo, err := GetFullSizeLogoPNG()
	if err != nil {
		panic(fmt.Sprintf("error reading logo: %v", err))
	}
	v.logoImg = canvas.NewImageFromImage(logo)
	v.logoImg.FillMode = canvas.ImageFillOriginal

	return &v
}

func (v *HomeView) Content() fyne.CanvasObject {
	v.mu.Lock()
	defer v.mu.Unlock()
	return v.content
}

func (v *HomeView) Foreground() {
	v.mu.Lock()
	defer v.mu.Unlock()
	if v.foreground() {
		v.content = container.NewCenter(
			container.NewVBox(
				v.logoImg,
				widget.NewButton("Today's List", func() {

				}),
				widget.NewButton("Create List", v.app.RenderCreateListView),
			),
		)
	}
}

func (v *HomeView) Background() {
	v.mu.Lock()
	defer v.mu.Unlock()
	if v.background() {
		v.content = widget.NewLabel("Loading...")
	}
}

var _ View = (*CreateTaskListView)(nil)

type CreateTaskListView struct {
	*baseView
	content fyne.CanvasObject
}

func NewCreateTaskListView(app *TaskApp) *CreateTaskListView {
	v := CreateTaskListView{
		baseView: newBaseView("Create Task List", app),
		content:  widget.NewLabel("Loading..."),
	}
	return &v
}

func (v *CreateTaskListView) Content() fyne.CanvasObject {
	v.mu.Lock()
	defer v.mu.Unlock()
	return v.content
}

func (v *CreateTaskListView) Foreground() {
	v.mu.Lock()
	defer v.mu.Unlock()
	if v.foreground() {
		v.render(nil)
	}
}

func (v *CreateTaskListView) Background() {
	v.mu.Lock()
	defer v.mu.Unlock()
	if v.background() {
		v.content = widget.NewLabel("Loading...")
	}
}

func (v *CreateTaskListView) render(err error) {
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

	v.content = content
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
	content  fyne.CanvasObject
}

func NewTaskListView(app *TaskApp, taskList TaskList) *TaskListView {
	v := TaskListView{
		baseView: newBaseView(fmt.Sprintf("Task List %s", taskList.Label), app),
		taskList: taskList,
		content:  widget.NewLabel("Loading..."),
	}
	return &v
}

func (v *TaskListView) Content() fyne.CanvasObject {
	v.mu.Lock()
	defer v.mu.Unlock()
	return v.content
}

func (v *TaskListView) Foreground() {
	v.mu.Lock()
	defer v.mu.Unlock()
	if v.foreground() {

	}
}

func (v *TaskListView) Background() {
	v.mu.Lock()
	defer v.mu.Unlock()
	if v.background() {

	}
}
