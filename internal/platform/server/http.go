package server

import (
	"github.com/gin-gonic/gin"
	"github.com/user/airbnb-test/internal/domain/booking"
	"github.com/user/airbnb-test/internal/domain/guest"
	"github.com/user/airbnb-test/internal/domain/host"
	"github.com/user/airbnb-test/internal/domain/property"
	"github.com/user/airbnb-test/internal/domain/wishlist"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// NewHTTPServer initializes Gin router with all domain routes
func NewHTTPServer(
	hostHTTP *host.HTTPServer,
	propertyHTTP *property.HTTPServer,
	guestHTTP *guest.HTTPServer,
	bookingHTTP *booking.HTTPServer,
	wishlistHTTP *wishlist.HTTPServer,
) *gin.Engine {

	r := gin.Default()

	// v1 API Group
	v1 := r.Group("/v1")
	{
		hostHTTP.RegisterRoutes(v1)
		propertyHTTP.RegisterRoutes(v1)
		guestHTTP.RegisterRoutes(v1)
		bookingHTTP.RegisterRoutes(v1)
		wishlistHTTP.RegisterRoutes(v1)
	}

	// Swagger documentation route
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	return r
}
