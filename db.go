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

type gormLogger struct {
	logMode glogger.LogLevel
}

func newGormLogger(logDebug bool) *gormLogger {
	gl := gormLogger{}
	if !logDebug {
		gl.logMode = glogger.Info
	}
	return &gl
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

func openDB(dbFile string, logDebug bool) (*gorm.DB, error) {
	log.Debug("Opening sqlite db...", "db", dbFile)

	conf := &gorm.Config{
		Logger: newGormLogger(logDebug),
	}
	db, err := gorm.Open(sqlite.Open(fmt.Sprintf("%s?_pragma=foreign_keys(1)", dbFile)), conf)
	if err != nil {
		return nil, err
	}

	log.Debug("Applying migrations...")

	if err = db.AutoMigrate(&TaskList{}, &Task{}); err != nil {
		defer tryCloseDB(db)
		return nil, fmt.Errorf("error applying migrations: %w", err)
	}

	return db, nil
}

func tryCloseDB(db *gorm.DB) {
	if db == nil {
		return
	}
	sdb, err := db.DB()
	if err != nil {
		return
	}
	_ = sdb.Close()
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

type TaskList struct {
	gorm.Model
	Label       string `gorm:"not null"`
	Date        time.Time
	Description string
	Tasks       []Task `gorm:"constraint:OnDelete:CASCADE"`
}

type Task struct {
	gorm.Model
	Label       string `gorm:"not null"`
	Description string
	Status      TaskStatus `gorm:"type:enum('todo','in progress','completed','skip');default:todo;type:TaskStatus"`
	Priority    uint       `gorm:"type:increment;uniqueIndex;not null"`

	TaskListID int
	TaskList   *TaskList
}
