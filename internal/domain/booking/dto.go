package booking

import (
	pb "github.com/user/airbnb-test/api/proto/v1"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func ToEntity(req *pb.CreateBookingRequest) *Entity {
	return &Entity{
		PropertyID:  req.PropertyId,
		GuestID:     req.GuestId,
		StartDate:   req.StartDate,
		EndDate:     req.EndDate,
		PaymentInfo: req.PaymentInfo,
	}
}

func ModelToEntity(m *Model) *Entity {
	return &Entity{
		BookingID:   m.BookingID,
		PropertyID:  m.PropertyID,
		GuestID:     m.GuestID,
		StartDate:   m.StartDate,
		EndDate:     m.EndDate,
		Status:      m.Status,
		PaymentInfo: m.PaymentInfo,
		CreatedAt:   m.CreatedAt,
		UpdatedAt:   m.UpdatedAt,
	}
}

func EntityToModel(e *Entity) *Model {
	return &Model{
		BookingID:   e.BookingID,
		PropertyID:  e.PropertyID,
		GuestID:     e.GuestID,
		StartDate:   e.StartDate,
		EndDate:     e.EndDate,
		Status:      e.Status,
		PaymentInfo: e.PaymentInfo,
		CreatedAt:   e.CreatedAt,
		UpdatedAt:   e.UpdatedAt,
	}
}

func EntityToProto(e *Entity) *pb.Booking {
	return &pb.Booking{
		BookingId:   e.BookingID,
		PropertyId:  e.PropertyID,
		GuestId:     e.GuestID,
		StartDate:   e.StartDate,
		EndDate:     e.EndDate,
		Status:      e.Status,
		PaymentInfo: e.PaymentInfo,
		CreatedAt:   timestamppb.New(e.CreatedAt),
		UpdatedAt:   timestamppb.New(e.UpdatedAt),
	}
}

func EntitiesToProtos(entities []*Entity) []*pb.Booking {
	protos := make([]*pb.Booking, len(entities))
	for i, e := range entities {
		protos[i] = EntityToProto(e)
	}
	return protos
}
