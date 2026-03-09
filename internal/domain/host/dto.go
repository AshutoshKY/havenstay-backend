package host

import (
	pb "github.com/user/airbnb-test/api/proto/v1"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// ToEntity maps from Proto to Domain Entity
func ToEntity(req *pb.CreateHostRequest) *Entity {
	return &Entity{
		UserName: req.UserName,
		Email:    req.Email,
		Phone:    req.Phone,
		Location: req.Location,
	}
}

// ModelToEntity maps from GORM Model to Domain Entity
func ModelToEntity(m *Model) *Entity {
	return &Entity{
		ID:         m.ID,
		UserName:   m.UserName,
		Email:      m.Email,
		Phone:      m.Phone,
		Location:   m.Location,
		IsVerified: m.IsVerified,
		CreatedAt:  m.CreatedAt,
		UpdatedAt:  m.UpdatedAt,
	}
}

// EntityToModel maps from Domain Entity to GORM Model
func EntityToModel(e *Entity) *Model {
	return &Model{
		ID:         e.ID,
		UserName:   e.UserName,
		Email:      e.Email,
		Phone:      e.Phone,
		Location:   e.Location,
		IsVerified: e.IsVerified,
		CreatedAt:  e.CreatedAt,
		UpdatedAt:  e.UpdatedAt,
	}
}

// EntityToProto maps from Domain Entity to Proto Message
func EntityToProto(e *Entity) *pb.Host {
	return &pb.Host{
		Id:         e.ID,
		UserName:   e.UserName,
		Email:      e.Email,
		Phone:      e.Phone,
		Location:   e.Location,
		IsVerified: e.IsVerified,
		CreatedAt:  timestamppb.New(e.CreatedAt),
		UpdatedAt:  timestamppb.New(e.UpdatedAt),
	}
}
