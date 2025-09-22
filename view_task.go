package main

import (
	"context"
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

var _ View = (*TaskView)(nil)

type TaskView struct {
	*baseView
	task     Task
	onDelete func()
}

func NewTaskView(ta *TaskApp, task Task, onDelete func()) *TaskView {
	v := TaskView{
		baseView: newBaseView("Task View", ta),
		task:     task,
		onDelete: onDelete,
	}
	return &v
}

func (v *TaskView) Title() string {
	return v.task.Label
}

func (v *TaskView) Foreground() fyne.CanvasObject {
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

	hdr := container.NewHBox(
		layout.NewSpacer(),
		newTaskPrioritySwitcherButton(v.app.DB(), &v.task),
		newTaskStatusSwitcherButton(v.app.DB(), &v.task),
	)

	body := container.NewVBox(
		FormLabel("List:"),
		widget.NewLabel(func() string {
			taskList := GetListForTask(ctx, v.app.DB(), v.task)
			if taskList == nil {
				return "None"
			}
			return taskList.Label
		}()),
		FormLabel("Due Date:"),
		widget.NewLabel(FormatDateTime(v.task.DueDate)),
		FormLabel("Description:"),
		container.NewHScroll(
			widget.NewRichTextFromMarkdown(v.task.Description),
		),
	)

	ftr := container.NewHBox(
		layout.NewSpacer(),
		widget.NewButtonWithIcon("Delete", theme.DeleteIcon(), func() {
			res := v.app.DB().Delete(v.task)
			if res.Error != nil {
				panic(fmt.Sprintf("Error deleting task %d: %v", v.task.ID, res.Error))
			}
			v.onDelete()
		}),
		widget.NewButtonWithIcon("Edit", IconEdit, func() {
			v.app.RenderMutateTaskView(&v.task, nil, v.onDelete)
		}),
	)

	return container.NewBorder(
		hdr,
		ftr,
		nil,
		nil,
		body,
	)
}

func (v *TaskView) Background() {
	v.mu.Lock()
	defer v.mu.Unlock()
	v.background()
}
