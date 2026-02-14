package repository

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/yourusername/go-enterprise-api/internal/models"
	apperrors "github.com/yourusername/go-enterprise-api/pkg/errors"
	"gorm.io/gorm"
)

// PostRepository interface defines post-specific repository methods
type PostRepository interface {
	Repository[models.Post]
	FindBySlug(ctx context.Context, slug string) (*models.Post, error)
	FindByUserID(ctx context.Context, userID uuid.UUID, page, pageSize int) ([]models.Post, int64, error)
	FindPublished(ctx context.Context, page, pageSize int) ([]models.Post, int64, error)
	FindByStatus(ctx context.Context, status models.PostStatus, page, pageSize int) ([]models.Post, int64, error)
	IncrementViewCount(ctx context.Context, postID uuid.UUID) error
	FindWithAuthor(ctx context.Context, id uuid.UUID) (*models.Post, error)
	FindAllWithAuthor(ctx context.Context, page, pageSize int) ([]models.Post, int64, error)
	SearchPosts(ctx context.Context, query string, page, pageSize int) ([]models.Post, int64, error)
	AddTag(ctx context.Context, postID, tagID uuid.UUID) error
	RemoveTag(ctx context.Context, postID, tagID uuid.UUID) error
	FindByTag(ctx context.Context, tagSlug string, page, pageSize int) ([]models.Post, int64, error)
}

// postRepository implements PostRepository
type postRepository struct {
	*BaseRepository[models.Post]
}

// NewPostRepository creates a new post repository
func NewPostRepository(db *gorm.DB) PostRepository {
	return &postRepository{
		BaseRepository: NewBaseRepository[models.Post](db),
	}
}

// FindBySlug finds a post by slug
func (r *postRepository) FindBySlug(ctx context.Context, slug string) (*models.Post, error) {
	var post models.Post
	err := r.DB.WithContext(ctx).Preload("User").Preload("Tags").Where("slug = ?", slug).First(&post).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperrors.ErrNotFound.WithDetails("Post not found")
		}
		return nil, err
	}
	return &post, nil
}

// FindByUserID finds posts by user ID
func (r *postRepository) FindByUserID(ctx context.Context, userID uuid.UUID, page, pageSize int) ([]models.Post, int64, error) {
	var posts []models.Post
	var total int64

	err := r.DB.WithContext(ctx).Model(&models.Post{}).Where("user_id = ?", userID).Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	err = r.DB.WithContext(ctx).
		Preload("Tags").
		Where("user_id = ?", userID).
		Order("created_at DESC").
		Offset(offset).Limit(pageSize).
		Find(&posts).Error

	return posts, total, err
}

// FindPublished finds all published posts
func (r *postRepository) FindPublished(ctx context.Context, page, pageSize int) ([]models.Post, int64, error) {
	return r.FindByStatus(ctx, models.PostStatusPublished, page, pageSize)
}

// FindByStatus finds posts by status
func (r *postRepository) FindByStatus(ctx context.Context, status models.PostStatus, page, pageSize int) ([]models.Post, int64, error) {
	var posts []models.Post
	var total int64

	err := r.DB.WithContext(ctx).Model(&models.Post{}).Where("status = ?", status).Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	err = r.DB.WithContext(ctx).
		Preload("User").
		Preload("Tags").
		Where("status = ?", status).
		Order("created_at DESC").
		Offset(offset).Limit(pageSize).
		Find(&posts).Error

	return posts, total, err
}

// IncrementViewCount increments the view count for a post
func (r *postRepository) IncrementViewCount(ctx context.Context, postID uuid.UUID) error {
	return r.DB.WithContext(ctx).Model(&models.Post{}).Where("id = ?", postID).UpdateColumn("view_count", gorm.Expr("view_count + ?", 1)).Error
}

// FindWithAuthor finds a post with its author
func (r *postRepository) FindWithAuthor(ctx context.Context, id uuid.UUID) (*models.Post, error) {
	var post models.Post
	err := r.DB.WithContext(ctx).Preload("User").Preload("Tags").First(&post, "id = ?", id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperrors.ErrNotFound.WithDetails("Post not found")
		}
		return nil, err
	}
	return &post, nil
}

// FindAllWithAuthor finds all posts with their authors
func (r *postRepository) FindAllWithAuthor(ctx context.Context, page, pageSize int) ([]models.Post, int64, error) {
	var posts []models.Post
	var total int64

	err := r.DB.WithContext(ctx).Model(&models.Post{}).Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	err = r.DB.WithContext(ctx).
		Preload("User").
		Preload("Tags").
		Order("created_at DESC").
		Offset(offset).Limit(pageSize).
		Find(&posts).Error

	return posts, total, err
}

// SearchPosts searches for posts by title or content
func (r *postRepository) SearchPosts(ctx context.Context, query string, page, pageSize int) ([]models.Post, int64, error) {
	var posts []models.Post
	var total int64

	searchQuery := "%" + query + "%"

	err := r.DB.WithContext(ctx).Model(&models.Post{}).
		Where("title LIKE ? OR content LIKE ?", searchQuery, searchQuery).
		Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	err = r.DB.WithContext(ctx).
		Preload("User").
		Preload("Tags").
		Where("title LIKE ? OR content LIKE ?", searchQuery, searchQuery).
		Order("created_at DESC").
		Offset(offset).Limit(pageSize).
		Find(&posts).Error

	return posts, total, err
}

// AddTag adds a tag to a post
func (r *postRepository) AddTag(ctx context.Context, postID, tagID uuid.UUID) error {
	return r.DB.WithContext(ctx).Exec("INSERT INTO post_tags (post_id, tag_id) VALUES (?, ?)", postID, tagID).Error
}

// RemoveTag removes a tag from a post
func (r *postRepository) RemoveTag(ctx context.Context, postID, tagID uuid.UUID) error {
	return r.DB.WithContext(ctx).Exec("DELETE FROM post_tags WHERE post_id = ? AND tag_id = ?", postID, tagID).Error
}

// FindByTag finds posts by tag slug
func (r *postRepository) FindByTag(ctx context.Context, tagSlug string, page, pageSize int) ([]models.Post, int64, error) {
	var posts []models.Post
	var total int64

	// Count total
	err := r.DB.WithContext(ctx).Model(&models.Post{}).
		Joins("JOIN post_tags ON post_tags.post_id = posts.id").
		Joins("JOIN tags ON tags.id = post_tags.tag_id").
		Where("tags.slug = ?", tagSlug).
		Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	// Get paginated results
	offset := (page - 1) * pageSize
	err = r.DB.WithContext(ctx).
		Preload("User").
		Preload("Tags").
		Joins("JOIN post_tags ON post_tags.post_id = posts.id").
		Joins("JOIN tags ON tags.id = post_tags.tag_id").
		Where("tags.slug = ?", tagSlug).
		Order("posts.created_at DESC").
		Offset(offset).Limit(pageSize).
		Find(&posts).Error

	return posts, total, err
}

// FindByID overrides base to include error handling
func (r *postRepository) FindByID(ctx context.Context, id uuid.UUID) (*models.Post, error) {
	var post models.Post
	err := r.DB.WithContext(ctx).First(&post, "id = ?", id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperrors.ErrNotFound.WithDetails("Post not found")
		}
		return nil, err
	}
	return &post, nil
}
