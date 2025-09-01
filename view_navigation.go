package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
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

func (v *NavigationView) Title() string {
	return "Navigation"
}

func (v *NavigationView) Foreground() fyne.CanvasObject {
	v.mu.Lock()
	defer v.mu.Unlock()
	if !v.foreground() {
		return nil
	}

	return container.NewBorder(
		container.NewBorder(
			nil,
			nil,
			nil,
			widget.NewButtonWithIcon("", theme.Icon(theme.IconNameCancel), func() {
				v.app.RenderPreviousView()
			}),
		),
		nil,
		nil,
		nil,
		container.NewVScroll(
			container.NewVBox(
				widget.NewButton("Today's Tasks", func() {
					v.app.RenderListOfTasksView("Today's List", nil, todaysTasksModelQueryOpt())
				}),
				widget.NewButton("Task Lists", func() {
					v.app.RenderTaskListsView()
				}),
				widget.NewButton("All Tasks", func() {
					v.app.RenderListOfTasksView("All tasks", nil)
				}),
			),
		),
	)
}

func (v *NavigationView) Background() {
	v.mu.Lock()
	defer v.mu.Unlock()
	v.background()
}
