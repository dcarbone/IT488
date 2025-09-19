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
	TaskPriorityLowest  = "lowest"
	TaskPriorityLow     = "low"
	TaskPriorityHigh    = "high"
	TaskPriorityHighest = "highest"
)

var (
	TaskPriorities = []string{
		strings.ToTitle(TaskPriorityLowest),
		strings.ToTitle(TaskPriorityLow),
		strings.ToTitle(TaskPriorityHigh),
		strings.ToTitle(TaskPriorityHighest),
	}

	TaskPriorityIconResourceLowest = EncodeImageToResource(
		"task_priority_lowest",
		GetConstrainedImage(TaskPriorityImage(TaskPriorityLowest), 50),
	)
	TaskPriorityIconResourceLow = EncodeImageToResource(
		"task_priority_low",
		GetConstrainedImage(TaskPriorityImage(TaskPriorityLow), 50),
	)
	TaskPriorityIconResourceHigh = EncodeImageToResource(
		"task_priority_high",
		GetConstrainedImage(TaskPriorityImage(TaskPriorityHigh), 50),
	)
	TaskPriorityIconResourceHighest = EncodeImageToResource(
		"task_priority_highest",
		GetConstrainedImage(TaskPriorityImage(TaskPriorityHighest), 50),
	)
)

func TaskPriorityNumber(priority string) uint {
	switch strings.ToLower(priority) {
	case TaskPriorityLowest:
		return 0
	case TaskPriorityLow:
		return 10
	case TaskPriorityHighest:
		return 30

	default:
		return 20
	}
}

func TaskPriorityName(priority uint) string {
	switch priority {
	case 0:
		return TaskPriorityLowest
	case 10:
		return TaskPriorityLow
	case 30:
		return TaskPriorityHighest

	default:
		return TaskPriorityHigh
	}
}

func TaskPriorityImage(priority string) image.Image {
	switch strings.ToLower(priority) {
	case TaskPriorityLowest:
		return AssetImagePriorityLowest
	case TaskPriorityLow:
		return AssetImagePriorityLow
	case TaskPriorityHighest:
		return AssetImagePriorityHighest

	default:
		return AssetImagePriorityHigh
	}
}

func TaskPriorityResource(priority string) *fyne.StaticResource {
	switch strings.ToLower(priority) {
	case TaskPriorityLowest:
		return TaskPriorityIconResourceLowest
	case TaskPriorityLow:
		return TaskPriorityIconResourceLow
	case TaskPriorityHighest:
		return TaskPriorityIconResourceHighest

	default:
		return TaskPriorityIconResourceHigh
	}
}

func newTaskPrioritySwitcherButton(db *gorm.DB, task *Task) *widget.Button {
	var priorityButton *widget.Button
	priorityIdx := slices.Index(TaskPriorities, strings.ToTitle(TaskPriorityName(task.UserPriority)))
	priorityButton = widget.NewButtonWithIcon(
		"",
		TaskPriorityResource(TaskPriorityName(task.UserPriority)),
		func() {
			priorityIdx++
			if priorityIdx == len(TaskPriorities) {
				priorityIdx = 0
			}
			task.UserPriority = TaskPriorityNumber(TaskPriorities[priorityIdx])
			res := db.Model(&task).Update("UserPriority", task.UserPriority)
			if res.Error != nil {
				panic(fmt.Sprintf("error updating task user priority: %v", res.Error))
			}
			priorityButton.SetIcon(TaskPriorityResource(TaskPriorities[priorityIdx]))
		},
	)
	priorityButton.Importance = widget.LowImportance

	return priorityButton
}
