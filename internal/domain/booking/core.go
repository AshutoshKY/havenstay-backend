package booking

import (
	"context"

	pb "github.com/user/airbnb-test/api/proto/v1"
)

type Service interface {
	CreateBooking(ctx context.Context, req *pb.CreateBookingRequest) (*pb.BookingResponse, error)
	GetBooking(ctx context.Context, req *pb.GetBookingRequest) (*pb.BookingResponse, error)
	ListGuestBookings(ctx context.Context, req *pb.ListGuestBookingsRequest) (*pb.ListBookingsResponse, error)
}

type coreService struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &coreService{repo: repo}
}

func (s *coreService) CreateBooking(ctx context.Context, req *pb.CreateBookingRequest) (*pb.BookingResponse, error) {
	entity := ToEntity(req)
	entity.Status = "PENDING" // Initial status logic

	createdEntity, err := s.repo.Create(ctx, entity)
	if err != nil {
		// handle fail state potentially
		return nil, err
	}

	return &pb.BookingResponse{
		Booking: EntityToProto(createdEntity),
	}, nil
}

func (s *coreService) GetBooking(ctx context.Context, req *pb.GetBookingRequest) (*pb.BookingResponse, error) {
	entity, err := s.repo.GetByID(ctx, req.BookingId)
	if err != nil {
		return nil, err
	}

	return &pb.BookingResponse{
		Booking: EntityToProto(entity),
	}, nil
}

func (s *coreService) ListGuestBookings(ctx context.Context, req *pb.ListGuestBookingsRequest) (*pb.ListBookingsResponse, error) {
	entities, err := s.repo.ListByGuestID(ctx, req.GuestId)
	if err != nil {
		return nil, err
	}

	return &pb.ListBookingsResponse{
		Bookings: EntitiesToProtos(entities),
	}, nil
}
