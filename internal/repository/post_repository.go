package repository

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/yourusername/go-enterprise-api/internal/database"
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
// Uses GORM Scopes instead of raw SQL LIKE queries
func (r *postRepository) SearchPosts(ctx context.Context, query string, page, pageSize int) ([]models.Post, int64, error) {
	var posts []models.Post
	var total int64

	// Define searchable fields
	searchFields := []string{"title", "content"}

	err := r.DB.WithContext(ctx).Model(&models.Post{}).
		Scopes(database.Search(searchFields, query)).
		Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	err = r.DB.WithContext(ctx).
		Preload("User").
		Preload("Tags").
		Scopes(database.Search(searchFields, query)).
		Order("created_at DESC").
		Offset(offset).Limit(pageSize).
		Find(&posts).Error

	return posts, total, err
}

// AddTag adds a tag to a post using GORM Association
// No raw SQL needed - uses GORM's built-in many-to-many support
func (r *postRepository) AddTag(ctx context.Context, postID, tagID uuid.UUID) error {
	post := &models.Post{}
	post.ID = postID

	tag := &models.Tag{}
	tag.ID = tagID

	return r.DB.WithContext(ctx).Model(post).Association("Tags").Append(tag)
}

// RemoveTag removes a tag from a post using GORM Association
// No raw SQL needed - uses GORM's built-in many-to-many support
func (r *postRepository) RemoveTag(ctx context.Context, postID, tagID uuid.UUID) error {
	post := &models.Post{}
	post.ID = postID

	tag := &models.Tag{}
	tag.ID = tagID

	return r.DB.WithContext(ctx).Model(post).Association("Tags").Delete(tag)
}

// FindByTag finds posts by tag slug using GORM Association
// No raw SQL JOIN needed - uses GORM's relationship features
func (r *postRepository) FindByTag(ctx context.Context, tagSlug string, page, pageSize int) ([]models.Post, int64, error) {
	var posts []models.Post
	var total int64

	// First, find the tag by slug
	var tag models.Tag
	err := r.DB.WithContext(ctx).Where("slug = ?", tagSlug).First(&tag).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return []models.Post{}, 0, nil // Return empty if tag not found
		}
		return nil, 0, err
	}

	// Count total posts with this tag using Association
	total = r.DB.WithContext(ctx).Model(&tag).Association("Posts").Count()

	// Get paginated posts using Association
	offset := (page - 1) * pageSize
	err = r.DB.WithContext(ctx).
		Model(&tag).
		Preload("User").
		Preload("Tags").
		Order("created_at DESC").
		Offset(offset).Limit(pageSize).
		Association("Posts").
		Find(&posts)

	// If Association.Find doesn't support Preload, use alternative approach
	if len(posts) > 0 {
		// Reload posts with relationships
		var postIDs []uuid.UUID
		for _, p := range posts {
			postIDs = append(postIDs, p.ID)
		}
		err = r.DB.WithContext(ctx).
			Preload("User").
			Preload("Tags").
			Where("id IN ?", postIDs).
			Order("created_at DESC").
			Find(&posts).Error
	}

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
