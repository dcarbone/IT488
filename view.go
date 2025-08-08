package main

import (
	"sync"
)

type ViewState int

const (
	ViewStateForeground ViewState = iota
	ViewStateBackground
)

type View interface {
	Foreground()
	Background()
}

var _ View = (*TaskListsView)(nil)

type TaskListsView struct {
	mu    sync.Mutex
	state ViewState
}

func (t *TaskListsView) Foreground() {
	//TODO implement me
	panic("implement me")
}

func (t *TaskListsView) Background() {
	//TODO implement me
	panic("implement me")
}
