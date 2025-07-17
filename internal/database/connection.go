package database

import (
	"fmt"

	"airbnb-clone/internal/config"
	"airbnb-clone/internal/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// creates a new database connection
func NewConnection(cfg config.DatabaseConfig) (*gorm.DB, error) {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=UTC",
		cfg.Host, cfg.User, cfg.Password, cfg.DBName, cfg.Port, cfg.SSLMode)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// get underlying sql.DB
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get underlying sql.DB: %w", err)
	}

	// set connection pool settings
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)

	return db, nil
}

// runs database migrations
func Migrate(db *gorm.DB) error {
	// enable UUID extension for PostgreSQL
	err := db.Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\"").Error
	if err != nil {
		return fmt.Errorf("failed to create uuid extension: %w", err)
	}

	err = db.AutoMigrate(
		&models.User{},
		&models.Property{},
		&models.Booking{},
		&models.Review{},
	)
	if err != nil {
		return fmt.Errorf("failed to migrate database: %w", err)
	}

	err = createIndexes(db)
	if err != nil {
		return fmt.Errorf("failed to create indexes: %w", err)
	}

	return nil
}

// creates additional indexes for better performance
func createIndexes(db *gorm.DB) error {
	// property indexes
	indexes := []string{
		"CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_properties_location ON properties (city, state, country)",
		"CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_properties_type ON properties (type)",
		"CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_properties_status ON properties (status)",
		"CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_properties_host_id ON properties (host_id)",
		
		// booking indexes
		"CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_bookings_dates ON bookings (check_in, check_out)",
		"CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_bookings_property_id ON bookings (property_id)",
		"CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_bookings_guest_id ON bookings (guest_id)",
		"CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_bookings_status ON bookings (status)",
		
		// review indexes
		"CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_reviews_property_id ON reviews (property_id)",
		"CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_reviews_reviewer_id ON reviews (reviewer_id)",
		"CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_reviews_rating ON reviews (rating)",
		
		// user indexes
		"CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_users_email ON users (email)",
		"CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_users_role ON users (role)",
	}

	for _, index := range indexes {
		if err := db.Exec(index).Error; err != nil {
			return fmt.Errorf("failed to create index: %w", err)
		}
	}

	return nil
}