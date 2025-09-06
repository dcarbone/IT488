package main

import (
	"image"
	"strings"

	"fyne.io/fyne/v2"
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

	TaskStatusIconResourceTodo = EncodeImageToResource(
		"task_status_todo",
		GetConstrainedImage(TaskStatusImage(TaskStatusTodo), 50),
	)
	TaskStatusIconResourceDone = EncodeImageToResource(
		"task_status_done",
		GetConstrainedImage(TaskStatusImage(TaskStatusDone), 50),
	)
	TaskStatusIconResourceSkip = EncodeImageToResource(
		"task_status_skip",
		GetConstrainedImage(TaskStatusImage(TaskStatusSkip), 50),
	)
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

func TaskStatusResource(status string) *fyne.StaticResource {
	switch strings.ToLower(status) {
	case TaskStatusDone:
		return TaskStatusIconResourceDone
	case TaskStatusSkip:
		return TaskStatusIconResourceSkip

	default:
		return TaskStatusIconResourceTodo
	}
}
