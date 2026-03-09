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

// @Summary Create a Booking Reservation
// @Description Securely initiates a booking transaction for a specific Property by a Guest.
// @Description The system will evaluate the requested dates, calculate total pricing, and attempt to reserve the asset.
// @Tags Bookings
// @Accept json
// @Produce json
// @Param request body pb.CreateBookingRequest true "Payload specifying the Guest ID, Property ID, Start Date, End Date, and calculated total price."
// @Success 201 {object} pb.BookingResponse "Reservation successfully created and locked."
// @Failure 400 {object} string "Validation Bad Request - Overlapping dates or missing required fields."
// @Failure 500 {object} string "Internal Server Error."
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

// @Summary Retrieve Booking Invoice
// @Description Fetch an exact copy of a completed or pending Booking reservation by its unique confirmation ID.
// @Description Typical use-case: Displaying a digital receipt or itinerary details to the Guest.
// @Tags Bookings
// @Produce json
// @Param id path int true "The numeric Reservation/Booking ID."
// @Success 200 {object} pb.BookingResponse "The retrieved itinerary and invoice details."
// @Failure 400 {object} string "Bad Request - Invalid ID format."
// @Failure 404 {object} string "Not Found - The reservation does not exist."
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

// @Summary Get a Guest's Booking History
// @Description Poll the history of all active, past, and cancelled reservations tied to a specific Guest profile.
// @Tags Bookings
// @Produce json
// @Param guestId query int true "The ID of the Guest requesting their itinerary."
// @Success 200 {object} pb.ListBookingsResponse "An array of all historical and upcoming bookings for the Guest."
// @Failure 400 {object} string "Bad Request - Missing or invalid guest ID."
// @Failure 500 {object} string "Internal Server Error."
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
