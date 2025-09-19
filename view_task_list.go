package main

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"gorm.io/gorm"
)

var _ View = (*TaskListView)(nil)

type TaskListView struct {
	*baseView
	taskList TaskList
	onDelete func()
}

func NewTaskListView(ta *TaskApp, taskList TaskList, onDelete func()) *TaskListView {
	v := TaskListView{
		baseView: newBaseView("Task List View", ta),
		taskList: taskList,
		onDelete: onDelete,
	}
	return &v
}

func (v *TaskListView) Title() string {
	return v.taskList.Label
}

func (v *TaskListView) Foreground() fyne.CanvasObject {
	v.mu.Lock()
	defer v.mu.Unlock()
	if !v.foreground() {
		return nil
	}

	body := container.NewVBox(
		FormLabel("Description:"),
		container.NewHScroll(
			widget.NewRichTextFromMarkdown(v.taskList.Description),
		),
	)

	ftr := container.NewHBox(
		layout.NewSpacer(),
		widget.NewButtonWithIcon("Delete", theme.DeleteIcon(), func() {
			res := v.app.DB().Delete(v.taskList)
			if res.Error != nil {
				panic(fmt.Sprintf("Error deleting task list %d: %v", v.taskList.ID, res.Error))
			}
			v.app.RenderTaskListsView()
		}),
		widget.NewButtonWithIcon("Edit", IconEdit, func() {
			v.app.RenderMutateTaskListView(&v.taskList)
		}),
		widget.NewButtonWithIcon("New task", theme.ContentAddIcon(), func() {
			v.app.RenderMutateTaskView(nil, &v.taskList, func() {
				v.app.RenderListOfTasksView(v.Name(), &v.taskList, func(db *gorm.DB) *gorm.DB {
					return db.Where("task_list_id = ?", v.taskList.ID)
				})
			})
		}),
	)

	return container.NewBorder(
		nil,
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
