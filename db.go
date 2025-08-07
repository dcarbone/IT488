package main

import (
	"context"
	"database/sql/driver"
	"fmt"
	"time"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
	glogger "gorm.io/gorm/logger"
)

var (
	dbFile string
)

type gormLogger struct {
	logMode glogger.LogLevel
}

func (gl gormLogger) LogMode(level glogger.LogLevel) glogger.Interface {
	gl.logMode = level
	return gl
}

func (gl gormLogger) Info(_ context.Context, f string, v ...any) {
	if gl.logMode <= glogger.Info {
		log.Info(fmt.Sprintf(f, v...))
	}
}

func (gl gormLogger) Warn(_ context.Context, f string, v ...any) {
	if gl.logMode <= glogger.Warn {
		log.Warn(fmt.Sprintf(f, v...))
	}
}

func (gl gormLogger) Error(_ context.Context, f string, v ...any) {
	if gl.logMode <= glogger.Error {
		log.Error(fmt.Sprintf(f, v...))
	}
}

func (gl gormLogger) Trace(_ context.Context, begin time.Time, fc func() (sql string, rowsAffected int64), err error) {
	sql, af := fc()
	log.Debug("Query trace", "begin", begin, "sql", sql, "rows_affected", af, "err", err)
}

func openDB() (*gorm.DB, error) {
	log.Debug("Opening sqlite db...", "db", dbFile)

	conf := &gorm.Config{
		Logger: gormLogger{},
	}
	db, err := gorm.Open(sqlite.Open(fmt.Sprintf("%s:?_pragma=foreign_keys(1)", dbFile)), conf)
	if err != nil {
		return nil, err
	}

	log.Debug("Applying migrations...")

	if err = db.AutoMigrate(&TaskList{}, &Task{}); err != nil {
		return nil, fmt.Errorf("error applying migrations: %w", err)
	}

	return db, nil
}

type TaskStatus string

const (
	TaskStatusTodo       TaskStatus = "todo"
	TaskStatusInProgress TaskStatus = "in progress"
	TaskStatusCompleted  TaskStatus = "completed"
	TaskStatusSkip       TaskStatus = "skip"
)

func (ct *TaskStatus) Scan(value any) error {
	*ct = TaskStatus(value.([]byte))
	return nil
}

func (ct TaskStatus) Value() (driver.Value, error) {
	return string(ct), nil
}

type Task struct {
	gorm.Model
	Label       string `gorm:"not null"`
	Description string
	Status      TaskStatus `gorm:"type:enum('todo','in progress','completed','skip');default:todo;type:task_status"`
	Priority    uint       `gorm:"type:autoIncrement;uniqueIndex;not null"`

	TaskListID TaskList
}

type TaskList struct {
	gorm.Model
	Label string `gorm:"not null"`
	Tasks []Task `gorm:"constraint:OnDelete:CASCADE"`
}
