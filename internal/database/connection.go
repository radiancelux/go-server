package database

import (
	"context"
	"fmt"
	"log"

	"github.com/go-redis/redis/v8"
	"github.com/jackc/pgx/v5/pgxpool"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// DatabaseManager manages database connections
type DatabaseManager struct {
	PostgresPool *pgxpool.Pool
	GormDB       *gorm.DB
	RedisClient  *redis.Client
	Config       *DatabaseConfig
}

// NewDatabaseManager creates a new database manager
func NewDatabaseManager(config *DatabaseConfig) *DatabaseManager {
	return &DatabaseManager{
		Config: config,
	}
}

// ConnectPostgres establishes PostgreSQL connection using pgxpool
func (dm *DatabaseManager) ConnectPostgres(ctx context.Context) error {
	config, err := pgxpool.ParseConfig(dm.Config.GetPostgresDSN())
	if err != nil {
		return fmt.Errorf("failed to parse postgres config: %w", err)
	}

	// Configure connection pool
	config.MaxConns = int32(dm.Config.MaxConnections)
	config.MinConns = 1
	config.MaxConnLifetime = dm.Config.ConnMaxLifetime
	config.MaxConnIdleTime = dm.Config.ConnMaxIdleTime

	// Connect to PostgreSQL
	dm.PostgresPool, err = pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		return fmt.Errorf("failed to connect to postgres: %w", err)
	}

	// Test connection
	if err := dm.PostgresPool.Ping(ctx); err != nil {
		return fmt.Errorf("failed to ping postgres: %w", err)
	}

	log.Println("✅ PostgreSQL connected successfully")
	return nil
}

// ConnectGorm establishes GORM connection for ORM operations
func (dm *DatabaseManager) ConnectGorm() error {
	var dialector gorm.Dialector

	// Use PostgreSQL in production, SQLite in development
	if dm.Config.PostgresHost != "localhost" || dm.Config.PostgresDB != "go_server" {
		dialector = postgres.Open(dm.Config.GetPostgresDSN())
	} else {
		// Use SQLite for development
		dialector = sqlite.Open("dev.db")
	}

	db, err := gorm.Open(dialector, &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}

	// Configure connection pool
	sqlDB, err := db.DB()
	if err != nil {
		return fmt.Errorf("failed to get underlying sql.DB: %w", err)
	}

	sqlDB.SetMaxOpenConns(dm.Config.MaxConnections)
	sqlDB.SetMaxIdleConns(dm.Config.MaxIdleConns)
	sqlDB.SetConnMaxLifetime(dm.Config.ConnMaxLifetime)
	sqlDB.SetConnMaxIdleTime(dm.Config.ConnMaxIdleTime)

	dm.GormDB = db
	log.Println("✅ GORM connected successfully")
	return nil
}

// ConnectRedis establishes Redis connection
func (dm *DatabaseManager) ConnectRedis(ctx context.Context) error {
	dm.RedisClient = redis.NewClient(&redis.Options{
		Addr:     dm.Config.GetRedisAddr(),
		Password: dm.Config.RedisPassword,
		DB:       dm.Config.RedisDB,
	})

	// Test connection
	if err := dm.RedisClient.Ping(ctx).Err(); err != nil {
		return fmt.Errorf("failed to connect to redis: %w", err)
	}

	log.Println("✅ Redis connected successfully")
	return nil
}

// ConnectAll establishes all database connections
func (dm *DatabaseManager) ConnectAll(ctx context.Context) error {
	// Connect to PostgreSQL
	if err := dm.ConnectPostgres(ctx); err != nil {
		return fmt.Errorf("postgres connection failed: %w", err)
	}

	// Connect to GORM
	if err := dm.ConnectGorm(); err != nil {
		return fmt.Errorf("gorm connection failed: %w", err)
	}

	// Connect to Redis
	if err := dm.ConnectRedis(ctx); err != nil {
		return fmt.Errorf("redis connection failed: %w", err)
	}

	return nil
}

// Close closes all database connections
func (dm *DatabaseManager) Close() error {
	var errs []error

	if dm.PostgresPool != nil {
		dm.PostgresPool.Close()
	}

	if dm.GormDB != nil {
		if sqlDB, err := dm.GormDB.DB(); err == nil {
			if err := sqlDB.Close(); err != nil {
				errs = append(errs, fmt.Errorf("failed to close gorm connection: %w", err))
			}
		}
	}

	if dm.RedisClient != nil {
		if err := dm.RedisClient.Close(); err != nil {
			errs = append(errs, fmt.Errorf("failed to close redis connection: %w", err))
		}
	}

	if len(errs) > 0 {
		return fmt.Errorf("errors closing connections: %v", errs)
	}

	log.Println("✅ All database connections closed")
	return nil
}

// HealthCheck performs health checks on all connections
func (dm *DatabaseManager) HealthCheck(ctx context.Context) map[string]string {
	health := make(map[string]string)

	// Check PostgreSQL
	if dm.PostgresPool != nil {
		if err := dm.PostgresPool.Ping(ctx); err != nil {
			health["postgres"] = "unhealthy: " + err.Error()
		} else {
			health["postgres"] = "healthy"
		}
	} else {
		health["postgres"] = "not connected"
	}

	// Check GORM
	if dm.GormDB != nil {
		if sqlDB, err := dm.GormDB.DB(); err == nil {
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

	// Check Redis
	if dm.RedisClient != nil {
		if err := dm.RedisClient.Ping(ctx).Err(); err != nil {
			health["redis"] = "unhealthy: " + err.Error()
		} else {
			health["redis"] = "healthy"
		}
	} else {
		health["redis"] = "not connected"
	}

	return health
}
