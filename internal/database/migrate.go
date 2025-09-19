package database

import (
	"fmt"
	"log"

	"go-server/internal/database/models"

	"gorm.io/gorm"
)

// MigrationManager handles database migrations
type MigrationManager struct {
	db     *gorm.DB
	config *DatabaseConfig
}

// NewMigrationManager creates a new migration manager
func NewMigrationManager(config *DatabaseConfig) *MigrationManager {
	return &MigrationManager{
		config: config,
	}
}

// SetupMigration initializes the migration system
func (mm *MigrationManager) SetupMigration(db *gorm.DB) error {
	mm.db = db
	return nil
}

// Up runs all pending migrations using GORM AutoMigrate
func (mm *MigrationManager) Up() error {
	if mm.db == nil {
		return fmt.Errorf("migration not initialized, call SetupMigration first")
	}

	log.Println("🔄 Running database migrations...")

	// Auto-migrate all models
	err := mm.db.AutoMigrate(
		&models.User{},
		&models.Post{},
		&models.Session{},
	)

	if err != nil {
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	log.Println("✅ Database migrations completed successfully")
	return nil
}

// Down drops all tables (use with caution!)
func (mm *MigrationManager) Down() error {
	if mm.db == nil {
		return fmt.Errorf("migration not initialized, call SetupMigration first")
	}

	log.Println("⚠️  Dropping all tables...")

	// Drop tables in reverse order to handle foreign key constraints
	err := mm.db.Migrator().DropTable(
		&models.Session{},
		&models.Post{},
		&models.User{},
	)

	if err != nil {
		return fmt.Errorf("failed to drop tables: %w", err)
	}

	log.Println("✅ All tables dropped")
	return nil
}

// Force recreates all tables
func (mm *MigrationManager) Force() error {
	if mm.db == nil {
		return fmt.Errorf("migration not initialized, call SetupMigration first")
	}

	log.Println("🔄 Force recreating all tables...")

	// Drop and recreate
	if err := mm.Down(); err != nil {
		return err
	}

	return mm.Up()
}

// Version returns migration status (simplified for GORM)
func (mm *MigrationManager) Version() (string, error) {
	if mm.db == nil {
		return "", fmt.Errorf("migration not initialized, call SetupMigration first")
	}

	// Check if tables exist
	var count int64
	err := mm.db.Raw("SELECT COUNT(*) FROM information_schema.tables WHERE table_schema = 'public'").Scan(&count).Error
	if err != nil {
		return "", fmt.Errorf("failed to check tables: %w", err)
	}

	return fmt.Sprintf("Tables: %d", count), nil
}

// Close closes the migration manager
func (mm *MigrationManager) Close() error {
	// GORM handles its own connection cleanup
	return nil
}
