package repositories

import (
	"context"
	"fmt"

	"github.com/go-redis/redis/v8"
	"github.com/jackc/pgx/v5/pgxpool"
	"gorm.io/gorm"
)

// RepositoryManager manages all repositories
type RepositoryManager struct {
	// Database connections
	PostgresPool *pgxpool.Pool
	GormDB       *gorm.DB
	RedisClient  *redis.Client

	// Repositories
	User    *UserRepository
	Post    *PostRepository
	Session *SessionRepository
	Cache   *CacheRepository
}

// NewRepositoryManager creates a new repository manager
func NewRepositoryManager(
	postgresPool *pgxpool.Pool,
	gormDB *gorm.DB,
	redisClient *redis.Client,
) *RepositoryManager {
	rm := &RepositoryManager{
		PostgresPool: postgresPool,
		GormDB:       gormDB,
		RedisClient:  redisClient,
	}

	// Initialize repositories
	rm.User = NewUserRepository(gormDB)
	rm.Post = NewPostRepository(gormDB)
	rm.Session = NewSessionRepository(gormDB)
	rm.Cache = NewCacheRepository(redisClient)

	return rm
}

// HealthCheck performs health checks on all repositories
func (rm *RepositoryManager) HealthCheck(ctx context.Context) map[string]string {
	health := make(map[string]string)

	// Check PostgreSQL connection
	if rm.PostgresPool != nil {
		if err := rm.PostgresPool.Ping(ctx); err != nil {
			health["postgres"] = "unhealthy: " + err.Error()
		} else {
			health["postgres"] = "healthy"
		}
	} else {
		health["postgres"] = "not connected"
	}

	// Check GORM connection
	if rm.GormDB != nil {
		if sqlDB, err := rm.GormDB.DB(); err == nil {
			if err := sqlDB.Ping(); err != nil {
				health["gorm"] = "unhealthy: " + err.Error()
			} else {
				health["gorm"] = "healthy"
			}
		} else {
			health["gorm"] = "unhealthy: " + err.Error()
		}
	} else {
		health["gorm"] = "not connected"
	}

	// Check Redis connection
	if rm.RedisClient != nil {
		if err := rm.RedisClient.Ping(ctx).Err(); err != nil {
			health["redis"] = "unhealthy: " + err.Error()
		} else {
			health["redis"] = "healthy"
		}
	} else {
		health["redis"] = "not connected"
	}

	return health
}

// Close closes all database connections
func (rm *RepositoryManager) Close() error {
	var errs []error

	if rm.PostgresPool != nil {
		rm.PostgresPool.Close()
	}

	if rm.GormDB != nil {
		if sqlDB, err := rm.GormDB.DB(); err == nil {
			if err := sqlDB.Close(); err != nil {
				errs = append(errs, err)
			}
		}
	}

	if rm.RedisClient != nil {
		if err := rm.RedisClient.Close(); err != nil {
			errs = append(errs, err)
		}
	}

	if len(errs) > 0 {
		return fmt.Errorf("errors closing connections: %v", errs)
	}

	return nil
}
