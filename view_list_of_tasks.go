package main

import (
	"context"
	"errors"
	"fmt"
	"image/color"
	"slices"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

func buildListOfTasksList(app *TaskApp, taskList *TaskList, tasks []Task, onDelete func()) fyne.CanvasObject {
	return widget.NewList(
		func() int {
			return len(tasks)
		},
		func() fyne.CanvasObject {
			return container.NewStack(widget.NewLabel("Loading..."))
		},
		func(id widget.ListItemID, object fyne.CanvasObject) {
			task := tasks[id]

			content := object.(*fyne.Container)

			var statusButton *widget.Button
			statusIdx := slices.Index(TaskStatuses, task.Status)
			statusButton = widget.NewButtonWithIcon(
				"",
				TaskStatusResource(task.Status),
				func() {
					statusIdx++
					if statusIdx == len(TaskStatuses) {
						statusIdx = 0
					}
					task.Status = TaskStatuses[statusIdx]
					res := app.DB().Model(&task).Update("Status", task.Status)
					if res.Error != nil {
						panic(fmt.Sprintf("error updating task status: %v", res.Error))
					}
					statusButton.SetIcon(TaskStatusResource(TaskStatuses[statusIdx]))
				},
			)
			statusButton.Importance = widget.LowImportance

			var priorityButton *widget.Button
			priorityIdx := slices.Index(TaskPriorities, strings.ToTitle(TaskPriorityName(task.UserPriority)))
			priorityButton = widget.NewButtonWithIcon(
				"",
				TaskPriorityResource(TaskPriorityName(task.UserPriority)),
				func() {
					priorityIdx++
					if priorityIdx == len(TaskPriorities) {
						priorityIdx = 0
					}
					task.UserPriority = TaskPriorityNumber(TaskPriorities[priorityIdx])
					res := app.DB().Model(&task).Update("UserPriority", task.UserPriority)
					if res.Error != nil {
						panic(fmt.Sprintf("error updating task user priority: %v", res.Error))
					}
					priorityButton.SetIcon(TaskPriorityResource(TaskPriorities[priorityIdx]))
				},
			)
			priorityButton.Importance = widget.LowImportance

			content.RemoveAll()
			content.Add(container.NewBorder(
				nil,
				nil,
				container.NewHBox(statusButton, priorityButton),
				container.NewHBox(
					widget.NewButtonWithIcon("", theme.Icon(theme.IconNameSettings), func() {
						if taskList != nil {
							app.RenderMutateTaskView(&task, taskList)
						} else {
							app.RenderMutateTaskView(&task, task.TaskList)
						}
					}),
					widget.NewButtonWithIcon("", theme.Icon(theme.IconNameDelete), func() {
						res := app.DB().Delete(&task)
						if res.Error != nil {
							panic(fmt.Sprintf("Error deleting task %d: %v", task.ID, res.Error))
						}
						onDelete()
					}),
				),
				widget.NewLabel(task.Label),
			))
		},
	)
}

var _ View = (*ListOfTasksView)(nil)

type ListOfTasksView struct {
	*baseView
	title    string
	taskList *TaskList
	opts     []ModelQueryOpt
}

func NewListOfTasksView(app *TaskApp, title string, taskList *TaskList, opts ...ModelQueryOpt) *ListOfTasksView {
	v := ListOfTasksView{
		baseView: newBaseView("Task List View", app),
		title:    title,
		taskList: taskList,
		opts:     append(opts, WithPreload("TaskList")),
	}
	return &v
}

func (v *ListOfTasksView) Title() string {
	return v.title
}

func (v *ListOfTasksView) Foreground() fyne.CanvasObject {
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

	taskCount, err := CountModel[Task](ctx, v.app.DB(), v.opts...)
	if err != nil {
		if errors.Is(err, context.Canceled) {
			return nil
		}
		panic(fmt.Sprintf("Error counting tasks: %v", err))
	}

	ftr := container.NewBorder(
		nil,
		nil,
		canvas.NewText(fmt.Sprintf("Total tasks: %d", taskCount), color.Black),
		widget.NewButtonWithIcon("New task", theme.Icon(theme.IconNameContentAdd), func() {
			v.app.RenderMutateTaskView(nil, v.taskList)
		}),
	)

	tasks, err := FindModel[Task](ctx, v.app.DB(), v.opts...)
	if err != nil {
		if errors.Is(err, context.Canceled) {
			return nil
		}
		panic(fmt.Sprintf("Error finding tasks: %v", err))
	}

	return container.NewBorder(
		nil,
		ftr,
		nil,
		nil,
		buildListOfTasksList(
			v.app,
			v.taskList,
			tasks,
			func() { v.app.RenderListOfTasksView(v.Name(), v.taskList, v.opts...) },
		),
	)
}

func (v *ListOfTasksView) Background() {
	v.mu.Lock()
	defer v.mu.Unlock()
	v.background()
}
