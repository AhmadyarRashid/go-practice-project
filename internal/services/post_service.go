package services

import (
	"context"
	"regexp"
	"strings"

	"github.com/google/uuid"
	"github.com/yourusername/go-enterprise-api/internal/models"
	"github.com/yourusername/go-enterprise-api/internal/repository"
	apperrors "github.com/yourusername/go-enterprise-api/pkg/errors"
	"github.com/yourusername/go-enterprise-api/pkg/logger"
)

// CreatePostRequest represents the create post request
type CreatePostRequest struct {
	Title         string   `json:"title" binding:"required"`
	Content       string   `json:"content" binding:"required"`
	Excerpt       string   `json:"excerpt"`
	FeaturedImage string   `json:"featured_image"`
	Status        string   `json:"status"`
	Tags          []string `json:"tags"`
}

// UpdatePostRequest represents the update post request
type UpdatePostRequest struct {
	Title         *string  `json:"title,omitempty"`
	Content       *string  `json:"content,omitempty"`
	Excerpt       *string  `json:"excerpt,omitempty"`
	FeaturedImage *string  `json:"featured_image,omitempty"`
	Status        *string  `json:"status,omitempty"`
	Tags          []string `json:"tags,omitempty"`
}

// PostService interface defines post service methods
type PostService interface {
	Create(ctx context.Context, userID uuid.UUID, req *CreatePostRequest) (*models.Post, error)
	GetByID(ctx context.Context, id uuid.UUID) (*models.Post, error)
	GetBySlug(ctx context.Context, slug string) (*models.Post, error)
	GetAll(ctx context.Context, page, pageSize int) ([]models.Post, int64, error)
	GetPublished(ctx context.Context, page, pageSize int) ([]models.Post, int64, error)
	GetByUser(ctx context.Context, userID uuid.UUID, page, pageSize int) ([]models.Post, int64, error)
	Update(ctx context.Context, id uuid.UUID, userID uuid.UUID, isAdmin bool, req *UpdatePostRequest) (*models.Post, error)
	Delete(ctx context.Context, id uuid.UUID, userID uuid.UUID, isAdmin bool) error
	Search(ctx context.Context, query string, page, pageSize int) ([]models.Post, int64, error)
	IncrementViews(ctx context.Context, id uuid.UUID) error
}

// postService implements PostService
type postService struct {
	postRepo repository.PostRepository
}

// NewPostService creates a new post service
func NewPostService(postRepo repository.PostRepository) PostService {
	return &postService{
		postRepo: postRepo,
	}
}

// Create creates a new post
func (s *postService) Create(ctx context.Context, userID uuid.UUID, req *CreatePostRequest) (*models.Post, error) {
	// Generate slug from title
	slug := generateSlug(req.Title)

	// Determine status
	status := models.PostStatusDraft
	if req.Status != "" {
		status = models.PostStatus(req.Status)
	}

	post := &models.Post{
		Title:         req.Title,
		Slug:          slug,
		Content:       req.Content,
		Excerpt:       req.Excerpt,
		FeaturedImage: req.FeaturedImage,
		Status:        status,
		UserID:        userID,
	}

	if err := s.postRepo.Create(ctx, post); err != nil {
		logger.Error("Failed to create post", logger.Err(err))
		return nil, apperrors.ErrInternal
	}

	// Fetch with relations
	return s.postRepo.FindWithAuthor(ctx, post.ID)
}

// GetByID retrieves a post by ID
func (s *postService) GetByID(ctx context.Context, id uuid.UUID) (*models.Post, error) {
	return s.postRepo.FindWithAuthor(ctx, id)
}

// GetBySlug retrieves a post by slug
func (s *postService) GetBySlug(ctx context.Context, slug string) (*models.Post, error) {
	return s.postRepo.FindBySlug(ctx, slug)
}

// GetAll retrieves all posts with pagination
func (s *postService) GetAll(ctx context.Context, page, pageSize int) ([]models.Post, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	return s.postRepo.FindAllWithAuthor(ctx, page, pageSize)
}

// GetPublished retrieves all published posts
func (s *postService) GetPublished(ctx context.Context, page, pageSize int) ([]models.Post, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	return s.postRepo.FindPublished(ctx, page, pageSize)
}

// GetByUser retrieves posts by user ID
func (s *postService) GetByUser(ctx context.Context, userID uuid.UUID, page, pageSize int) ([]models.Post, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	return s.postRepo.FindByUserID(ctx, userID, page, pageSize)
}

// Update updates a post
func (s *postService) Update(ctx context.Context, id uuid.UUID, userID uuid.UUID, isAdmin bool, req *UpdatePostRequest) (*models.Post, error) {
	post, err := s.postRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Check ownership or admin
	if post.UserID != userID && !isAdmin {
		return nil, apperrors.ErrForbidden
	}

	// Update fields if provided
	if req.Title != nil {
		post.Title = *req.Title
		post.Slug = generateSlug(*req.Title)
	}
	if req.Content != nil {
		post.Content = *req.Content
	}
	if req.Excerpt != nil {
		post.Excerpt = *req.Excerpt
	}
	if req.FeaturedImage != nil {
		post.FeaturedImage = *req.FeaturedImage
	}
	if req.Status != nil {
		post.Status = models.PostStatus(*req.Status)
	}

	if err := s.postRepo.Update(ctx, post); err != nil {
		logger.Error("Failed to update post", logger.Err(err))
		return nil, apperrors.ErrInternal
	}

	return s.postRepo.FindWithAuthor(ctx, post.ID)
}

// Delete deletes a post
func (s *postService) Delete(ctx context.Context, id uuid.UUID, userID uuid.UUID, isAdmin bool) error {
	post, err := s.postRepo.FindByID(ctx, id)
	if err != nil {
		return err
	}

	// Check ownership or admin
	if post.UserID != userID && !isAdmin {
		return apperrors.ErrForbidden
	}

	if err := s.postRepo.Delete(ctx, id); err != nil {
		logger.Error("Failed to delete post", logger.Err(err))
		return apperrors.ErrInternal
	}

	return nil
}

// Search searches for posts
func (s *postService) Search(ctx context.Context, query string, page, pageSize int) ([]models.Post, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	return s.postRepo.SearchPosts(ctx, query, page, pageSize)
}

// IncrementViews increments the view count
func (s *postService) IncrementViews(ctx context.Context, id uuid.UUID) error {
	return s.postRepo.IncrementViewCount(ctx, id)
}

// generateSlug generates a URL-friendly slug from a title
func generateSlug(title string) string {
	// Convert to lowercase
	slug := strings.ToLower(title)

	// Replace spaces with hyphens
	slug = strings.ReplaceAll(slug, " ", "-")

	// Remove non-alphanumeric characters except hyphens
	reg := regexp.MustCompile("[^a-z0-9-]+")
	slug = reg.ReplaceAllString(slug, "")

	// Remove multiple consecutive hyphens
	reg = regexp.MustCompile("-+")
	slug = reg.ReplaceAllString(slug, "-")

	// Trim hyphens from start and end
	slug = strings.Trim(slug, "-")

	// Add UUID suffix for uniqueness
	slug = slug + "-" + uuid.New().String()[:8]

	return slug
}
