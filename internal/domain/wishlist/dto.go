package wishlist

import (
	pb "github.com/user/airbnb-test/api/proto/v1"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func ToEntity(req *pb.AddToWishlistRequest) *Entity {
	return &Entity{
		GuestID:    req.GuestId,
		PropertyID: req.PropertyId,
	}
}

func ModelToEntity(m *Model) *Entity {
	return &Entity{
		ID:         m.ID,
		GuestID:    m.GuestID,
		PropertyID: m.PropertyID,
		CreatedAt:  m.CreatedAt,
	}
}

func EntityToModel(e *Entity) *Model {
	return &Model{
		ID:         e.ID,
		GuestID:    e.GuestID,
		PropertyID: e.PropertyID,
		CreatedAt:  e.CreatedAt,
	}
}

func EntityToProto(e *Entity) *pb.WishlistItem {
	return &pb.WishlistItem{
		Id:         e.ID,
		GuestId:    e.GuestID,
		PropertyId: e.PropertyID,
		CreatedAt:  timestamppb.New(e.CreatedAt),
	}
}

func EntitiesToProtos(entities []*Entity) []*pb.WishlistItem {
	protos := make([]*pb.WishlistItem, len(entities))
	for i, e := range entities {
		protos[i] = EntityToProto(e)
	}
	return protos
}
