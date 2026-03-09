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

// @Summary Create a Host
// @Description Register a new host on the platform
// @Tags hosts
// @Accept json
// @Produce json
// @Param request body pb.CreateHostRequest true "Host creation payload"
// @Success 201 {object} pb.HostResponse "Successfully created host"
// @Failure 400 {object} string "Bad Request"
// @Failure 500 {object} string "Internal Server Error"
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

// @Summary Get a Host
// @Description Get details of a single host by ID
// @Tags hosts
// @Produce json
// @Param id path int true "Host ID"
// @Success 200 {object} pb.HostResponse "Successfully retrieved host"
// @Failure 400 {object} string "Bad Request"
// @Failure 404 {object} string "Host Not Found"
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
