package repository

import (
	"context"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Repository is the base repository interface
type Repository[T any] interface {
	Create(ctx context.Context, entity *T) error
	FindByID(ctx context.Context, id uuid.UUID) (*T, error)
	FindAll(ctx context.Context, page, pageSize int) ([]T, int64, error)
	Update(ctx context.Context, entity *T) error
	Delete(ctx context.Context, id uuid.UUID) error
	Count(ctx context.Context) (int64, error)
}

// BaseRepository provides common repository functionality
type BaseRepository[T any] struct {
	DB *gorm.DB
}

// NewBaseRepository creates a new base repository
func NewBaseRepository[T any](db *gorm.DB) *BaseRepository[T] {
	return &BaseRepository[T]{DB: db}
}

// Create creates a new entity
func (r *BaseRepository[T]) Create(ctx context.Context, entity *T) error {
	return r.DB.WithContext(ctx).Create(entity).Error
}

// FindByID finds an entity by ID
func (r *BaseRepository[T]) FindByID(ctx context.Context, id uuid.UUID) (*T, error) {
	var entity T
	err := r.DB.WithContext(ctx).First(&entity, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &entity, nil
}

// FindAll finds all entities with pagination
func (r *BaseRepository[T]) FindAll(ctx context.Context, page, pageSize int) ([]T, int64, error) {
	var entities []T
	var total int64

	// Get total count
	if err := r.DB.WithContext(ctx).Model(new(T)).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get paginated results
	offset := (page - 1) * pageSize
	if err := r.DB.WithContext(ctx).Offset(offset).Limit(pageSize).Find(&entities).Error; err != nil {
		return nil, 0, err
	}

	return entities, total, nil
}

// Update updates an entity
func (r *BaseRepository[T]) Update(ctx context.Context, entity *T) error {
	return r.DB.WithContext(ctx).Save(entity).Error
}

// Delete soft deletes an entity
func (r *BaseRepository[T]) Delete(ctx context.Context, id uuid.UUID) error {
	return r.DB.WithContext(ctx).Delete(new(T), "id = ?", id).Error
}

// HardDelete permanently deletes an entity
func (r *BaseRepository[T]) HardDelete(ctx context.Context, id uuid.UUID) error {
	return r.DB.WithContext(ctx).Unscoped().Delete(new(T), "id = ?", id).Error
}

// Count counts all entities
func (r *BaseRepository[T]) Count(ctx context.Context) (int64, error) {
	var count int64
	err := r.DB.WithContext(ctx).Model(new(T)).Count(&count).Error
	return count, err
}

// FindByField finds entities by a specific field
func (r *BaseRepository[T]) FindByField(ctx context.Context, field string, value interface{}) ([]T, error) {
	var entities []T
	err := r.DB.WithContext(ctx).Where(field+" = ?", value).Find(&entities).Error
	return entities, err
}

// FindOneByField finds a single entity by a specific field
func (r *BaseRepository[T]) FindOneByField(ctx context.Context, field string, value interface{}) (*T, error) {
	var entity T
	err := r.DB.WithContext(ctx).Where(field+" = ?", value).First(&entity).Error
	if err != nil {
		return nil, err
	}
	return &entity, nil
}

// Exists checks if an entity exists by ID
func (r *BaseRepository[T]) Exists(ctx context.Context, id uuid.UUID) (bool, error) {
	var count int64
	err := r.DB.WithContext(ctx).Model(new(T)).Where("id = ?", id).Count(&count).Error
	return count > 0, err
}

// Transaction executes a function within a transaction
func (r *BaseRepository[T]) Transaction(ctx context.Context, fn func(tx *gorm.DB) error) error {
	return r.DB.WithContext(ctx).Transaction(fn)
}
