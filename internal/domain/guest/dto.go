package guest

import (
	pb "github.com/user/airbnb-test/api/proto/v1"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func ToEntity(req *pb.CreateGuestRequest) *Entity {
	return &Entity{
		Name:    req.Name,
		Email:   req.Email,
		Phone:   req.Phone,
		Country: req.Country,
	}
}

func ModelToEntity(m *Model) *Entity {
	return &Entity{
		GuestID:   m.GuestID,
		Name:      m.Name,
		Email:     m.Email,
		Phone:     m.Phone,
		Country:   m.Country,
		CreatedAt: m.CreatedAt,
	}
}

func EntityToModel(e *Entity) *Model {
	return &Model{
		GuestID:   e.GuestID,
		Name:      e.Name,
		Email:     e.Email,
		Phone:     e.Phone,
		Country:   e.Country,
		CreatedAt: e.CreatedAt,
	}
}

func EntityToProto(e *Entity) *pb.Guest {
	return &pb.Guest{
		GuestId:   e.GuestID,
		Name:      e.Name,
		Email:     e.Email,
		Phone:     e.Phone,
		Country:   e.Country,
		CreatedAt: timestamppb.New(e.CreatedAt),
	}
}
