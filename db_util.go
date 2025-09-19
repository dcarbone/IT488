package main

import (
	"context"
	"fmt"

	"gorm.io/gorm"
)

type ModelQueryOpt func(db *gorm.DB) *gorm.DB

type AssociationQueryOpt func(assoc *gorm.Association) *gorm.Association

func WithLimit(limit int) ModelQueryOpt {
	return func(db *gorm.DB) *gorm.DB {
		return db.Limit(limit)
	}
}

func WithSort(clause any) ModelQueryOpt {
	return func(db *gorm.DB) *gorm.DB {
		return db.Order(clause)
	}
}

func WithPreload(query string, args ...any) ModelQueryOpt {
	return func(db *gorm.DB) *gorm.DB {
		return db.Preload(query, args...)
	}
}

func CountModel[T any](ctx context.Context, db *gorm.DB, opts ...ModelQueryOpt) (int64, error) {
	qdb := db.WithContext(ctx).Model(new(T))
	for _, opt := range opts {
		qdb = opt(qdb)
		if qdb.Error != nil {
			return 0, db.Error
		}
	}
	var count int64
	return count, qdb.Count(&count).Error
}

func FindModel[T any](ctx context.Context, db *gorm.DB, opts ...ModelQueryOpt) ([]T, error) {
	qdb := db.WithContext(ctx).Model(new(T))
	for _, opt := range opts {
		qdb = opt(qdb)
		if qdb.Error != nil {
			return nil, qdb.Error
		}
	}
	out := make([]T, 0)
	return out, qdb.Find(&out).Error
}

func FindOneModel[T any](ctx context.Context, db *gorm.DB, opts ...ModelQueryOpt) (*T, error) {
	models, err := FindModel[T](ctx, db, append(opts, WithLimit(1))...)
	if err != nil || len(models) == 0 {
		return nil, err
	}
	return &models[0], nil
}

func CountAssociation[T any](ctx context.Context, db *gorm.DB, base T, column string, opts ...AssociationQueryOpt) (int64, error) {
	assoc := db.WithContext(ctx).Model(&base).Association(column)
	if assoc.Error != nil {
		return 0, assoc.Error
	}
	for _, opt := range opts {
		assoc = opt(assoc)
		if assoc.Error != nil {
			return 0, assoc.Error
		}
	}
	return assoc.Count(), assoc.Error
}

func FindAssociation[T, E any](ctx context.Context, db *gorm.DB, base T, column string, opts ...AssociationQueryOpt) ([]E, error) {
	assoc := db.WithContext(ctx).Model(&base).Association(column)
	if assoc.Error != nil {
		return nil, assoc.Error
	}
	for _, opt := range opts {
		assoc = opt(assoc)
		if assoc.Error != nil {
			return nil, assoc.Error
		}
	}
	out := make([]E, 0)
	return out, assoc.Find(&out)
}

func todaysTasksModelQueryOpt() ModelQueryOpt {
	return func(db *gorm.DB) *gorm.DB {
		return WithSort("due_date asc")(WithSort("id asc")(WithPreload("TaskList")(db))).
			Where("date(`tasks`.`due_date`, 'localtime') = date('now', 'localtime')")
	}
}

func GetListForTask(ctx context.Context, db *gorm.DB, task Task) *TaskList {
	if task.TaskList != nil {
		return task.TaskList
	}
	if task.TaskListID.Valid {
		taskList, err := FindOneModel[TaskList](ctx, db, func(db *gorm.DB) *gorm.DB {
			return db.Where("ID = ?", task.TaskListID.V)
		})
		if err != nil {
			panic(fmt.Sprintf("error loading task list with ID %d: %v", task.TaskListID.V, err))
		}
		return taskList
	}
	return nil
}
