package main

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

func buildListOfTasksList(app *TaskApp, tasks []Task, onDelete func()) fyne.CanvasObject {
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

			content.RemoveAll()
			content.Add(container.NewBorder(
				nil,
				nil,
				nil,
				container.NewHBox(
					widget.NewButtonWithIcon("", theme.Icon(theme.IconNameSettings), func() {
						app.RenderMutateTaskView(&task, nil)
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
