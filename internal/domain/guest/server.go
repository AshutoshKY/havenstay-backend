package guest

import (
	"context"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	pb "github.com/user/airbnb-test/api/proto/v1"
)

type GRPCServer struct {
	pb.UnimplementedGuestServiceServer
	svc Service
}

func NewGRPCServer(svc Service) *GRPCServer {
	return &GRPCServer{svc: svc}
}

func (s *GRPCServer) CreateGuest(ctx context.Context, req *pb.CreateGuestRequest) (*pb.GuestResponse, error) {
	return s.svc.CreateGuest(ctx, req)
}

func (s *GRPCServer) GetGuest(ctx context.Context, req *pb.GetGuestRequest) (*pb.GuestResponse, error) {
	return s.svc.GetGuest(ctx, req)
}

type HTTPServer struct {
	svc Service
}

func NewHTTPServer(svc Service) *HTTPServer {
	return &HTTPServer{svc: svc}
}

func (h *HTTPServer) RegisterRoutes(router *gin.RouterGroup) {
	group := router.Group("/guests")
	{
		group.POST("", h.createGuest)
		group.GET("/:id", h.getGuest)
	}
}

// @Summary Create a Guest
// @Description Register a new guest on the platform
// @Tags guests
// @Accept json
// @Produce json
// @Param request body pb.CreateGuestRequest true "Guest creation payload"
// @Success 201 {object} pb.GuestResponse "Successfully created guest"
// @Failure 400 {object} string "Bad Request"
// @Failure 500 {object} string "Internal Server Error"
// @Router /v1/guests [post]
func (h *HTTPServer) createGuest(c *gin.Context) {
	var req pb.CreateGuestRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp, err := h.svc.CreateGuest(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, resp)
}

// @Summary Get a Guest
// @Description Get details of a single guest by ID
// @Tags guests
// @Produce json
// @Param id path int true "Guest ID"
// @Success 200 {object} pb.GuestResponse "Successfully retrieved guest"
// @Failure 400 {object} string "Bad Request"
// @Failure 404 {object} string "Guest Not Found"
// @Router /v1/guests/{id} [get]
func (h *HTTPServer) getGuest(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid guest ID"})
		return
	}

	resp, err := h.svc.GetGuest(c.Request.Context(), &pb.GetGuestRequest{GuestId: id})
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, resp)
}
