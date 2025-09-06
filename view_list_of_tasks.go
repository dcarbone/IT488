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
	return widget.NewList(
		func() int {
			return len(tasks)
		},
		func() fyne.CanvasObject {
			return container.NewStack(widget.NewLabel("Loading..."))
		},
		func(id widget.ListItemID, object fyne.CanvasObject) {
			var statusPickerModal *widget.PopUp

			task := tasks[id]

			content := object.(*fyne.Container)

			statusIcon := NewTappableIcon(TaskStatusResource(task.Status), func(ev *fyne.PointEvent) {
				statusPickerModal.Show()
			})

			var pickedStatus string
			statusPickerList := widget.NewList(
				func() int {
					return len(TaskStatuses)
				},
				func() fyne.CanvasObject {
					return container.NewStack(widget.NewLabel("Loading..."))
				},
				func(id widget.ListItemID, object fyne.CanvasObject) {
					content := object.(*fyne.Container)
					content.RemoveAll()
					content.Add(container.NewHBox(
						widget.NewIcon(TaskStatusResource(TaskStatuses[id])),
						widget.NewLabel(TaskStatuses[id]),
					))
				},
			)
			statusPickerList.OnSelected = func(id widget.ListItemID) {
				pickedStatus = TaskStatuses[id]
			}
			statusPickerList.Select(func() widget.ListItemID {
				for i, stat := range TaskStatuses {
					if stat == task.Status {
						return i
					}
				}
				return 0
			}())

			statusPickerContainer := container.NewBorder(
				nil,
				container.NewHBox(
					layout.NewSpacer(),
					widget.NewButton("OK", func() {
						statusPickerModal.Hide()
						task.Status = pickedStatus
						res := app.DB().Model(&task).Update("Status", task.Status)
						if res.Error != nil {
							panic(fmt.Sprintf("error updating task status: %v", res.Error))
						}
						statusIcon.SetResource(TaskStatusResource(pickedStatus))
					}),
				),
				nil,
				nil,
				statusPickerList,
			)

			statusPickerModal = widget.NewModalPopUp(
				statusPickerContainer,
				app.window.Canvas(),
			)

			content.RemoveAll()
			content.Add(container.NewBorder(
				nil,
				nil,
				statusIcon,
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
				container.NewHBox(
					task.PriorityIcon(),
					widget.NewLabel(task.Label),
				),
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
