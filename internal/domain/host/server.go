package host

import (
	"context"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	pb "github.com/user/airbnb-test/api/proto/v1"
)

// GRPCServer implements the gRPC HostServiceServer interface
type GRPCServer struct {
	pb.UnimplementedHostServiceServer
	svc Service
}

// NewGRPCServer creates a new gRPC server for the Host domain
func NewGRPCServer(svc Service) *GRPCServer {
	return &GRPCServer{svc: svc}
}

func (s *GRPCServer) CreateHost(ctx context.Context, req *pb.CreateHostRequest) (*pb.HostResponse, error) {
	return s.svc.CreateHost(ctx, req)
}

func (s *GRPCServer) GetHost(ctx context.Context, req *pb.GetHostRequest) (*pb.HostResponse, error) {
	return s.svc.GetHost(ctx, req)
}

// HTTPServer handles HTTP REST requests for Host domain
type HTTPServer struct {
	svc Service
}

// NewHTTPServer creates HTTP handlers
func NewHTTPServer(svc Service) *HTTPServer {
	return &HTTPServer{svc: svc}
}

// RegisterRoutes registers the REST endpoints with Gin
func (h *HTTPServer) RegisterRoutes(router *gin.RouterGroup) {
	hostGroup := router.Group("/hosts") // using /hosts as standard
	{
		hostGroup.POST("", h.createHost)
		hostGroup.GET("/:id", h.getHost)
	}
}

// @Summary Register a new Host
// @Description Creates a new Host profile on the platform. Hosts are independent users capable of listing and managing multiple Properties.
// @Description Note: A newly created host is automatically marked as `is_verified: false` pending administrator approval.
// @Tags Hosts
// @Accept json
// @Produce json
// @Param request body pb.CreateHostRequest true "Payload containing the Host's basic contact and profile information (Name, Email, Phone, Location)."
// @Success 201 {object} pb.HostResponse "Host successfully registered and assigned a unique ID."
// @Failure 400 {object} string "Validation Error - Invalid email format or missing required fields."
// @Failure 500 {object} string "Internal Server Error - Database connection failure."
// @Router /v1/hosts [post]
func (h *HTTPServer) createHost(c *gin.Context) {
	var req pb.CreateHostRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp, err := h.svc.CreateHost(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, resp)
}

// @Summary Retrieve Host details
// @Description Fetches the complete public profile of a Host by their unique identifier natively from the core system.
// @Description Use this endpoint to display a Host's verification status, contact info, and join date on their public landing page.
// @Tags Hosts
// @Produce json
// @Param id path int true "The unique numeric ID of the Host to retrieve."
// @Success 200 {object} pb.HostResponse "The requested Host details were found and returned successfully."
// @Failure 400 {object} string "Bad Request - The provided ID was not a valid integer."
// @Failure 404 {object} string "Not Found - No Host exists with the provided ID."
// @Router /v1/hosts/{id} [get]
func (h *HTTPServer) getHost(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid host ID"})
		return
	}

	resp, err := h.svc.GetHost(c.Request.Context(), &pb.GetHostRequest{Id: id})
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}
