package main

import (
	"log"
	"time"

	"github.com/user/airbnb-test/config"
	"github.com/user/airbnb-test/internal/domain/booking"
	"github.com/user/airbnb-test/internal/domain/guest"
	"github.com/user/airbnb-test/internal/domain/host"
	"github.com/user/airbnb-test/internal/domain/property"
	"github.com/user/airbnb-test/internal/platform/database"
	"github.com/user/airbnb-test/internal/platform/migration"
)

func main() {
	cfg := config.LoadConfig()

	db, err := database.Connect(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// 1. Run migrations first
	migration.RunMigrations(db)

	log.Println("Starting database seed...")

	// 2. Clear existing (optional, but good for idempotent seeds)
	db.Exec("SET FOREIGN_KEY_CHECKS = 0")
	db.Exec("TRUNCATE TABLE bookings")
	db.Exec("TRUNCATE TABLE properties")
	db.Exec("TRUNCATE TABLE guests")
	db.Exec("TRUNCATE TABLE hosts")
	db.Exec("SET FOREIGN_KEY_CHECKS = 1")

	// 3. Insert Hosts
	h1 := host.Model{
		UserName:   "johndoe",
		Email:      "john@example.com",
		Phone:      "+1-555-0101",
		Location:   "New York, USA",
		IsVerified: true,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}
	db.Create(&h1)

	// 4. Insert Properties
	p1 := property.Model{
		HostID:      h1.ID,
		Name:        "Cozy Central Park Apartment",
		Location:    "New York, USA",
		Description: "A beautiful 2-bedroom near the park.",
		Price:       150,
		CreatedAt:   time.Now(),
	}
	db.Create(&p1)

	// 5. Insert Guests
	g1 := guest.Model{
		Name:      "Alice Smith",
		Email:     "alice@example.com",
		Phone:     "+44-7700-900000",
		Country:   "UK",
		CreatedAt: time.Now(),
	}
	db.Create(&g1)

	// 6. Insert Bookings
	b1 := booking.Model{
		PropertyID:  p1.ID,
		GuestID:     g1.GuestID,
		StartDate:   "2026-04-01",
		EndDate:     "2026-04-07",
		Status:      "CONFIRMED",
		PaymentInfo: "CREDIT_CARD",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	db.Create(&b1)

	log.Println("Database seed completed successfully!")
}
