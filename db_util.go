package main

import (
	"context"

	"gorm.io/gorm"
)

type ModelQueryOpt func(db *gorm.DB) *gorm.DB

type AssociationQueryOpt func(assoc *gorm.Association) *gorm.Association

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
		opt(qdb)
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
		opt(qdb)
		if qdb.Error != nil {
			return nil, qdb.Error
		}
	}
	out := make([]T, 0)
	return out, qdb.Find(&out).Error
}

func CountAssociation[T any](ctx context.Context, db *gorm.DB, column string, opts ...AssociationQueryOpt) (int64, error) {
	assoc := db.WithContext(ctx).Model(new(T)).Association(column)
	if assoc.Error != nil {
		return 0, assoc.Error
	}
	for _, opt := range opts {
		opt(assoc)
		if assoc.Error != nil {
			return 0, assoc.Error
		}
	}
	return assoc.Count(), assoc.Error
}

func FindAssociation[T, E any](ctx context.Context, db *gorm.DB, column string, opts ...AssociationQueryOpt) ([]E, error) {
	assoc := db.WithContext(ctx).Model(new(T)).Association(column)
	if assoc.Error != nil {
		return nil, assoc.Error
	}
	for _, opt := range opts {
		opt(assoc)
		if assoc.Error != nil {
			return nil, assoc.Error
		}
	}
	out := make([]E, 0)
	return out, assoc.Find(&out)
}
