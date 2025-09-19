package main

import (
	"context"
	"log"
	"time"

	"go-server/internal/database"
	"go-server/internal/database/models"
	"go-server/internal/database/repositories"
)

func main() {
	log.Println("🔍 Testing Database Integration...")

	// Load database configuration
	dbConfig := database.NewDatabaseConfig()
	log.Printf("📋 Database Config: PostgreSQL=%s:%d/%s, Redis=%s:%d", 
		dbConfig.PostgresHost, dbConfig.PostgresPort, dbConfig.PostgresDB,
		dbConfig.RedisHost, dbConfig.RedisPort)

	// Create database manager
	dbManager := database.NewDatabaseManager(dbConfig)

	// Test database connections
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	log.Println("🔌 Connecting to databases...")
	if err := dbManager.ConnectAll(ctx); err != nil {
		log.Printf("❌ Database connection failed: %v", err)
		log.Println("💡 This is expected if databases are not running")
		log.Println("💡 To test with databases, start PostgreSQL and Redis")
		return
	}

	log.Println("✅ Successfully connected to all databases!")

	// Test database operations
	log.Println("🧪 Testing database operations...")

	// Test PostgreSQL connection
	if dbManager.PostgresPool != nil {
		log.Println("✅ PostgreSQL connection pool is active")
	} else {
		log.Println("❌ PostgreSQL connection pool is nil")
	}

	// Test GORM connection
	if dbManager.GormDB != nil {
		log.Println("✅ GORM database connection is active")
	} else {
		log.Println("❌ GORM database connection is nil")
	}

	// Test Redis connection
	if dbManager.RedisClient != nil {
		log.Println("✅ Redis client is active")
	} else {
		log.Println("❌ Redis client is nil")
	}

	// Test migrations
	log.Println("🔄 Testing database migrations...")
	// Note: Migrations would be handled by the migrate package
	// For now, we'll skip this test
	log.Println("⏭️ Skipping migrations test (requires migrate package)")

	// Test basic operations
	log.Println("🧪 Testing basic database operations...")

	// Test user creation
	userRepo := repositories.NewUserRepository(dbManager.GormDB)
	testUser := &models.User{
		Username:  "testuser",
		Email:     "test@example.com",
		FirstName: "Test",
		LastName:  "User",
		IsActive:  true,
	}

	if err := userRepo.CreateUser(ctx, testUser); err != nil {
		log.Printf("❌ User creation failed: %v", err)
	} else {
		log.Println("✅ User created successfully")
	}

	// Test user retrieval
	retrievedUser, err := userRepo.GetUserByEmail(ctx, "test@example.com")
	if err != nil {
		log.Printf("❌ User retrieval failed: %v", err)
	} else {
		log.Printf("✅ User retrieved: %s (%s)", retrievedUser.Username, retrievedUser.Email)
	}

	// Test Redis operations
	cacheRepo := repositories.NewCacheRepository(dbManager.RedisClient)
	if err := cacheRepo.Set(ctx, "test:key", "test:value", 5*time.Minute); err != nil {
		log.Printf("❌ Redis set failed: %v", err)
	} else {
		log.Println("✅ Redis set operation successful")
	}

	value, err := cacheRepo.Get(ctx, "test:key")
	if err != nil {
		log.Printf("❌ Redis get failed: %v", err)
	} else {
		log.Printf("✅ Redis get operation successful: %s", value)
	}

	// Cleanup
	log.Println("🧹 Cleaning up test data...")
	if err := userRepo.DeleteUser(ctx, testUser.ID); err != nil {
		log.Printf("⚠️ User cleanup failed: %v", err)
	} else {
		log.Println("✅ Test user cleaned up")
	}

	if err := cacheRepo.Delete(ctx, "test:key"); err != nil {
		log.Printf("⚠️ Redis cleanup failed: %v", err)
	} else {
		log.Println("✅ Redis test data cleaned up")
	}

	// Close connections
	log.Println("🔌 Closing database connections...")
	if err := dbManager.Close(); err != nil {
		log.Printf("❌ Error closing connections: %v", err)
	} else {
		log.Println("✅ All database connections closed successfully")
	}

	log.Println("🎉 Database integration test completed!")
}
