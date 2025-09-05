package main

import (
	"image"
	"strings"
)

const (
	TaskStatusTodo       = "todo"
	TaskStatusInProgress = "in progress"
	TaskStatusCompleted  = "completed"
	TaskStatusSkip       = "skip"
)

var (
	TaskStatuses = []string{
		strings.ToTitle(TaskStatusTodo),
		strings.ToTitle(TaskStatusInProgress),
		strings.ToTitle(TaskStatusCompleted),
		strings.ToTitle(TaskStatusSkip),
	}
)

func TaskStatusImage(status string) image.Image {
	switch strings.ToLower(status) {
	case TaskStatusInProgress:
		return AssetImageStatusInProgress
	case TaskStatusCompleted:
		return AssetImageStatusDone
	case TaskStatusSkip:
		return AssetImageStatusSkip

	default:
		return AssetImageStatusTodo
	}
}
