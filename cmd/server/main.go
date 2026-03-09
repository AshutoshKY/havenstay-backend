package main

import (
	"fmt"
	"log"
	"net"

	"github.com/user/airbnb-test/config"
	"github.com/user/airbnb-test/internal/domain/booking"
	"github.com/user/airbnb-test/internal/domain/guest"
	"github.com/user/airbnb-test/internal/domain/host"
	"github.com/user/airbnb-test/internal/domain/property"
	"github.com/user/airbnb-test/internal/domain/wishlist"
	"github.com/user/airbnb-test/internal/platform/database"
	"github.com/user/airbnb-test/internal/platform/migration"
	"github.com/user/airbnb-test/internal/platform/server"

	_ "github.com/user/airbnb-test/docs" // Required for Swaggo
)

// @title HavenStay - Hotel & Room Booking API
// @version 1.0.0
// @description Welcome to the HavenStay Backend API documentation! 🚀
// @description
// @description HavenStay is a modern, high-performance Domain-Driven Design (DDD) backend built in Go. It empowers Hosts to manage luxury properties and Guests to seamlessly book remote getaways.
// @description
// @description ### Key Features:
// @description - **🏘️ Property Management**: Full CRUD lifecycle for short-term rentals and luxury homes.
// @description - **👩‍💼 Host & Guest Profiles**: Independent onboarding and profile management domains.
// @description - **📅 Booking Engine**: Real-time reservation creation.
// @description - **❤️ Wishlists**: Guests can curate lists of favorite properties.
// @description
// @description **Built With**: Go 1.21+, Gin (REST), gRPC, Protobufs, MySQL 8, GORM.
// @description
// @description 📥 **[Download Postman Collection](https://raw.githubusercontent.com/AshutoshKY/havenstay-backend/master/postman_collection.json)**
// @termsOfService https://havenstay.localhost/terms

// @contact.name HavenStay Platform Support
// @contact.url https://github.com/AshutoshKY/havenstay-backend
// @contact.email api-support@havenstay.local

// @license.name MIT License
// @license.url https://opensource.org/licenses/MIT

// @host localhost:8080
// @BasePath /

func main() {
	// 1. Load config
	cfg := config.LoadConfig()

	// 2. Setup Database & run migrations
	db, err := database.Connect(cfg)
	if err != nil {
		log.Fatalf("Database connection failed: %v", err)
	}
	migration.RunMigrations(db)

	// 3. Initialize Repositories
	hostRepo := host.NewMySQLRepository(db)
	propertyRepo := property.NewMySQLRepository(db)
	guestRepo := guest.NewMySQLRepository(db)
	bookingRepo := booking.NewMySQLRepository(db)
	wishlistRepo := wishlist.NewMySQLRepository(db)

	// 4. Initialize Core Services
	hostSvc := host.NewService(hostRepo)
	propertySvc := property.NewService(propertyRepo)
	guestSvc := guest.NewService(guestRepo)
	bookingSvc := booking.NewService(bookingRepo)
	wishlistSvc := wishlist.NewService(wishlistRepo)

	// 5. Initialize Domain Servers (gRPC & HTTP handlers)
	hostGRPC := host.NewGRPCServer(hostSvc)
	hostHTTP := host.NewHTTPServer(hostSvc)

	propertyGRPC := property.NewGRPCServer(propertySvc)
	propertyHTTP := property.NewHTTPServer(propertySvc)

	guestGRPC := guest.NewGRPCServer(guestSvc)
	guestHTTP := guest.NewHTTPServer(guestSvc)

	bookingGRPC := booking.NewGRPCServer(bookingSvc)
	bookingHTTP := booking.NewHTTPServer(bookingSvc)

	wishlistGRPC := wishlist.NewGRPCServer(wishlistSvc)
	wishlistHTTP := wishlist.NewHTTPServer(wishlistSvc)

	// 6. Setup Root Servers
	grpcServer := server.NewGRPCServer(
		hostGRPC,
		propertyGRPC,
		guestGRPC,
		bookingGRPC,
		wishlistGRPC,
	)

	httpServer := server.NewHTTPServer(
		hostHTTP,
		propertyHTTP,
		guestHTTP,
		bookingHTTP,
		wishlistHTTP,
	)

	// 7. Start gRPC server in a goroutine
	go func() {
		lis, err := net.Listen("tcp", fmt.Sprintf(":%s", cfg.GRPCPort))
		if err != nil {
			log.Fatalf("failed to listen for gRPC: %v", err)
		}
		log.Printf("Starting gRPC server on port %s", cfg.GRPCPort)
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("failed to serve gRPC: %v", err)
		}
	}()

	// 8. Start HTTP server (blocking)
	log.Printf("Starting HTTP server on port %s", cfg.HTTPPort)
	if err := httpServer.Run(fmt.Sprintf(":%s", cfg.HTTPPort)); err != nil {
		log.Fatalf("failed to serve HTTP: %v", err)
	}
}
