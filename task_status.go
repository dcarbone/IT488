package main

import (
	"image"
	"strings"

	"fyne.io/fyne/v2"
)

const (
	TaskStatusTodo uint = 0
	TaskStatusDone uint = 10
	TaskStatusSkip uint = 20

	TaskStatusTitleTodo = "Todo"
	TaskStatusTitleDone = "Done"
	TaskStatusTitleSkip = "Skip"
)

var (
	TaskStatusTitles = []string{
		TaskStatusTitleTodo,
		TaskStatusTitleDone,
		TaskStatusTitleSkip,
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

func TaskStatusNumber(status string) uint {
	switch strings.ToLower(status) {
	case strings.ToLower(TaskStatusTitleSkip):
		return TaskStatusSkip
	case strings.ToLower(TaskStatusTitleDone):
		return TaskStatusDone

	default:
		return TaskStatusTodo
	}
}

func TaskStatusTitle(status uint) string {
	switch status {
	case TaskStatusSkip:
		return TaskStatusTitleSkip
	case TaskStatusDone:
		return TaskStatusTitleDone

	default:
		return TaskStatusTitleTodo
	}
}

func TaskStatusImage(status uint) image.Image {
	switch status {
	case TaskStatusDone:
		return AssetImageStatusDone
	case TaskStatusSkip:
		return AssetImageStatusSkip

	default:
		return AssetImageStatusTodo
	}
}

func TaskStatusResource(status uint) *fyne.StaticResource {
	switch status {
	case TaskStatusDone:
		return TaskStatusIconResourceDone
	case TaskStatusSkip:
		return TaskStatusIconResourceSkip

	default:
		return TaskStatusIconResourceTodo
	}
}
