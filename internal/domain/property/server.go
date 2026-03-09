package property

import (
	"context"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	pb "github.com/user/airbnb-test/api/proto/v1"
)

// GRPCServer implements pb.PropertyServiceServer
type GRPCServer struct {
	pb.UnimplementedPropertyServiceServer
	svc Service
}

func NewGRPCServer(svc Service) *GRPCServer {
	return &GRPCServer{svc: svc}
}

func (s *GRPCServer) GetProperty(ctx context.Context, req *pb.GetPropertyRequest) (*pb.PropertyResponse, error) {
	return s.svc.GetProperty(ctx, req)
}

func (s *GRPCServer) ListPropertiesByHost(ctx context.Context, req *pb.ListPropertiesByHostRequest) (*pb.ListPropertiesResponse, error) {
	return s.svc.ListPropertiesByHost(ctx, req)
}

func (s *GRPCServer) CreateProperty(ctx context.Context, req *pb.CreatePropertyRequest) (*pb.PropertyResponse, error) {
	return s.svc.CreateProperty(ctx, req)
}

func (s *GRPCServer) ListProperties(ctx context.Context, req *pb.ListPropertiesRequest) (*pb.ListPropertiesResponse, error) {
	return s.svc.ListProperties(ctx, req)
}

// HTTPServer handles HTTP REST requests for Property domain
type HTTPServer struct {
	svc Service
}

func NewHTTPServer(svc Service) *HTTPServer {
	return &HTTPServer{svc: svc}
}

func (h *HTTPServer) RegisterRoutes(router *gin.RouterGroup) {
	propertyGroup := router.Group("/properties")
	{
		propertyGroup.GET("", h.listProperties)
		propertyGroup.POST("", h.createProperty)
		propertyGroup.GET("/:id", h.getProperty)
		propertyGroup.GET("/host/:hostId", h.listPropertiesByHost)
	}
}

// @Summary List properties
// @Description Get a list of all properties
// @Tags properties
// @Produce json
// @Success 200 {object} pb.ListPropertiesResponse "Successfully retrieved list of properties"
// @Failure 500 {object} string "Internal Server Error"
// @Router /v1/properties [get]
func (h *HTTPServer) listProperties(c *gin.Context) {
	resp, err := h.svc.ListProperties(c.Request.Context(), &pb.ListPropertiesRequest{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, resp)
}

// @Summary Get Property Details
// @Description Get details of a single property by its ID
// @Tags properties
// @Produce json
// @Param id path int true "Property ID"
// @Success 200 {object} pb.PropertyResponse "Successfully retrieved property"
// @Failure 400 {object} string "Bad Request"
// @Failure 404 {object} string "Property Not Found"
// @Router /v1/properties/{id} [get]
func (h *HTTPServer) getProperty(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid property ID"})
		return
	}

	resp, err := h.svc.GetProperty(c.Request.Context(), &pb.GetPropertyRequest{PropertyId: id})
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, resp)
}

// @Summary Create a Property
// @Description Create a new property listing
// @Tags properties
// @Accept json
// @Produce json
// @Param request body pb.CreatePropertyRequest true "Property creation payload"
// @Success 201 {object} pb.PropertyResponse "Successfully created property"
// @Failure 400 {object} string "Bad Request"
// @Failure 500 {object} string "Internal Server Error"
// @Router /v1/properties [post]
func (h *HTTPServer) createProperty(c *gin.Context) {
	var req pb.CreatePropertyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp, err := h.svc.CreateProperty(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, resp)
}

// @Summary List Properties by Host
// @Description Get all properties listed by a specific host, with optional location filtering
// @Tags properties
// @Produce json
// @Param hostId path int true "Host ID"
// @Param location query string false "Location filter"
// @Success 200 {object} pb.ListPropertiesResponse "Successfully retrieved host properties"
// @Failure 400 {object} string "Bad Request"
// @Failure 500 {object} string "Internal Server Error"
// @Router /v1/properties/host/{hostId} [get]
func (h *HTTPServer) listPropertiesByHost(c *gin.Context) {
	hostIdStr := c.Param("hostId")
	hostId, err := strconv.ParseInt(hostIdStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid host ID"})
		return
	}

	locationQuery := c.Query("location")

	resp, err := h.svc.ListPropertiesByHost(c.Request.Context(), &pb.ListPropertiesByHostRequest{
		HostId:   hostId,
		Location: locationQuery,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, resp)
}
