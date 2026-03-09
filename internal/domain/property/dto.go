package property

import (
	pb "github.com/user/airbnb-test/api/proto/v1"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func ToEntity(req *pb.CreatePropertyRequest) *Entity {
	return &Entity{
		HostID:      req.HostId,
		Name:        req.Name,
		Location:    req.Location,
		Description: req.Description,
		Price:       req.Price,
	}
}

func ModelToEntity(m *Model) *Entity {
	return &Entity{
		ID:          m.ID,
		HostID:      m.HostID,
		Name:        m.Name,
		Location:    m.Location,
		Description: m.Description,
		Price:       m.Price,
		CreatedAt:   m.CreatedAt,
	}
}

func EntityToModel(e *Entity) *Model {
	return &Model{
		ID:          e.ID,
		HostID:      e.HostID,
		Name:        e.Name,
		Location:    e.Location,
		Description: e.Description,
		Price:       e.Price,
		CreatedAt:   e.CreatedAt,
	}
}

func EntityToProto(e *Entity) *pb.Property {
	return &pb.Property{
		Id:          e.ID,
		HostId:      e.HostID,
		Name:        e.Name,
		Location:    e.Location,
		Description: e.Description,
		Price:       e.Price,
		CreatedAt:   timestamppb.New(e.CreatedAt),
	}
}

func EntitiesToProtos(entities []*Entity) []*pb.Property {
	protos := make([]*pb.Property, len(entities))
	for i, e := range entities {
		protos[i] = EntityToProto(e)
	}
	return protos
}
