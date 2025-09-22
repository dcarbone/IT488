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
	TaskPriorityNeutral = "neutral"
	TaskPriorityHigh    = "high"
	TaskPriorityHighest = "highest"
)

var (
	TaskPriorities = []string{
		strings.ToTitle(TaskPriorityLowest),
		strings.ToTitle(TaskPriorityLow),
		strings.ToTitle(TaskPriorityNeutral),
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
	TaskPriorityIconResourceLNeutral = EncodeImageToResource(
		"task_priority_neutral",
		GetConstrainedImage(TaskPriorityImage(TaskPriorityNeutral), 50),
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
	case TaskPriorityHigh:
		return 20
	case TaskPriorityHighest:
		return 30

	default:
		return 15
	}
}

func TaskPriorityName(priority uint) string {
	switch priority {
	case 0:
		return TaskPriorityLowest
	case 10:
		return TaskPriorityLow
	case 20:
		return TaskPriorityHigh
	case 30:
		return TaskPriorityHighest

	default:
		return TaskPriorityNeutral
	}
}

func TaskPriorityImage(priority string) image.Image {
	switch strings.ToLower(priority) {
	case TaskPriorityLowest:
		return AssetImagePriorityLowest
	case TaskPriorityLow:
		return AssetImagePriorityLow
	case TaskPriorityHigh:
		return AssetImagePriorityHigh
	case TaskPriorityHighest:
		return AssetImagePriorityHighest

	default:
		return AssetImagePriorityNeutral
	}
}

func TaskPriorityResource(priority string) *fyne.StaticResource {
	switch strings.ToLower(priority) {
	case TaskPriorityLowest:
		return TaskPriorityIconResourceLowest
	case TaskPriorityLow:
		return TaskPriorityIconResourceLow
	case TaskPriorityHigh:
		return TaskPriorityIconResourceHigh
	case TaskPriorityHighest:
		return TaskPriorityIconResourceHighest

	default:
		return TaskPriorityIconResourceLNeutral
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
