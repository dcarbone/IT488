package main

import (
	"context"
	"errors"
	"fmt"
	"image/color"
	"log/slog"
	"regexp"
	"strconv"
	"strings"
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
					ctx, cancel := context.WithCancel(context.Background())
					defer cancel()
					go func() {
						<-v.deactivated
						cancel()
					}()
					latestList, err := FindOneModel[TaskList](ctx, v.app.DB(), WithSort("Date desc"))
					if err != nil {
						log.Error("Error finding latest task list", "err", err)
						panic(fmt.Sprintf("Error finding latest task list: %v", err))
					}
					if latestList == nil {
						v.app.RenderCreateListView()
						return
					}
					v.app.RenderTaskListView(*latestList)
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
			if v.state != ViewStateForeground {
				v.mu.Unlock()
				return
			}
			tl, err := v.createList(nameInput.Text, descInput.Text)
			if err != nil {
				v.render(err)
				v.mu.Unlock()
				return
			}

			v.mu.Unlock()
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
	if !v.foreground() {
		v.mu.Unlock()
		return nil
	}
	v.mu.Unlock()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		<-v.deactivated
		cancel()
	}()

	hdr := container.NewHBox(
		HeaderCanvas(v.taskList.Label),
		widget.NewButtonWithIcon("", theme.Icon(theme.IconNameContentAdd), func() {
			v.app.RenderMutateTaskModal(&v.taskList)
		}),
	)

	log.Debug("Counting tasks...", "task_list", v.taskList.Label)
	taskCount, err := CountAssociation[TaskList](ctx, v.app.DB(), v.taskList, "Tasks")
	if err != nil {
		if errors.Is(err, context.Canceled) {
			return nil
		}
		log.Error("Error counting tasks in list", "list", v.taskList.Label, "err", err)
		panic(fmt.Sprintf("Error counting tasks in list %q: %v", v.taskList.Label, err))
	}

	log.Debug("Got task count", "task_list", v.taskList.Label, "count", taskCount)

	ftr := canvas.NewText(fmt.Sprintf("Total tasks: %d", taskCount), color.Black)

	log.Debug("Finding tasks...", "task_list", v.taskList.Label)

	tasks, err := FindAssociation[TaskList, Task](ctx, v.app.DB(), v.taskList, "Tasks")
	if err != nil {
		if errors.Is(err, context.Canceled) {
			return nil
		}
		log.Error("Error finding tasks in list", "list", v.taskList.Label, "err", err)
		panic(fmt.Sprintf("Error finding tasks in list %q: %v", v.taskList.Label, err))
	}

	log.Debug("Found tasks", "task_list", v.taskList.Label, "task_count", len(tasks))

	taskRows := make([]fyne.CanvasObject, 0)

	for _, task := range tasks {
		taskRows = append(taskRows, container.NewBorder(
			nil,
			canvas.NewText(fmt.Sprintf("Created: %s", task.CreatedAt), color.Black),

			widget.NewButtonWithIcon("", theme.Icon(theme.IconNameSettings), func() {
				// TODO: implement edit
			}),

			widget.NewButtonWithIcon("", theme.Icon(theme.IconNameDelete), func() {
				if v.state == ViewStateBackground {
					return
				}
				res := v.app.DB().Delete(task)
				if res.Error != nil {
					panic(fmt.Sprintf("Error deleting task %d: %v", task.ID, err))
				}
				v.app.RenderTaskListView(v.taskList)
			}),

			canvas.NewText(task.Label, color.Black),
		))
	}

	body := container.NewHScroll(container.NewHBox(taskRows...))

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

var _ View = (*MutateTaskView)(nil)

type MutateTaskView struct {
	*baseView
	taskList *TaskList
}

func NewMutateTaskView(app *TaskApp, taskList *TaskList) *MutateTaskView {
	v := MutateTaskView{
		baseView: newBaseView("Mutate Task Modal", app),
		taskList: taskList,
	}
	return &v
}

var taskSelectRe = regexp.MustCompile("\\((\\d+)\\)$")

func (v *MutateTaskView) Foreground() fyne.CanvasObject {
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

	allTaskLists, err := FindModel[TaskList](ctx, v.app.DB())
	if err != nil {
		log.Error("Error fetching task lists", "err", err)
		panic(fmt.Sprintf("Error fetching task lists: %v", err))
	}

	listNames := make([]string, 0)
	for _, tl := range allTaskLists {
		listNames = append(listNames, fmt.Sprintf("%s (%d)", tl.Label, tl.ID))
	}

	hdr := container.NewHBox(
		HeaderCanvas("Create New Task"),
		widget.NewButtonWithIcon("", theme.Icon(theme.IconNameCancel), v.app.RenderPreviousView),
	)

	titleLabel := canvas.NewText("Title", color.Black)
	titleInput := widget.NewEntry()
	titleInput.OnChanged = func(s string) {
		if len(s) > 50 {
			titleInput.SetText(s[:50])
		}
	}
	titleInput.PlaceHolder = "Task Title"

	chosenTaskList := v.taskList

	tlSelectLabel := canvas.NewText("Choose Task List", color.Black)
	tlSelect := widget.NewSelect(listNames, func(s string) {
		id, _ := strconv.ParseUint(taskSelectRe.FindStringSubmatch(s)[1], 10, 64)
		for _, tl := range allTaskLists {
			if uint(id) == tl.ID {
				chosenTaskList = &tl
				break
			}
		}
	})

	if chosenTaskList != nil {
		tlSelect.Selected = fmt.Sprintf("%s (%d)", chosenTaskList.Label, chosenTaskList.ID)
	}

	chosenStatus := TaskStatusTodo
	statusSelectLabel := canvas.NewText("Status", color.Black)
	statusSelect := widget.NewSelect(TaskStatuses, func(s string) {
		chosenStatus = TaskStatus(strings.ToLower(s))
	})
	statusSelect.Selected = strings.ToTitle(string(chosenStatus))

	descLabel := canvas.NewText("Description:", color.Black)
	descInput := widget.NewMultiLineEntry()
	descInput.PlaceHolder = "Task description in Markdown"
	descInput.OnChanged = func(s string) {
		if len(s) > 500 {
			descInput.SetText(s[:500])
		}
	}

	body := container.NewVBox(
		titleLabel,
		titleInput,

		tlSelectLabel,
		tlSelect,

		statusSelectLabel,
		statusSelect,

		descLabel,
		descInput,
	)

	ftr := widget.NewButtonWithIcon("Save", theme.Icon(theme.IconNameDocumentSave), func() {
		task := Task{
			Label:       titleLabel.Text,
			Description: descInput.Text,
			Status:      string(chosenStatus),
			TaskList:    chosenTaskList,
			Priority:    uint(taskPrioritySrc.Add(1)),
		}
		res := v.app.DB().Create(&task)
		if res.Error != nil {
			log.Error("Error saving task", "err", res.Error)
			panic(fmt.Sprintf("Error saving task: %v", res.Error))
		}

		v.app.RenderPreviousView()
	})

	content := container.NewBorder(
		hdr,
		ftr,
		nil,
		nil,
		body,
	)

	return content
}

func (v *MutateTaskView) Background() {
	v.mu.Lock()
	defer v.mu.Unlock()
	v.background()
}
