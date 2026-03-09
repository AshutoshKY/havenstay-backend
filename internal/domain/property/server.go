package property

import (
	"context"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	pb "github.com/user/airbnb-test/api/proto/v1"
	"github.com/user/airbnb-test/internal/pkg/response"
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

// @Summary Browse all Properties
// @Description Retrieve a sweeping list of all active property listings on the HavenStay platform.
// @Description *(Note: In future versions, this will include global pagination and bounding-box map filters).*
// @Tags Properties
// @Produce json
// @Success 200 {object} pb.ListPropertiesResponse "A comprehensive array of all available properties."
// @Failure 500 {object} response.HTTPError "Internal Server Error while querying the database."
// @Router /v1/properties [get]
func (h *HTTPServer) listProperties(c *gin.Context) {
	resp, err := h.svc.ListProperties(c.Request.Context(), &pb.ListPropertiesRequest{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.HTTPError{Code: http.StatusInternalServerError, Message: "Internal Server Error", Details: err.Error()})
		return
	}
	c.JSON(http.StatusOK, resp)
}

// @Summary Get Property Details
// @Description Drill down into a specific property. Returns the exact title, host bindings, textual description, nightly price, and location coordinates formatted for front-end rendering.
// @Tags Properties
// @Produce json
// @Param id path int true "The exact Property ID to inspect."
// @Success 200 {object} pb.PropertyResponse "Full property overview returned."
// @Failure 400 {object} response.HTTPError "Bad Request."
// @Failure 404 {object} response.HTTPError "Property Not Found - the listing may have been deactivated."
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

// @Summary List a new Property
// @Description An authenticated Host can use this endpoint to publish a new real-estate listing.
// @Description Note: You MUST pass a valid `host_id` representing an already-registered Host on the platform, or the database foreign-key constraint will reject the listing.
// @Tags Properties
// @Accept json
// @Produce json
// @Param request body pb.CreatePropertyRequest true "The core details of the listing: Name, Description, Location String, and Nightly Price (Float)."
// @Success 201 {object} pb.PropertyResponse "Property was published and is now live on the index."
// @Failure 400 {object} response.HTTPError "Validation Bad Request - Usually missing Host ID or malformed price."
// @Failure 500 {object} response.HTTPError "Internal Server Error."
// @Router /v1/properties [post]
func (h *HTTPServer) createProperty(c *gin.Context) {
	var req pb.CreatePropertyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.HTTPError{Code: http.StatusBadRequest, Message: "Validation Error", Details: err.Error()})
		return
	}

	resp, err := h.svc.CreateProperty(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.HTTPError{Code: http.StatusInternalServerError, Message: "Internal Server Error", Details: err.Error()})
		return
	}
	c.JSON(http.StatusCreated, resp)
}

// @Summary Host's Dashboard Properties
// @Description Get all properties listed specifically by one Host.
// @Description This endpoint powers the "Host Dashboard", allowing property owners to see all assets they have under management. Optionally, you can pass a `location` query to filter the host's properties by a specific city.
// @Tags Properties
// @Produce json
// @Param hostId path int true "The Host ID managing these properties."
// @Param location query string false "Optional string to filter the resulting list by a specific city or region."
// @Success 200 {object} pb.ListPropertiesResponse "A list of properties bound to the specified host."
// @Failure 400 {object} response.HTTPError "Bad Request."
// @Failure 500 {object} response.HTTPError "Internal Server Error."
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
