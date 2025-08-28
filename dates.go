package main

import (
	"time"
)

const (
	TimestampDisplayFormat = "Jan _2 3:04:05PM"
)

func FormatDateTime(tm time.Time) string {
	return tm.Format(TimestampDisplayFormat)
}
