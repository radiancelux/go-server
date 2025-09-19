package repositories

import (
	"context"

	"go-server/internal/database/models"
	"gorm.io/gorm"
)

// PostRepository handles post-related database operations
type PostRepository struct {
	db *gorm.DB
}

// NewPostRepository creates a new post repository
func NewPostRepository(db *gorm.DB) *PostRepository {
	return &PostRepository{db: db}
}

// CreatePost creates a new post
func (pr *PostRepository) CreatePost(ctx context.Context, post *models.Post) error {
	return pr.db.WithContext(ctx).Create(post).Error
}

// GetPostByID retrieves a post by ID
func (pr *PostRepository) GetPostByID(ctx context.Context, id uint) (*models.Post, error) {
	var post models.Post
	err := pr.db.WithContext(ctx).
		Preload("Author").
		Preload("Categories").
		First(&post, id).Error
	if err != nil {
		return nil, err
	}
	return &post, nil
}

// GetPostBySlug retrieves a post by slug
func (pr *PostRepository) GetPostBySlug(ctx context.Context, slug string) (*models.Post, error) {
	var post models.Post
	err := pr.db.WithContext(ctx).
		Preload("Author").
		Preload("Categories").
		Where("slug = ?", slug).
		First(&post).Error
	if err != nil {
		return nil, err
	}
	return &post, nil
}

// UpdatePost updates a post
func (pr *PostRepository) UpdatePost(ctx context.Context, post *models.Post) error {
	return pr.db.WithContext(ctx).Save(post).Error
}

// DeletePost soft deletes a post
func (pr *PostRepository) DeletePost(ctx context.Context, id uint) error {
	return pr.db.WithContext(ctx).Delete(&models.Post{}, id).Error
}

// ListPosts retrieves posts with pagination
func (pr *PostRepository) ListPosts(ctx context.Context, offset, limit int) ([]models.Post, error) {
	var posts []models.Post
	err := pr.db.WithContext(ctx).
		Preload("Author").
		Preload("Categories").
		Offset(offset).
		Limit(limit).
		Find(&posts).Error
	return posts, err
}

// ListPublishedPosts retrieves published posts with pagination
func (pr *PostRepository) ListPublishedPosts(ctx context.Context, offset, limit int) ([]models.Post, error) {
	var posts []models.Post
	err := pr.db.WithContext(ctx).
		Preload("Author").
		Preload("Categories").
		Where("status = ? AND published_at IS NOT NULL", "published").
		Order("published_at DESC").
		Offset(offset).
		Limit(limit).
		Find(&posts).Error
	return posts, err
}

// ListPostsByAuthor retrieves posts by author
func (pr *PostRepository) ListPostsByAuthor(ctx context.Context, authorID uint, offset, limit int) ([]models.Post, error) {
	var posts []models.Post
	err := pr.db.WithContext(ctx).
		Preload("Author").
		Preload("Categories").
		Where("author_id = ?", authorID).
		Offset(offset).
		Limit(limit).
		Find(&posts).Error
	return posts, err
}

// IncrementViewCount increments the view count for a post
func (pr *PostRepository) IncrementViewCount(ctx context.Context, id uint) error {
	return pr.db.WithContext(ctx).
		Model(&models.Post{}).
		Where("id = ?", id).
		Update("view_count", gorm.Expr("view_count + 1")).Error
}

// CountPosts returns the total number of posts
func (pr *PostRepository) CountPosts(ctx context.Context) (int64, error) {
	var count int64
	err := pr.db.WithContext(ctx).Model(&models.Post{}).Count(&count).Error
	return count, err
}

// CountPublishedPosts returns the total number of published posts
func (pr *PostRepository) CountPublishedPosts(ctx context.Context) (int64, error) {
	var count int64
	err := pr.db.WithContext(ctx).
		Model(&models.Post{}).
		Where("status = ? AND published_at IS NOT NULL", "published").
		Count(&count).Error
	return count, err
}
