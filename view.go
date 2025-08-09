package main

import (
	"log/slog"
	"sync"

	"fyne.io/fyne/v2"
)

type ViewState int

const (
	ViewStateForeground ViewState = iota
	ViewStateBackground
)

type View interface {
	Name() string
	State() ViewState
	Content() fyne.CanvasObject
	Do(act Action)
	Foreground()
	Background()
}

type baseView struct {
	mu          sync.Mutex
	log         *slog.Logger
	state       ViewState
	name        string
	acts        chan Action
	deactivated chan struct{}
	children    []View
}

func newBaseView(name string, actFn ActionHandler) *baseView {
	v := baseView{
		log:      log.With("view", name),
		state:    ViewStateBackground,
		name:     name,
		acts:     make(chan Action, 100),
		children: make([]View, 0),
	}
	go func() {
		defer actFn(nil, true)
		for act := range v.acts {
			actFn(act, false)
		}
	}()
	return &v
}

func (v *baseView) Name() string {
	return v.name
}

func (v *baseView) State() ViewState {
	v.mu.Lock()
	defer v.mu.Unlock()
	return v.state
}

func (v *baseView) Do(act Action) {
	v.mu.Lock()
	defer v.mu.Unlock()
	v.acts <- act
	for i := range v.children {
		v.children[i].Do(v)
	}
}

func (v *baseView) foreground() bool {
	if v.state == ViewStateForeground {
		v.log.Debug("View already in foreground")
		return false
	}
	v.log.Debug("Setting view state to foreground")
	v.state = ViewStateForeground
	v.deactivated = make(chan struct{})
	return true
}

func (v *baseView) background() bool {
	if v.state == ViewStateBackground {
		v.log.Debug("View already in background")
		return false
	}
	v.log.Debug("Setting view state to background")
	v.state = ViewStateBackground
	for _, child := range v.children {
		v.log.Debug("Backgrounding view child...", "child", child.Name())
		child.Background()
	}
	close(v.deactivated)
	v.children = make([]View, 0)
	return true
}

var _ View = (*TaskListsView)(nil)

type TaskListsView struct {
	*baseView
	content fyne.CanvasObject
}

func NewTaskListView() *TaskListsView {
	v := TaskListsView{}
	v.baseView = newBaseView("Task List", v.do)
	return &v
}

func (v *TaskListsView) Content() fyne.CanvasObject {
	//TODO implement me
	panic("implement me")
}

func (v *TaskListsView) do(act Action, closed bool) {

}

func (v *TaskListsView) Foreground() {
	//TODO implement me
	panic("implement me")
}

func (v *TaskListsView) Background() {
	//TODO implement me
	panic("implement me")
}
