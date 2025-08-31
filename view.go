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
	Title() string
	State() ViewState
	Foreground() fyne.CanvasObject
	Background()
}

type baseView struct {
	mu          sync.Mutex
	log         *slog.Logger
	state       ViewState
	name        string
	app         *TaskApp
	deactivated chan struct{}
	children    []View
}

func newBaseView(name string, app *TaskApp) *baseView {
	v := baseView{
		log:      log.With("view", name),
		state:    ViewStateBackground,
		name:     name,
		app:      app,
		children: make([]View, 0),
	}
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
