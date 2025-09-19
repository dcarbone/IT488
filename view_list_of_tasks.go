package main

import (
	"context"
	"errors"
	"fmt"
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

func buildListOfTasksList(app *TaskApp, taskList *TaskList, tasks []Task, onDelete func()) fyne.CanvasObject {
	list := widget.NewList(
		func() int {
			return len(tasks)
		},
		func() fyne.CanvasObject {
			return container.NewStack(widget.NewLabel("Loading..."))
		},
		func(id widget.ListItemID, object fyne.CanvasObject) {
			task := &tasks[id]

			content := object.(*fyne.Container)

			content.RemoveAll()
			content.Add(container.NewBorder(
				nil,
				nil,
				container.NewHBox(
					newTaskStatusSwitcherButton(app.DB(), task),
					newTaskPrioritySwitcherButton(app.DB(), task),
				),
				container.NewHBox(
					widget.NewButtonWithIcon("", IconEdit, func() {
						if taskList != nil {
							app.RenderMutateTaskView(task, taskList, onDelete)
						} else {
							app.RenderMutateTaskView(task, task.TaskList, onDelete)
						}
					}),
				),
				widget.NewLabel(task.Label),
			))
		},
	)

	list.OnSelected = func(id widget.ListItemID) {
		app.RenderTaskView(tasks[id], onDelete)
	}

	return list
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
		opts:     opts,
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

	ftr := container.NewHBox(
		canvas.NewText(fmt.Sprintf("Total tasks: %d", taskCount), color.Black),
		layout.NewSpacer(),
	)
	if v.taskList != nil {
		ftr.Add(widget.NewButtonWithIcon("Edit", IconEdit, func() {
			v.app.RenderMutateTaskListView(v.taskList)
		}))
	}
	ftr.Add(widget.NewButtonWithIcon("New task", theme.ContentAddIcon(), func() {
		v.app.RenderMutateTaskView(nil, v.taskList, func() {
			v.app.RenderListOfTasksView(v.Name(), v.taskList, v.opts...)
		})
	}))

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
			func() { v.app.RenderListOfTasksView(v.Name(), v.taskList, append(v.opts, WithPreload("TaskList"))...) },
		),
	)
}

func (v *ListOfTasksView) Background() {
	v.mu.Lock()
	defer v.mu.Unlock()
	v.background()
}
