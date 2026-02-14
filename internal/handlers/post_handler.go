package handlers

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/yourusername/go-enterprise-api/internal/middleware"
	"github.com/yourusername/go-enterprise-api/internal/models"
	"github.com/yourusername/go-enterprise-api/internal/services"
	"github.com/yourusername/go-enterprise-api/pkg/response"
	"github.com/yourusername/go-enterprise-api/pkg/validator"
)

// PostHandler handles post-related requests
type PostHandler struct {
	postService services.PostService
}

// NewPostHandler creates a new post handler
func NewPostHandler(postService services.PostService) *PostHandler {
	return &PostHandler{
		postService: postService,
	}
}

// CreatePostRequest represents the create post request body
type CreatePostRequest struct {
	Title         string   `json:"title" binding:"required"`
	Content       string   `json:"content" binding:"required"`
	Excerpt       string   `json:"excerpt"`
	FeaturedImage string   `json:"featured_image"`
	Status        string   `json:"status"`
	Tags          []string `json:"tags"`
}

// UpdatePostRequest represents the update post request body
type UpdatePostRequest struct {
	Title         *string  `json:"title,omitempty"`
	Content       *string  `json:"content,omitempty"`
	Excerpt       *string  `json:"excerpt,omitempty"`
	FeaturedImage *string  `json:"featured_image,omitempty"`
	Status        *string  `json:"status,omitempty"`
	Tags          []string `json:"tags,omitempty"`
}

// Create creates a new post
// @Summary Create a new post
// @Description Create a new blog post
// @Tags posts
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body CreatePostRequest true "Post data"
// @Success 201 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Router /posts [post]
func (h *PostHandler) Create(c *gin.Context) {
	var req CreatePostRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request body")
		return
	}

	// Validate request
	v := validator.New()
	v.Required("title", req.Title, "")
	v.MaxLength("title", req.Title, 255, "")
	v.Required("content", req.Content, "")

	if errs := v.Validate(); errs != nil {
		response.ValidationError(c, errs)
		return
	}

	user := middleware.MustGetUser(c)

	serviceReq := &services.CreatePostRequest{
		Title:         req.Title,
		Content:       req.Content,
		Excerpt:       req.Excerpt,
		FeaturedImage: req.FeaturedImage,
		Status:        req.Status,
		Tags:          req.Tags,
	}

	post, err := h.postService.Create(c.Request.Context(), user.ID, serviceReq)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Created(c, gin.H{
		"post": post.ToResponse(),
	})
}

// GetAll returns all posts with pagination
// @Summary Get all posts
// @Description Get a paginated list of posts (published only for non-authenticated users)
// @Tags posts
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Page size" default(10)
// @Success 200 {object} response.Response
// @Router /posts [get]
func (h *PostHandler) GetAll(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	// Check if user is authenticated and is admin
	user, exists := middleware.GetUser(c)
	isAdmin := exists && user.IsAdmin()

	var posts []models.Post
	var total int64
	var err error

	if isAdmin {
		// Admin can see all posts
		posts, total, err = h.postService.GetAll(c.Request.Context(), page, pageSize)
	} else {
		// Non-admin only sees published posts
		posts, total, err = h.postService.GetPublished(c.Request.Context(), page, pageSize)
	}

	if err != nil {
		response.Error(c, err)
		return
	}

	// Convert to response
	postResponses := make([]*models.PostResponse, len(posts))
	for i, post := range posts {
		postResponses[i] = post.ToResponse()
	}

	response.Paginated(c, postResponses, page, pageSize, total)
}

// GetByID returns a post by ID
// @Summary Get post by ID
// @Description Get a specific post by its ID
// @Tags posts
// @Accept json
// @Produce json
// @Param id path string true "Post ID"
// @Success 200 {object} response.Response
// @Failure 404 {object} response.Response
// @Router /posts/{id} [get]
func (h *PostHandler) GetByID(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "Invalid post ID")
		return
	}

	post, err := h.postService.GetByID(c.Request.Context(), id)
	if err != nil {
		response.Error(c, err)
		return
	}

	// Check if post is published or user has permission
	user, exists := middleware.GetUser(c)
	if !post.IsPublished() {
		if !exists || (post.UserID != user.ID && !user.IsAdmin()) {
			response.NotFound(c, "Post not found")
			return
		}
	}

	// Increment view count
	_ = h.postService.IncrementViews(c.Request.Context(), id)

	response.Success(c, gin.H{
		"post": post.ToResponse(),
	})
}

