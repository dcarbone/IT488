package main

import (
	"fmt"
	"image/color"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"gorm.io/gorm"
)

var _ View = (*MutateTaskListView)(nil)

type MutateTaskListView struct {
	*baseView
	taskList *TaskList
}

func NewMutateTaskListView(app *TaskApp, taskList *TaskList) *MutateTaskListView {
	v := MutateTaskListView{
		baseView: newBaseView("Create Task List", app),
		taskList: taskList,
	}
	return &v
}

func (v *MutateTaskListView) Title() string {
	if v.taskList == nil {
		return "Create task list"
	}
	return fmt.Sprintf("Edit task list %s", v.taskList.Label)
}

func (v *MutateTaskListView) Foreground() fyne.CanvasObject {
	v.mu.Lock()
	defer v.mu.Unlock()
	if v.foreground() {
		return v.render(nil)
	}
	return nil
}

func (v *MutateTaskListView) Background() {
	v.mu.Lock()
	defer v.mu.Unlock()
	v.background()
}

func (v *MutateTaskListView) render(err error) fyne.CanvasObject {
	content := container.NewVBox()

	if err != nil {
		content.Add(canvas.NewText("Error:", ColorRed))
		content.Add(canvas.NewText(err.Error(), ColorRed))
	}

	content.Add(canvas.NewText("Name:", color.Black))

	labelInput := widget.NewEntry()
	labelInput.PlaceHolder = "Enter task list name."
	labelInput.OnChanged = func(s string) {
		if len(s) > 50 {
			labelInput.SetText(s[:50])
		}
	}
	if v.taskList != nil {
		labelInput.SetText(v.taskList.Label)
	}

	content.Add(labelInput)

	content.Add(canvas.NewText("Description:", color.Black))

	descInput := widget.NewMultiLineEntry()
	descInput.PlaceHolder = "Enter Markdown formatted text."
	descInput.OnChanged = func(s string) {
		if len(s) > 500 {
			descInput.SetText(s[:500])
		}
	}
	if v.taskList != nil {
		descInput.SetText(v.taskList.Description)
	}

	content.Add(descInput)

	ftr := container.NewHBox(layout.NewSpacer())

	if v.taskList != nil {
		ftr.Add(widget.NewButtonWithIcon("Delete", theme.DeleteIcon(), func() {
			res := v.app.DB().Delete(v.taskList)
			if res.Error != nil {
				panic(fmt.Sprintf("Error deleting task list %d: %v", v.taskList.ID, res.Error))
			}
			v.app.RenderTaskListsView()
		}))
	}

	ftr.Add(widget.NewButtonWithIcon("Cancel", theme.CancelIcon(), func() {
		v.app.RenderPreviousView()
	}))
	ftr.Add(widget.NewButtonWithIcon(
		"Save",
		theme.DocumentSaveIcon(),
		func() {
			var res *gorm.DB
			if v.taskList != nil {
				v.taskList.Label = labelInput.Text
				v.taskList.Description = descInput.Text
				res = v.app.DB().Updates(v.taskList)
			} else {
				v.taskList = &TaskList{
					Label:       labelInput.Text,
					Date:        time.Now(),
					Description: descInput.Text,
				}
				res = v.app.DB().Create(v.taskList)
			}
			if res.Error != nil {
				v.render(err)
				return
			}

			v.app.RenderListOfTasksView(v.taskList.Label, v.taskList, func(db *gorm.DB) *gorm.DB {
				return db.Where("task_list_id = ?", v.taskList.ID)
			})
		},
	))

	return container.NewBorder(
		nil,
		ftr,
		nil,
		nil,
		content,
	)
}
