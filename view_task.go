package main

import (
	"fyne.io/fyne/v2"
)

var _ View = (*TaskView)(nil)

type TaskView struct {
	*baseView
	task Task
}

func NewTaskView(ta *TaskApp, task Task) *TaskView {
	v := TaskView{
		baseView: newBaseView("Task View", ta),
		task:     task,
	}
	return &v
}

func (v *TaskView) Title() string {
	return v.task.Label
}

func (v *TaskView) Foreground() fyne.CanvasObject {
	v.mu.Lock()
	defer v.mu.Unlock()
	if !v.foreground() {
		return nil
	}
	// todo: finish me!
	return nil
}

func (v *TaskView) Background() {
	v.mu.Lock()
	defer v.mu.Unlock()
	v.background()
}
