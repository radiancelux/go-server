package repositories

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
)

// CacheRepository handles Redis cache operations
type CacheRepository struct {
	client *redis.Client
}

// NewCacheRepository creates a new cache repository
func NewCacheRepository(client *redis.Client) *CacheRepository {
	return &CacheRepository{client: client}
}

// Set stores a value in cache with expiration
func (cr *CacheRepository) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	return cr.client.Set(ctx, key, value, expiration).Err()
}

// Get retrieves a value from cache
func (cr *CacheRepository) Get(ctx context.Context, key string) (string, error) {
	return cr.client.Get(ctx, key).Result()
}

// Delete removes a value from cache
func (cr *CacheRepository) Delete(ctx context.Context, key string) error {
	return cr.client.Del(ctx, key).Err()
}

// Exists checks if a key exists in cache
func (cr *CacheRepository) Exists(ctx context.Context, key string) (bool, error) {
	result, err := cr.client.Exists(ctx, key).Result()
	return result > 0, err
}

// SetUserSession stores a user session in cache
func (cr *CacheRepository) SetUserSession(ctx context.Context, userID uint, sessionID string, expiration time.Duration) error {
	key := fmt.Sprintf("session:%d:%s", userID, sessionID)
	return cr.Set(ctx, key, "active", expiration)
}

// GetUserSession retrieves a user session from cache
func (cr *CacheRepository) GetUserSession(ctx context.Context, userID uint, sessionID string) (bool, error) {
	key := fmt.Sprintf("session:%d:%s", userID, sessionID)
	exists, err := cr.Exists(ctx, key)
	return exists, err
}

// DeleteUserSession removes a user session from cache
func (cr *CacheRepository) DeleteUserSession(ctx context.Context, userID uint, sessionID string) error {
	key := fmt.Sprintf("session:%d:%s", userID, sessionID)
	return cr.Delete(ctx, key)
}

// SetPostCache stores a post in cache
func (cr *CacheRepository) SetPostCache(ctx context.Context, postID uint, post interface{}, expiration time.Duration) error {
	key := fmt.Sprintf("post:%d", postID)
	return cr.Set(ctx, key, post, expiration)
}

// GetPostCache retrieves a post from cache
func (cr *CacheRepository) GetPostCache(ctx context.Context, postID uint) (string, error) {
	key := fmt.Sprintf("post:%d", postID)
	return cr.Get(ctx, key)
}

// DeletePostCache removes a post from cache
func (cr *CacheRepository) DeletePostCache(ctx context.Context, postID uint) error {
	key := fmt.Sprintf("post:%d", postID)
	return cr.Delete(ctx, key)
}

// SetUserCache stores a user in cache
func (cr *CacheRepository) SetUserCache(ctx context.Context, userID uint, user interface{}, expiration time.Duration) error {
	key := fmt.Sprintf("user:%d", userID)
	return cr.Set(ctx, key, user, expiration)
}

// GetUserCache retrieves a user from cache
func (cr *CacheRepository) GetUserCache(ctx context.Context, userID uint) (string, error) {
	key := fmt.Sprintf("user:%d", userID)
	return cr.Get(ctx, key)
}

// DeleteUserCache removes a user from cache
func (cr *CacheRepository) DeleteUserCache(ctx context.Context, userID uint) error {
	key := fmt.Sprintf("user:%d", userID)
	return cr.Delete(ctx, key)
}

// SetListCache stores a list in cache
func (cr *CacheRepository) SetListCache(ctx context.Context, listKey string, data interface{}, expiration time.Duration) error {
	key := fmt.Sprintf("list:%s", listKey)
	return cr.Set(ctx, key, data, expiration)
}

// GetListCache retrieves a list from cache
func (cr *CacheRepository) GetListCache(ctx context.Context, listKey string) (string, error) {
	key := fmt.Sprintf("list:%s", listKey)
	return cr.Get(ctx, key)
}

// DeleteListCache removes a list from cache
func (cr *CacheRepository) DeleteListCache(ctx context.Context, listKey string) error {
	key := fmt.Sprintf("list:%s", listKey)
	return cr.Delete(ctx, key)
}

// FlushAll clears all cache entries
func (cr *CacheRepository) FlushAll(ctx context.Context) error {
	return cr.client.FlushAll(ctx).Err()
}

// Ping checks if Redis is accessible
func (cr *CacheRepository) Ping(ctx context.Context) error {
	return cr.client.Ping(ctx).Err()
}
