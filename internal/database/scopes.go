package database

import (
	"strings"

	"gorm.io/gorm"
)

// Search creates a search scope for LIKE queries across multiple fields
// This replaces raw SQL LIKE queries with a reusable GORM scope
//
// Usage:
//
//	db.Scopes(Search([]string{"first_name", "last_name", "email"}, "john")).Find(&users)
func Search(fields []string, query string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if query == "" || len(fields) == 0 {
			return db
		}

		searchQuery := "%" + query + "%"
		conditions := make([]string, len(fields))
		args := make([]interface{}, len(fields))

		for i, field := range fields {
			conditions[i] = field + " LIKE ?"
			args[i] = searchQuery
		}

		return db.Where(strings.Join(conditions, " OR "), args...)
	}
}

// OrderBy creates a scope for ordering results
//
// Usage:
//
//	db.Scopes(OrderBy("created_at", "desc")).Find(&posts)
func OrderBy(field string, direction string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if field == "" {
			return db
		}
		if direction == "" {
			direction = "asc"
		}
		return db.Order(field + " " + direction)
	}
}

// Status filters by status field
//
// Usage:
//
//	db.Scopes(Status("active")).Find(&users)
func Status(status string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if status == "" {
			return db
		}
		return db.Where("status = ?", status)
	}
}

// DateRange filters by date range
//
// Usage:
//
//	db.Scopes(DateRange("created_at", startDate, endDate)).Find(&posts)
func DateRange(field string, start, end interface{}) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if start != nil && end != nil {
			return db.Where(field+" BETWEEN ? AND ?", start, end)
		}
		if start != nil {
			return db.Where(field+" >= ?", start)
		}
		if end != nil {
			return db.Where(field+" <= ?", end)
		}
		return db
	}
}

// WithDeleted includes soft-deleted records
//
// Usage:
//
//	db.Scopes(WithDeleted()).Find(&users)
func WithDeleted() func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Unscoped()
	}
}

// OnlyDeleted returns only soft-deleted records
//
// Usage:
//
//	db.Scopes(OnlyDeleted()).Find(&users)
func OnlyDeleted() func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Unscoped().Where("deleted_at IS NOT NULL")
	}
}
