package booking

import (
	"context"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	pb "github.com/user/airbnb-test/api/proto/v1"
)

type GRPCServer struct {
	pb.UnimplementedBookingServiceServer
	svc Service
}

func NewGRPCServer(svc Service) *GRPCServer {
	return &GRPCServer{svc: svc}
}

func (s *GRPCServer) CreateBooking(ctx context.Context, req *pb.CreateBookingRequest) (*pb.BookingResponse, error) {
	return s.svc.CreateBooking(ctx, req)
}

func (s *GRPCServer) GetBooking(ctx context.Context, req *pb.GetBookingRequest) (*pb.BookingResponse, error) {
	return s.svc.GetBooking(ctx, req)
}

func (s *GRPCServer) ListGuestBookings(ctx context.Context, req *pb.ListGuestBookingsRequest) (*pb.ListBookingsResponse, error) {
	return s.svc.ListGuestBookings(ctx, req)
}

type HTTPServer struct {
	svc Service
}

func NewHTTPServer(svc Service) *HTTPServer {
	return &HTTPServer{svc: svc}
}

func (h *HTTPServer) RegisterRoutes(router *gin.RouterGroup) {
	group := router.Group("/bookings")
	{
		group.GET("", h.listGuestBookings) // e.g. /bookings?guestId=123
		group.POST("", h.createBooking)    // create booking
		group.GET("/:id", h.getBooking)    // get details
	}
}

// @Summary Create a Booking
// @Description Make a booking for a property
// @Tags bookings
// @Accept json
// @Produce json
// @Param request body pb.CreateBookingRequest true "Booking creation payload"
// @Success 201 {object} pb.BookingResponse "Successfully created booking"
// @Failure 400 {object} string "Bad Request"
// @Failure 500 {object} string "Internal Server Error"
// @Router /v1/bookings [post]
func (h *HTTPServer) createBooking(c *gin.Context) {
	var req pb.CreateBookingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp, err := h.svc.CreateBooking(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, resp)
}

// @Summary Get a Booking
// @Description Get details of a single booking by ID
// @Tags bookings
// @Produce json
// @Param id path int true "Booking ID"
// @Success 200 {object} pb.BookingResponse "Successfully retrieved booking"
// @Failure 400 {object} string "Bad Request"
// @Failure 404 {object} string "Booking Not Found"
// @Router /v1/bookings/{id} [get]
func (h *HTTPServer) getBooking(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid booking ID"})
		return
	}

	resp, err := h.svc.GetBooking(c.Request.Context(), &pb.GetBookingRequest{BookingId: id})
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, resp)
}

// @Summary List Guest Bookings
// @Description Get all bookings made by a specific guest
// @Tags bookings
// @Produce json
// @Param guestId query int true "Guest ID"
// @Success 200 {object} pb.ListBookingsResponse "Successfully retrieved guest bookings"
// @Failure 400 {object} string "Bad Request"
// @Failure 500 {object} string "Internal Server Error"
// @Router /v1/bookings [get]
func (h *HTTPServer) listGuestBookings(c *gin.Context) {
	guestIdStr := c.Query("guestId")
	guestId, err := strconv.ParseInt(guestIdStr, 10, 64)
	if err != nil {
		// Could also list ALL bookings if no guestId is provided, based on requirement: "All the bookings user has done"
		c.JSON(http.StatusBadRequest, gin.H{"error": "guestId query parameter is required"})
		return
	}

	resp, err := h.svc.ListGuestBookings(c.Request.Context(), &pb.ListGuestBookingsRequest{GuestId: guestId})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, resp)
}
