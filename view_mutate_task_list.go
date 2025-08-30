package main

import (
	"image/color"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"gorm.io/gorm"
)

var _ View = (*CreateTaskListView)(nil)

type CreateTaskListView struct {
	*baseView
	taskList *TaskList
}

func NewMutateTaskListView(app *TaskApp, taskList *TaskList) *CreateTaskListView {
	v := CreateTaskListView{
		baseView: newBaseView("Create Task List", app),
		taskList: taskList,
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
			v.app.RenderTaskListView(tl.Label, func(db *gorm.DB) *gorm.DB {
				return db.Where("task_list_id = ?", tl.ID)
			})
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
