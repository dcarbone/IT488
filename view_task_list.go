package main

import (
	"context"
	"errors"
	"fmt"
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
)

var _ View = (*TaskListView)(nil)

type TaskListView struct {
	*baseView
	title string
	opts  []ModelQueryOpt
}

func NewTaskListView(app *TaskApp, title string, opts ...ModelQueryOpt) *TaskListView {
	v := TaskListView{
		baseView: newBaseView("Task List View", app),
		title:    title,
		opts:     opts,
	}
	return &v
}

func (v *TaskListView) Title() string {
	return v.title
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
	//HeaderCanvas(v.taskList.Label),
	//widget.NewButtonWithIcon("", theme.Icon(theme.IconNameContentAdd), func() {
	//	v.app.RenderMutateTaskView(nil, &v.taskList)
	//}),
	)

	taskCount, err := CountModel[Task](ctx, v.app.DB(), v.opts...)
	if err != nil {
		if errors.Is(err, context.Canceled) {
			return nil
		}
		panic(fmt.Sprintf("Error counting tasks: %v", err))
	}

	ftr := canvas.NewText(fmt.Sprintf("Total tasks: %d", taskCount), color.Black)

	tasks, err := FindModel[Task](ctx, v.app.DB(), v.opts...)
	if err != nil {
		if errors.Is(err, context.Canceled) {
			return nil
		}
		panic(fmt.Sprintf("Error finding tasks: %v", err))
	}

	return container.NewBorder(
		hdr,
		ftr,
		nil,
		nil,
		buildListOfTasksList(v.app, tasks, func() { v.app.RenderTaskListView(v.Name(), v.opts...) }),
	)
}

func (v *TaskListView) Background() {
	v.mu.Lock()
	defer v.mu.Unlock()
	v.background()
}