// GetBySlug returns a post by slug
// @Summary Get post by slug
// @Description Get a specific post by its URL slug
// @Tags posts
// @Accept json
// @Produce json
// @Param slug path string true "Post slug"
// @Success 200 {object} response.Response
// @Failure 404 {object} response.Response
// @Router /posts/slug/{slug} [get]
func (h *PostHandler) GetBySlug(c *gin.Context) {
	slug := c.Param("slug")

	post, err := h.postService.GetBySlug(c.Request.Context(), slug)
	if err != nil {
		response.Error(c, err)
		return
	}

	// Check if post is published or user has permission
	user, exists := middleware.GetUser(c)
	if !post.IsPublished() {
		if !exists || (post.UserID != user.ID && !user.IsAdmin()) {
			response.NotFound(c, "Post not found")
			return
		}
	}

	// Increment view count
	_ = h.postService.IncrementViews(c.Request.Context(), post.ID)

	response.Success(c, gin.H{
		"post": post.ToResponse(),
	})
}

// Update updates a post
// @Summary Update post
// @Description Update a post (owner or admin only)
// @Tags posts
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Post ID"
// @Param request body UpdatePostRequest true "Update data"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 403 {object} response.Response
// @Failure 404 {object} response.Response
// @Router /posts/{id} [put]
func (h *PostHandler) Update(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "Invalid post ID")
		return
	}

	var req UpdatePostRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request body")
		return
	}

	user := middleware.MustGetUser(c)

	serviceReq := &services.UpdatePostRequest{
		Title:         req.Title,
		Content:       req.Content,
		Excerpt:       req.Excerpt,
		FeaturedImage: req.FeaturedImage,
		Status:        req.Status,
		Tags:          req.Tags,
	}

	post, err := h.postService.Update(c.Request.Context(), id, user.ID, user.IsAdmin(), serviceReq)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, gin.H{
		"post": post.ToResponse(),
	})
}

// Delete deletes a post
// @Summary Delete post
// @Description Delete a post (owner or admin only)
// @Tags posts
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Post ID"
// @Success 204
// @Failure 401 {object} response.Response
// @Failure 403 {object} response.Response
// @Failure 404 {object} response.Response
// @Router /posts/{id} [delete]
func (h *PostHandler) Delete(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "Invalid post ID")
		return
	}

	user := middleware.MustGetUser(c)

	if err := h.postService.Delete(c.Request.Context(), id, user.ID, user.IsAdmin()); err != nil {
		response.Error(c, err)
		return
	}

	response.NoContent(c)
}

// GetMyPosts returns the current user's posts
// @Summary Get my posts
// @Description Get the authenticated user's posts
// @Tags posts
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Page size" default(10)
// @Success 200 {object} response.Response
// @Failure 401 {object} response.Response
// @Router /posts/my [get]
func (h *PostHandler) GetMyPosts(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	user := middleware.MustGetUser(c)

	posts, total, err := h.postService.GetByUser(c.Request.Context(), user.ID, page, pageSize)
	if err != nil {
		response.Error(c, err)
		return
	}

	// Convert to response
	postResponses := make([]*models.PostResponse, len(posts))
	for i, post := range posts {
		postResponses[i] = post.ToResponse()
	}

	response.Paginated(c, postResponses, page, pageSize, total)
}

// Search searches for posts
// @Summary Search posts
// @Description Search for posts by title or content
// @Tags posts
// @Accept json
// @Produce json
// @Param q query string true "Search query"
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Page size" default(10)
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Router /posts/search [get]
func (h *PostHandler) Search(c *gin.Context) {
	query := c.Query("q")
	if query == "" {
		response.BadRequest(c, "Search query is required")
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	posts, total, err := h.postService.Search(c.Request.Context(), query, page, pageSize)
	if err != nil {
		response.Error(c, err)
		return
	}

	// Convert to response
	postResponses := make([]*models.PostResponse, len(posts))
	for i, post := range posts {
		postResponses[i] = post.ToResponse()
	}

	response.Paginated(c, postResponses, page, pageSize, total)
}
