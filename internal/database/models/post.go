package models

import (
	"time"
)

// Post represents a blog post or article
type Post struct {
	BaseModel
	Title       string     `json:"title" gorm:"not null" validate:"required,min=1,max=200"`
	Slug        string     `json:"slug" gorm:"uniqueIndex;not null" validate:"required"`
	Content     string     `json:"content" gorm:"type:text" validate:"required"`
	Excerpt     string     `json:"excerpt" gorm:"type:text" validate:"max=500"`
	Status      string     `json:"status" gorm:"default:'draft'" validate:"oneof=draft published archived"`
	AuthorID    uint       `json:"author_id" gorm:"not null"`
	Author      User       `json:"author" gorm:"foreignKey:AuthorID"`
	PublishedAt *time.Time `json:"published_at,omitempty"`
	ViewCount   int        `json:"view_count" gorm:"default:0"`
}

// TableName returns the table name for Post
func (Post) TableName() string {
	return "posts"
}

// IsPublished checks if the post is published
func (p *Post) IsPublished() bool {
	return p.Status == "published" && p.PublishedAt != nil
}

// IsDraft checks if the post is a draft
func (p *Post) IsDraft() bool {
	return p.Status == "draft"
}

// IsArchived checks if the post is archived
func (p *Post) IsArchived() bool {
	return p.Status == "archived"
}

// GetExcerpt returns the excerpt or a truncated version of content
func (p *Post) GetExcerpt() string {
	if p.Excerpt != "" {
		return p.Excerpt
	}

	// Truncate content to 150 characters
	if len(p.Content) > 150 {
		return p.Content[:150] + "..."
	}
	return p.Content
}

// Publish marks the post as published
func (p *Post) Publish() {
	now := time.Now()
	p.Status = "published"
	p.PublishedAt = &now
}

// Archive marks the post as archived
func (p *Post) Archive() {
	p.Status = "archived"
}

// IncrementViewCount increments the view count
func (p *Post) IncrementViewCount() {
	p.ViewCount++
}
