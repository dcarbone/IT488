package main

import (
	"image"
	"strings"
)

const (
	TaskStatusTodo = "todo"
	TaskStatusDone = "done"
	TaskStatusSkip = "skip"
)

var (
	TaskStatuses = []string{
		strings.ToTitle(TaskStatusTodo),
		strings.ToTitle(TaskStatusDone),
		strings.ToTitle(TaskStatusSkip),
	}
)

func TaskStatusImage(status string) image.Image {
	switch strings.ToLower(status) {
	case TaskStatusDone:
		return AssetImageStatusDone
	case TaskStatusSkip:
		return AssetImageStatusSkip

	default:
		return AssetImageStatusTodo
	}
}
