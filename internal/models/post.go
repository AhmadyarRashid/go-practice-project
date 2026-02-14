package models

import (
	"github.com/google/uuid"
)

// PostStatus represents post status
type PostStatus string

const (
	PostStatusDraft     PostStatus = "draft"
	PostStatusPublished PostStatus = "published"
	PostStatusArchived  PostStatus = "archived"
)

// Post represents a blog post or article
type Post struct {
	BaseModel
	Title       string     `gorm:"not null;size:255" json:"title"`
	Slug        string     `gorm:"uniqueIndex;not null;size:255" json:"slug"`
	Content     string     `gorm:"type:text" json:"content"`
	Excerpt     string     `gorm:"size:500" json:"excerpt"`
	FeaturedImage string   `gorm:"size:500" json:"featured_image,omitempty"`
	Status      PostStatus `gorm:"type:varchar(20);default:draft" json:"status"`
	ViewCount   int        `gorm:"default:0" json:"view_count"`

	// Foreign keys
	UserID      uuid.UUID  `gorm:"type:uuid;not null" json:"user_id"`

	// Relations
	User        *User      `gorm:"foreignKey:UserID" json:"author,omitempty"`
	Tags        []Tag      `gorm:"many2many:post_tags;" json:"tags,omitempty"`
}

// TableName returns the table name for Post model
func (Post) TableName() string {
	return "posts"
}

// IsPublished checks if the post is published
func (p *Post) IsPublished() bool {
	return p.Status == PostStatusPublished
}

// PostResponse is the response structure for post data
type PostResponse struct {
	ID            uuid.UUID     `json:"id"`
	Title         string        `json:"title"`
	Slug          string        `json:"slug"`
	Content       string        `json:"content"`
	Excerpt       string        `json:"excerpt"`
	FeaturedImage string        `json:"featured_image,omitempty"`
	Status        PostStatus    `json:"status"`
	ViewCount     int           `json:"view_count"`
	Author        *UserResponse `json:"author,omitempty"`
	Tags          []TagResponse `json:"tags,omitempty"`
	CreatedAt     string        `json:"created_at"`
	UpdatedAt     string        `json:"updated_at"`
}

// ToResponse converts Post to PostResponse
func (p *Post) ToResponse() *PostResponse {
	response := &PostResponse{
		ID:            p.ID,
		Title:         p.Title,
		Slug:          p.Slug,
		Content:       p.Content,
		Excerpt:       p.Excerpt,
		FeaturedImage: p.FeaturedImage,
		Status:        p.Status,
		ViewCount:     p.ViewCount,
		CreatedAt:     p.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:     p.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}

	if p.User != nil {
		response.Author = p.User.ToResponse()
	}

	if len(p.Tags) > 0 {
		response.Tags = make([]TagResponse, len(p.Tags))
		for i, tag := range p.Tags {
			response.Tags[i] = *tag.ToResponse()
		}
	}

	return response
}

// Tag represents a tag for categorizing posts
type Tag struct {
	BaseModel
	Name        string `gorm:"uniqueIndex;not null;size:100" json:"name"`
	Slug        string `gorm:"uniqueIndex;not null;size:100" json:"slug"`
	Description string `gorm:"size:500" json:"description,omitempty"`

	// Relations
	Posts       []Post `gorm:"many2many:post_tags;" json:"posts,omitempty"`
}

// TableName returns the table name for Tag model
func (Tag) TableName() string {
	return "tags"
}

// TagResponse is the response structure for tag data
type TagResponse struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Slug        string    `json:"slug"`
	Description string    `json:"description,omitempty"`
}

// ToResponse converts Tag to TagResponse
func (t *Tag) ToResponse() *TagResponse {
	return &TagResponse{
		ID:          t.ID,
		Name:        t.Name,
		Slug:        t.Slug,
		Description: t.Description,
	}
}
