package main

import (
	"context"
	"fmt"
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"gorm.io/gorm"
)

var _ View = (*TaskListsView)(nil)

type TaskListsView struct {
	*baseView
}

func NewTaskListsView(ta *TaskApp) *TaskListsView {
	v := TaskListsView{
		baseView: newBaseView("Task Lists", ta),
	}
	return &v
}

func (v *TaskListsView) Title() string {
	return "Task lists"
}

func (v *TaskListsView) Foreground() fyne.CanvasObject {
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

	listCount, err := CountModel[TaskList](ctx, v.app.DB())
	if err != nil {
		panic(fmt.Sprintf("Error counting task lists: %v", err))
	}

	taskLists, err := FindModel[TaskList](ctx, v.app.DB())
	if err != nil {
		panic(fmt.Sprintf("Error fetching tasks: %v", err))
	}

	listView := widget.NewList(
		func() int {
			return int(listCount)
		},
		func() fyne.CanvasObject {
			return container.NewStack(widget.NewLabel("Loading..."))
		},
		func(id widget.ListItemID, object fyne.CanvasObject) {
			taskList := taskLists[id]

			content := object.(*fyne.Container)

			content.RemoveAll()
			content.Add(container.NewBorder(
				nil,
				nil,
				nil,
				container.NewHBox(
					widget.NewButtonWithIcon("", theme.ListIcon(), func() {
						v.app.RenderListOfTasksView(taskList.Label, &taskList, func(db *gorm.DB) *gorm.DB {
							return db.Where("task_list_id = ?", taskList.ID)
						})
					}),
					widget.NewButtonWithIcon("", IconEdit, func() {
						v.app.RenderMutateTaskListView(&taskList)
					}),
				),
				widget.NewLabel(taskList.Label),
			))
		},
	)

	ftr := container.NewBorder(
		nil,
		nil,
		canvas.NewText(fmt.Sprintf("Total lists: %d", listCount), color.Black),
		widget.NewButtonWithIcon("New list", theme.ContentAddIcon(), func() {
			v.app.RenderMutateTaskListView(nil)
		}),
	)

	return container.NewBorder(
		nil,
		ftr,
		nil,
		nil,
		listView,
	)
}

func (v *TaskListsView) Background() {
	v.mu.Lock()
	defer v.mu.Unlock()
	v.background()
}
