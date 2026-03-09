package wishlist

import (
	"context"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	pb "github.com/user/airbnb-test/api/proto/v1"
	"github.com/user/airbnb-test/internal/pkg/response"
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

// @Summary Save Property to Wishlist
// @Description Bookmarks a specific property to a Guest's curated wishlist for later viewing or booking.
// @Description Useful for encouraging users to return to properties they found engaging.
// @Tags Wishlists
// @Accept json
// @Produce json
// @Param request body pb.AddToWishlistRequest true "Payload connecting a Guest ID to a Property ID."
// @Success 201 {object} pb.WishlistResponse "A confirmation object encapsulating the newly bookmarked item constraint."
// @Failure 400 {object} response.HTTPError "Bad Request - Duplicate bookmark or invalid identifiers."
// @Failure 500 {object} response.HTTPError "Internal Server Error."
// @Router /v1/lists [post]
func (h *HTTPServer) addToWishlist(c *gin.Context) {
	var req pb.AddToWishlistRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.HTTPError{Code: http.StatusBadRequest, Message: "Validation Bad Request", Details: err.Error()})
		return
	}

	resp, err := h.svc.AddToWishlist(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.HTTPError{Code: http.StatusInternalServerError, Message: "Internal Server Error", Details: err.Error()})
		return
	}
	c.JSON(http.StatusCreated, resp)
}

// @Summary View Guest Wishlist
// @Description Fetches the complete array of bookmarked properties saved under a specific Guest account.
// @Description Acts as the primary data source for the "Saved Homes" UI tab.
// @Tags Wishlists
// @Produce json
// @Param guestId query int true "The target Guest's unique Identifier."
// @Success 200 {object} pb.GetWishlistResponse "An array wrapping the nested Property details saved by the user."
// @Failure 400 {object} response.HTTPError "Bad Request."
// @Failure 500 {object} response.HTTPError "Internal Server Error."
// @Router /v1/lists [get]
func (h *HTTPServer) getWishlist(c *gin.Context) {
	guestIdStr := c.Query("guestId")
	guestId, err := strconv.ParseInt(guestIdStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.HTTPError{Code: http.StatusBadRequest, Message: "Missing or Invalid guestId query param", Details: err.Error()})
		return
	}

	resp, err := h.svc.GetWishlist(c.Request.Context(), &pb.GetWishlistRequest{GuestId: guestId})
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.HTTPError{Code: http.StatusInternalServerError, Message: "Internal Server Error", Details: err.Error()})
		return
	}
	c.JSON(http.StatusOK, resp)
}
