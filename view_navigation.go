package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"gorm.io/gorm"
)

var _ View = (*NavigationView)(nil)

type NavigationView struct {
	*baseView
}

func NewNavigationView(ta *TaskApp) *NavigationView {
	v := NavigationView{
		baseView: newBaseView("Navigation", ta),
	}
	return &v
}

func (v *NavigationView) Title() []fyne.CanvasObject {
	return []fyne.CanvasObject{
		HeaderCanvas("Navigation"),
		layout.NewSpacer(),
		widget.NewButtonWithIcon("", theme.CancelIcon(), func() {
			v.app.RenderPreviousView()
		}),
	}
}

func (v *NavigationView) Foreground() fyne.CanvasObject {
	v.mu.Lock()
	defer v.mu.Unlock()
	if !v.foreground() {
		return nil
	}

	return container.NewVScroll(
		container.NewVBox(
			widget.NewSeparator(),
			widget.NewSeparator(),
			widget.NewSeparator(),

			widget.NewButtonWithIcon("Home", theme.HomeIcon(), func() {
				v.app.RenderHomeView()
			}),

			widget.NewSeparator(),
			widget.NewSeparator(),
			widget.NewSeparator(),

			widget.NewButton("Lists", func() {
				v.app.RenderTaskListsView()
			}),

			widget.NewSeparator(),
			widget.NewSeparator(),
			widget.NewSeparator(),

			widget.NewButton("Today's Tasks", func() {
				v.app.RenderListOfTasksView("Today's Tasks", nil, todaysTasksModelQueryOpt())
			}),
			widget.NewButton("Todo Tasks", func() {
				v.app.RenderListOfTasksView("Todo Tasks", nil, func(db *gorm.DB) *gorm.DB {
					return WithSort("due_date asc")(WithSort("id asc")(WithPreload("TaskList")(db))).
						Where("Status = ?", TaskStatusTodo)
				})
			}),
			widget.NewButton("Done Tasks", func() {
				v.app.RenderListOfTasksView("Done Tasks", nil, func(db *gorm.DB) *gorm.DB {
					return WithSort("due_date asc")(WithSort("id asc")(WithPreload("TaskList")(db))).
						Where("Status in ?", []uint{TaskStatusSkip, TaskStatusDone})
				})
			}),

			widget.NewSeparator(),
			widget.NewSeparator(),
			widget.NewSeparator(),
		),
	)
}

func (v *NavigationView) Background() {
	v.mu.Lock()
	defer v.mu.Unlock()
	v.background()
}
