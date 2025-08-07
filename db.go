package main

import (
	"gorm.io/gorm"
)

var (
	dbFile string
)

type TaskStatus string

const (
	TaskStatusTodo       TaskStatus = "todo"
	TaskStatusInProgress TaskStatus = "in progress"
	TaskStatusCompleted  TaskStatus = "completed"
	TaskStatusSkip       TaskStatus = "skip"
)

type Task struct {
	gorm.
	ID          int
	Label       string
	Description string
	Status      TaskStatus
}
