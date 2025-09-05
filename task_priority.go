package main

import (
	"image"
	"strings"
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
