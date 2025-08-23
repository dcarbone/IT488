package main

import (
	"context"
	"errors"
	"fmt"
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

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
			v.app.RenderMutateTaskModal(nil, &v.taskList)
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

	var listOfTasks *widget.List
	listOfTasks = widget.NewList(
		func() int {
			return int(taskCount)
		},
		func() fyne.CanvasObject {
			return container.NewStack(widget.NewLabel("Loading..."))
		},
		func(id widget.ListItemID, object fyne.CanvasObject) {
			task := tasks[id]

			content := object.(*fyne.Container)

			content.RemoveAll()
			content.Add(container.NewBorder(
				nil,
				canvas.NewText(fmt.Sprintf("Created: %s", task.CreatedAt), color.Black),

				widget.NewButtonWithIcon("", theme.Icon(theme.IconNameSettings), func() {
					v.app.RenderMutateTaskModal(&task, &v.taskList)
				}),

				widget.NewButtonWithIcon("", theme.Icon(theme.IconNameDelete), func() {
					res := v.app.DB().Delete(&task)
					if res.Error != nil {
						panic(fmt.Sprintf("Error deleting task %d: %v", task.ID, res.Error))
					}
					v.app.RenderTaskListView(v.taskList)
				}),

				canvas.NewText(task.Label, color.Black),
			))
		},
	)

	return container.NewBorder(
		hdr,
		ftr,
		nil,
		nil,
		listOfTasks,
	)
}

func (v *TaskListView) Background() {
	v.mu.Lock()
	defer v.mu.Unlock()
	v.background()
}
