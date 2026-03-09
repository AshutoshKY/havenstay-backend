package server

import (
	pb "github.com/user/airbnb-test/api/proto/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

// NewGRPCServer initializes the grpc Server with all domain services
func NewGRPCServer(
	hostServer pb.HostServiceServer,
	propertyServer pb.PropertyServiceServer,
	guestServer pb.GuestServiceServer,
	bookingServer pb.BookingServiceServer,
	wishlistServer pb.WishlistServiceServer,
) *grpc.Server {

	s := grpc.NewServer()

	pb.RegisterHostServiceServer(s, hostServer)
	pb.RegisterPropertyServiceServer(s, propertyServer)
	pb.RegisterGuestServiceServer(s, guestServer)
	pb.RegisterBookingServiceServer(s, bookingServer)
	pb.RegisterWishlistServiceServer(s, wishlistServer)

	// Enable reflection for tools like grpcurl
	reflection.Register(s)

	return s
}
