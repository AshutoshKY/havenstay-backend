package migration

import (
	"log"

	"github.com/user/airbnb-test/internal/domain/booking"
	"github.com/user/airbnb-test/internal/domain/guest"
	"github.com/user/airbnb-test/internal/domain/host"
	"github.com/user/airbnb-test/internal/domain/property"
	"github.com/user/airbnb-test/internal/domain/wishlist"
	"gorm.io/gorm"
)

// RunMigrations auto-migrates all application models
func RunMigrations(db *gorm.DB) {
	log.Println("Running database migrations...")

	err := db.AutoMigrate(
		&host.Model{},
		&property.Model{},
		&guest.Model{},
		&booking.Model{},
		&wishlist.Model{},
	)

	if err != nil {
		log.Fatalf("Failed to run database migrations: %v", err)
	}

	log.Println("Database migrations completed.")
}
