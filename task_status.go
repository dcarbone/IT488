package main

import (
	"fmt"
	"image"
	"slices"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
	"gorm.io/gorm"
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

func newTaskStatusSwitcherButton(db *gorm.DB, task *Task) *widget.Button {
	var statusButton *widget.Button
	var statusIdx = slices.Index(TaskStatusTitles, TaskStatusTitle(task.Status))
	statusButton = widget.NewButtonWithIcon(
		"",
		TaskStatusResource(task.Status),
		func() {
			statusIdx++
			if statusIdx == len(TaskStatusTitles) {
				statusIdx = 0
			}
			task.Status = TaskStatusNumber(TaskStatusTitles[statusIdx])
			res := db.Model(&task).Update("Status", task.Status)
			if res.Error != nil {
				panic(fmt.Sprintf("error updating task status: %v", res.Error))
			}
			statusButton.SetIcon(TaskStatusResource(task.Status))
		},
	)
	statusButton.Importance = widget.LowImportance

	return statusButton
}
