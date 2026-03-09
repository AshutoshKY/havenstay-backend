package wishlist

import (
	"context"

	pb "github.com/user/airbnb-test/api/proto/v1"
)

type Service interface {
	AddToWishlist(ctx context.Context, req *pb.AddToWishlistRequest) (*pb.WishlistResponse, error)
	GetWishlist(ctx context.Context, req *pb.GetWishlistRequest) (*pb.GetWishlistResponse, error)
}

type coreService struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &coreService{repo: repo}
}

func (s *coreService) AddToWishlist(ctx context.Context, req *pb.AddToWishlistRequest) (*pb.WishlistResponse, error) {
	entity := ToEntity(req)

	createdEntity, err := s.repo.Add(ctx, entity)
	if err != nil {
		return nil, err
	}

	return &pb.WishlistResponse{
		Item: EntityToProto(createdEntity),
	}, nil
}

func (s *coreService) GetWishlist(ctx context.Context, req *pb.GetWishlistRequest) (*pb.GetWishlistResponse, error) {
	entities, err := s.repo.ListByGuestID(ctx, req.GuestId)
	if err != nil {
		return nil, err
	}

	return &pb.GetWishlistResponse{
		Items: EntitiesToProtos(entities),
	}, nil
}
