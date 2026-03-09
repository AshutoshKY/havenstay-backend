package wishlist

import (
	"context"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	pb "github.com/user/airbnb-test/api/proto/v1"
)

type GRPCServer struct {
	pb.UnimplementedWishlistServiceServer
	svc Service
}

func NewGRPCServer(svc Service) *GRPCServer {
	return &GRPCServer{svc: svc}
}

func (s *GRPCServer) AddToWishlist(ctx context.Context, req *pb.AddToWishlistRequest) (*pb.WishlistResponse, error) {
	return s.svc.AddToWishlist(ctx, req)
}

func (s *GRPCServer) GetWishlist(ctx context.Context, req *pb.GetWishlistRequest) (*pb.GetWishlistResponse, error) {
	return s.svc.GetWishlist(ctx, req)
}

type HTTPServer struct {
	svc Service
}

func NewHTTPServer(svc Service) *HTTPServer {
	return &HTTPServer{svc: svc}
}

// RegisterRoutes handles endpoints like /lists instead of /wishlist to match the specific list requirement
func (h *HTTPServer) RegisterRoutes(router *gin.RouterGroup) {
	group := router.Group("/lists")
	{
		group.GET("", h.getWishlist) // /lists?guestId=123
		group.POST("", h.addToWishlist)
	}
}

// @Summary Add to Wishlist
// @Description Add a property to a guest's wishlist
// @Tags wishlists
// @Accept json
// @Produce json
// @Param request body pb.AddToWishlistRequest true "Wishlist addition payload"
// @Success 201 {object} pb.WishlistResponse "Successfully added to wishlist"
// @Failure 400 {object} string "Bad Request"
// @Failure 500 {object} string "Internal Server Error"
// @Router /v1/lists [post]
func (h *HTTPServer) addToWishlist(c *gin.Context) {
	var req pb.AddToWishlistRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp, err := h.svc.AddToWishlist(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, resp)
}

// @Summary Get Wishlist
// @Description Retrieve the wishlist of a specific guest
// @Tags wishlists
// @Produce json
// @Param guestId query int true "Guest ID"
// @Success 200 {object} pb.GetWishlistResponse "Successfully retrieved guest's wishlist"
// @Failure 400 {object} string "Bad Request"
// @Failure 500 {object} string "Internal Server Error"
// @Router /v1/lists [get]
func (h *HTTPServer) getWishlist(c *gin.Context) {
	guestIdStr := c.Query("guestId")
	guestId, err := strconv.ParseInt(guestIdStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "guestId query param is required"})
		return
	}

	resp, err := h.svc.GetWishlist(c.Request.Context(), &pb.GetWishlistRequest{GuestId: guestId})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, resp)
}
