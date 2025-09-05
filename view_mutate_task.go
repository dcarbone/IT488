package main

import (
	"context"
	"fmt"
	"image/color"
	"regexp"
	"strconv"
	"strings"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/sdassow/fyne-datepicker"
	"gorm.io/gorm"
)

var (
	_ View = (*MutateTaskView)(nil)

	taskSelectRe = regexp.MustCompile("\\((\\d+)\\)$")
)

type MutateTaskView struct {
	*baseView
	task     *Task
	taskList *TaskList
}

func NewMutateTaskView(app *TaskApp, task *Task, taskList *TaskList) *MutateTaskView {
	v := MutateTaskView{
		baseView: newBaseView("Mutate Task Modal", app),
		task:     task,
		taskList: taskList,
	}
	return &v
}

func (v *MutateTaskView) Title() string {
	if v.task == nil {
		return "Create new task"
	}
	return fmt.Sprintf("Edit task %s", v.task.Label)
}

func (v *MutateTaskView) Foreground() fyne.CanvasObject {
	v.mu.Lock()
	defer v.mu.Unlock()
	if !v.foreground() {
		return nil
	}

	var (
		content *fyne.Container
	)

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

	titleLabel := canvas.NewText("Title", color.Black)
	titleInput := widget.NewEntry()
	if v.task != nil {
		titleInput.SetText(v.task.Label)
	}
	titleInput.OnChanged = func(s string) {
		if len(s) > 50 {
			titleInput.SetText(s[:50])
		}
	}
	titleInput.PlaceHolder = "Task Title"

	chosenTaskList := v.taskList

	tlSelectLabel := canvas.NewText("Choose Task List", color.Black)
	tlSelect := widget.NewSelect(listNames, func(s string) {
		var id uint
		if v.task != nil {
			id = uint(v.task.TaskListID)
		} else {
			idd, _ := strconv.ParseUint(taskSelectRe.FindStringSubmatch(s)[1], 10, 64)
			id = uint(idd)
		}
		for _, tl := range allTaskLists {
			if id == tl.ID {
				chosenTaskList = &tl
				break
			}
		}
	})

	if chosenTaskList != nil {
		tlSelect.Selected = fmt.Sprintf("%s (%d)", chosenTaskList.Label, chosenTaskList.ID)
	}

	chosenStatus := TaskStatusTodo
	if v.task != nil {
		chosenStatus = v.task.Status
	}
	statusSelectLabel := widget.NewLabel("Status")
	statusSelect := widget.NewSelect(TaskStatuses, func(s string) {
		chosenStatus = strings.ToLower(s)
	})
	statusSelect.Selected = strings.ToTitle(chosenStatus)

	var priorityContainer *fyne.Container
	chosenPriority := TaskPriorityHigh
	chosenPriorityImageContainer := container.NewStack(
		GetAssetImageCanvas(
			GetConstrainedImage(TaskPriorityImage(chosenPriority), 50),
		),
	)
	if v.task != nil {
		chosenPriority = TaskPriorityName(v.task.UserPriority)
	}
	prioritySelectLabel := widget.NewLabel("Priority")
	prioritySelect := widget.NewSelect(TaskPriorities, func(s string) {
		chosenPriority = strings.ToLower(s)
		chosenPriorityImageContainer.RemoveAll()
		chosenPriorityImageContainer.Add(
			GetAssetImageCanvas(
				GetConstrainedImage(TaskPriorityImage(chosenPriority), 50),
			),
		)
	})
	prioritySelect.Selected = strings.ToTitle(chosenPriority)
	priorityContainer = container.NewBorder(
		nil,
		nil,
		nil,
		chosenPriorityImageContainer,
		prioritySelect,
	)

	chosenDueDate := time.Now()
	if v.task != nil && !v.task.DueDate.IsZero() {
		chosenDueDate = v.task.DueDate
	}
	dtpLabel := FormLabel("Due Date:")
	dueDateDisplay := widget.NewLabel(FormatDateTime(chosenDueDate))
	var datePickerModal *widget.PopUp
	dtp := datepicker.NewDateTimePicker(chosenDueDate, time.Sunday, func(t time.Time, b bool) {
		chosenDueDate = t
		dueDateDisplay.SetText(FormatDateTime(chosenDueDate))
	})
	dtpSaveBtn := widget.NewButtonWithIcon("Save", theme.Icon(theme.IconNameDocumentSave), func() {
		dtp.OnActioned(true)
		datePickerModal.Hide()
	})
	dtpContainer := container.NewBorder(
		container.NewBorder(
			nil,
			nil,
			nil,
			widget.NewButtonWithIcon("", theme.Icon(theme.IconNameCancel), func() {
				datePickerModal.Hide()
			}),
		),
		container.NewBorder(
			nil,
			nil,
			nil,
			dtpSaveBtn,
		),
		nil,
		nil,
		dtp,
	)
	datePickerModal = widget.NewModalPopUp(
		dtpContainer,
		v.app.window.Canvas(),
	)

	dtpButton := widget.NewButtonWithIcon("", theme.Icon(theme.IconNameCalendar), datePickerModal.Show)
	dueDateContainer := container.NewBorder(nil, nil, dueDateDisplay, dtpButton)

	descLabel := FormLabel("Description:")
	descInput := widget.NewMultiLineEntry()
	if v.task != nil {
		descInput.SetText(v.task.Description)
	}
	descInput.PlaceHolder = "Task description in Markdown"
	descInput.OnChanged = func(s string) {
		if len(s) > 500 {
			descInput.SetText(s[:500])
		}
	}
	descInput.SetMinRowsVisible(10)

	body := container.NewVScroll(
		container.NewVBox(
			titleLabel,
			titleInput,

			tlSelectLabel,
			tlSelect,

			statusSelectLabel,
			statusSelect,

			prioritySelectLabel,
			priorityContainer,

			dtpLabel,
			dueDateContainer,

			descLabel,
			descInput,
		),
	)

	ftr := container.NewHBox(
		layout.NewSpacer(),
		widget.NewButtonWithIcon("Cancel", theme.Icon(theme.IconNameCancel), v.app.RenderPreviousView),
		widget.NewButtonWithIcon("Save", theme.Icon(theme.IconNameDocumentSave), func() {
			var res *gorm.DB
			if v.task != nil {
				v.task.Label = titleInput.Text
				v.task.Description = descInput.Text
				v.task.Status = chosenStatus
				v.task.UserPriority = TaskPriorityNumber(chosenPriority)
				v.task.TaskList = chosenTaskList
				v.task.DueDate = chosenDueDate
				res = v.app.DB().Updates(v.task)
			} else {
				task := Task{
					Label:        titleInput.Text,
					Description:  descInput.Text,
					Status:       chosenStatus,
					UserPriority: TaskPriorityNumber(chosenPriority),
					TaskList:     chosenTaskList,
					DueDate:      chosenDueDate,
					Priority:     GetNextTaskOrderNum(),
				}
				res = v.app.DB().Create(&task)
			}
			if res.Error != nil {
				log.Error("Error saving task", "err", res.Error)
				panic(fmt.Sprintf("Error saving task: %v", res.Error))
			}

			v.app.RenderPreviousView()
		}),
	)

	content = container.NewBorder(
		nil,
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
